package clusterlock

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
	"github.com/smiecj/go_common/util/net"
	timeutil "github.com/smiecj/go_common/util/time"
)

type IntervalType int

const (
	IntervalShort IntervalType = iota
	IntervalMedium

	// 表
	dbMeta    = "d_meta"
	tableLock = "t_lock"
)

var (
	// lock 超时时间（也是slave 去定期检查当前时间戳是否超时的时间）和 更新周期类型 的对应关系，单位分钟
	expireTimeMap = map[IntervalType]time.Duration{
		IntervalShort:  10 * time.Minute,
		IntervalMedium: time.Hour,
	}
	// lock 更新时间戳 和 更新周期类型 的对应关系，单位分钟
	updateIntervalMap = map[IntervalType]time.Duration{
		IntervalShort:  time.Minute,
		IntervalMedium: 10 * time.Minute,
	}

	// lock manager 单例
	lockManagerInstance *lockManager
	lockOnce            sync.Once
)

// 锁管理对象
type lockManager struct {
	connector db.RDBConnector
	envName   string
	// 当前进程占用的所有 lock 的更新周期类型
	lockNameIntervalTypeMap map[string]IntervalType
	mapLock                 sync.RWMutex

	// 返回给调用方的错误信息
	errorChan chan<- error
}

// 锁对象
type lock struct {
	Name       string `gorm:"column:name"`
	Version    int64  `gorm:"column:version"`
	UpdateTime string `gorm:"column:update_time"`
	Env        string `gorm:"column:env`
}

// 获取 lock manager 单例
func GetLockManager(connector db.RDBConnector, errorChan chan<- error) *lockManager {
	lockOnce.Do(func() {
		lockManagerInstance = new(lockManager)
		lockManagerInstance.connector = connector
		lockManagerInstance.envName, _ = net.GetLocalIPV4()
		lockManagerInstance.errorChan = errorChan
		lockManagerInstance.lockNameIntervalTypeMap = make(map[string]IntervalType)
	})
	return lockManagerInstance
}

/**
 * 入口: 占用指定名称的锁
 * 会一直持续 直到成功占用到锁才返回
 * 后续: 控制同时可进行的占锁任务数量
 */
func (manager *lockManager) Lock(lockName string, intervalType IntervalType, job func(ctx context.Context) error) {
	// 判断当前锁是否已经占用过，不会重复占用
	manager.mapLock.Lock()
	if _, ok := manager.lockNameIntervalTypeMap[lockName]; ok {
		log.Info("[Lock] lock name: %s, 程序已经占用", lockName)
		return
	}
	manager.lockNameIntervalTypeMap[lockName] = intervalType
	manager.mapLock.Unlock()

	for {
		lockErr := manager.tryLock(lockName)
		if nil != lockErr {
			log.Warn("[Lock] lock: %s, 获取lock乐观锁失败原因: %s", lockName, lockErr.Error())
		} else {
			log.Info("[Lock] lock: %s, 占锁成功！", lockName)
			break
		}
		time.Sleep(expireTimeMap[intervalType])
	}

	log.Info("[lock] 获取Lock乐观锁成功")

	log.Info("[Lock] 启动用户任务")
	ctx, cancel := context.WithCancel(context.Background())
	go job(ctx)

	log.Info("[Lock] 启动更新 lock 状态任务")
	go func() {
		updateTimeFailedTime := 0
		for {
			time.Sleep(updateIntervalMap[intervalType])
			updateErr := manager.updateLockTime(lockName)
			if updateErr != nil {
				log.Warn("[updateLockTime] master节点更新时间失败, 失败原因: %s", updateErr.Error())
				updateTimeFailedTime++
			} else {
				updateTimeFailedTime = 0
			}
			// 连续3次更新状态失败，将认为之前占的锁已经失效，关闭当前的 channel 并返回错误信息
			if updateTimeFailedTime > 3 {
				log.Warn("[updateLockTime] lock name: %s, 已经连续3次更新失败，将取消更新，并停止任务", lockName)
				cancel()
				go func() {
					// 向调用方发送错误信息
					manager.errorChan <- errorcode.BuildErrorWithMsg(errorcode.ServiceError,
						fmt.Sprintf("lock: %s, update failed", lockName))
				}()
				go func() {
					// 过超时时间后清空信息
					time.Sleep(expireTimeMap[manager.lockNameIntervalTypeMap[lockName]])
					manager.mapLock.Lock()
					delete(manager.lockNameIntervalTypeMap, lockName)
					manager.mapLock.Unlock()
				}()
				return
			}
		}
	}()
}

/**
 * 占用分布式锁
 */
func (manager *lockManager) tryLock(lockName string) error {
	connector := manager.connector
	localEnv := manager.envName
	countRet, dbErr := connector.Count(db.SearchSetSpace(dbMeta, tableLock), db.SearchSetCondition("name", "=", lockName))

	if dbErr != nil {
		return dbErr
	}

	needUpdateLockVersion := false
	var oldLock lock
	if countRet.Len != 0 {
		// 之前不是本机占用的: 判断是否过期; 是本机占用的: 直接更新
		searchRet, dbErr := connector.Search(db.SearchSetSpace(dbMeta, tableLock),
			db.SearchSetCondition("name", "=", lockName), db.SearchSetObjectArrType([]lock{}))
		if nil != dbErr || searchRet.Len != 1 {
			return dbErr
		}

		oldLock = searchRet.ObjectArr.([]lock)[0]
		if oldLock.Env != localEnv {
			dur, _ := timeutil.CompareTimestampWithNow(oldLock.UpdateTime)
			if dur > expireTimeMap[manager.lockNameIntervalTypeMap[lockName]] {
				needUpdateLockVersion = true
			}
		} else {
			// 锁在之前就是由当前节点占用的，直接更新
			needUpdateLockVersion = true
		}
	} else {
		// 没有锁记录, 直接插入锁
		newLock := lock{Name: lockName, UpdateTime: timeutil.CurrentTimestamp(), Env: localEnv}
		insertRet, dbErr := connector.Insert(
			db.InsertSetSpace(dbMeta, tableLock),
			db.InsertAddKeyArr([]string{"name", "update_time", "env"}),
			db.InsertAddObjectArr([]lock{newLock}),
		)

		if dbErr != nil {
			return dbErr
		}
		if insertRet.AffectedRows > 0 {
			return nil
		} else {
			err := errorcode.BuildErrorWithMsg(errorcode.InnerError,
				fmt.Sprintf("lock insert failed, affected row: %d", insertRet.AffectedRows))
			return err
		}
	}
	// 锁过期 or 同个节点直接占用锁
	if needUpdateLockVersion {
		return manager.updateLockVersion(oldLock)
	} else {
		return errorcode.BuildErrorWithMsg(errorcode.InnerError,
			fmt.Sprintf("lock failed, maybe another node lock %s", lockName))
	}
}

/**
 * 更新锁的版本信息，一般在激活之前失效的锁的时候使用
 */
func (manager *lockManager) updateLockVersion(oldLock lock) error {
	// copy lock
	newLock := oldLock
	newLock.Version++
	newLock.UpdateTime, newLock.Env = timeutil.CurrentTimestamp(), manager.envName

	updateRet, dbErr := manager.connector.Update(
		db.UpdateSetSpace(dbMeta, tableLock),
		db.UpdateSetCondition("name", "=", oldLock.Name, "and", "version", "=", strconv.Itoa(int(oldLock.Version))),
		db.UpdateAddKeyArr([]string{"update_time", "env", "version"}),
		db.UpdateAddObject(newLock),
	)
	if nil != dbErr {
		return dbErr
	}

	if updateRet.AffectedRows == 1 {
		return nil
	} else {
		return errorcode.BuildErrorWithMsg(errorcode.InnerError, fmt.Sprintf("lock update failed, name: %s", oldLock.Name))
	}
}

/**
 * 已经占用成功的 lock 执行
 * 仅更新时间，保活作用
 */
func (manager *lockManager) updateLockTime(lockName string) error {
	lock := lock{Name: lockName, UpdateTime: timeutil.CurrentTimestamp(), Env: manager.envName}
	updateRet, dbErr := manager.connector.Update(
		db.UpdateSetSpace(dbMeta, tableLock),
		db.UpdateSetCondition("name", "=", lockName, "and", "env", "=", manager.envName),
		db.UpdateAddKeyArr([]string{"update_time"}),
		db.UpdateAddObject(lock),
	)
	if nil != dbErr {
		return dbErr
	}

	if updateRet.AffectedRows == 1 {
		return nil
	} else {
		return errorcode.BuildErrorWithMsg(errorcode.InnerError, fmt.Sprintf("lock update time failed, name: %s", lockName))
	}
}
