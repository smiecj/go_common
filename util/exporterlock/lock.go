package exportersql

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/log"
	"gorm.io/gorm"
)

const (
	IntervalShort = iota
	IntervalMedium
)

var (
	// exporter 超时时间（也是slave 去定期检查当前时间戳是否超时的时间）和 更新周期类型 的对应关系，单位分钟
	expireTimeMap = map[int]int64{
		IntervalShort:  10,
		IntervalMedium: 60,
	}
	// exporter 更新时间戳 和 更新周期类型 的对应关系，单位分钟
	updateIntervalMap = map[config.ExporterIntervalType]int64{
		IntervalShort:  1,
		IntervalMedium: 10,
	}

	// 当前进程占用的所有 exporter 的更新周期类型
	exporterIntervalTypeMap = make(map[string]int)

	// lock manager 单例
	lockManagerInstance *lockManager
	lockOnce            sync.Once
)

type lockManager struct {
	db db.RDBConnector
}

// 获取 lock manager 单例
// 当前: 完善lock manager 初始化逻辑
func GetExporterLockManager(db db.RDBConnector) string {
	localEnvName := system.GetEnvName()
	hasLockExporterName := ""
	intervalType = getValidIntervalType(intervalType)
	for {
		for _, exporterName := range exporterNameArr {
			isGetExporterLockSuccess, failedReason := updateExporterVersion(db, exporterName, localEnvName)
			if !isGetExporterLockSuccess {
				log.Info("[InitExporterLock] 当前任务: %s, 获取Exporter乐观锁失败原因: %s", exporterName, failedReason)
			} else {
				log.Info("[InitExporterLock] 当前任务: %s, 占锁成功！", exporterName)
				hasLockExporterName = exporterName
				break
			}
		}
		if "" != hasLockExporterName {
			break
		}
		time.Sleep(time.Duration(EXPIRE_TIME_MAP[intervalType]) * time.Minute)
	}

	log.Info("[InitExporterLock] 获取Exporter乐观锁成功！开始执行主程序")
	go func() {
		updateTimeFailedTime := 0
		for {
			time.Sleep(time.Duration(UPDATE_INTERVAL_MAP[intervalType]) * time.Minute)
			isUpdateTimeSuccess, failedReason := updateExporterTime(db, hasLockExporterName, localEnvName)
			if !isUpdateTimeSuccess {
				log.Info("[updateExporterTime] master节点更新时间失败, 失败原因: %s", failedReason)
			}
			if !isUpdateTimeSuccess {
				updateTimeFailedTime++
			} else {
				updateTimeFailedTime = 0
			}
			// 连续3次更新状态失败，将认为之前占的锁已经失效，强行退出程序
			if updateTimeFailedTime > 3 {
				log.Info("[updateExporterTime] 当前已经连续3次更新时间失败，直接退出！")
				os.Exit(1)
			}
		}
	}()
	exporterIntervalTypeMap[hasLockExporterName] = intervalType
	return hasLockExporterName
}

/**
 * exporter 启动或者检查是否需要抢占时执行
 * 保证 exporter 服务部署到Sumeru 之后，同一时间，只有一个节点的exporter 在采集数据，
 * 避免多个节点都在采集同样的数据，造成资源浪费
 */
func updateExporterVersion(db *gorm.DB, exporterName string, localIp string) (bool, string) {
	// 查询 -> 插入 -> 更新，只有到最后一步成功了才算成功
	currentTime := time.Now()
	currentTimeStr := currentTime.Format(timeutil.TIME_NORMAL_FORMAT)

	exporterLock := new(monitorcfgmodel.ExporterLock)
	dbErr := db.Where("name = ?", exporterName).First(&exporterLock).Error
	log.Info("[updateExporterVersion] 获取指定任务的当前状态: %s", exporterName)

	var needUpdate bool

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		log.Error("[updateExporterVersion] 查询字段失败，失败原因: %s, exporter name: %s",
			dbErr.Error(), exporterName)
		return false, "查询语句执行失败"
	}
	if "" != exporterLock.Name {
		// 之前不是本机占用的: 判断是否过期; 是本机占用的: 直接更新
		if exporterLock.Ip != localIp {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			parsedTime, _ := time.ParseInLocation(timeutil.TIME_NORMAL_FORMAT, exporterLock.UpdateTime, loc)
			dur := currentTime.Unix() - parsedTime.Unix()
			// 已经超过10分钟没有更新
			if dur/60 > EXPIRE_TIME_MAP[exporterIntervalTypeMap[exporterName]] {
				needUpdate = true
			} else {
				needUpdate = false
			}
		} else {
			needUpdate = true
		}
	} else {
		// 一条数据都没有：先插入数据
		exporterLock.InitExporterLock(exporterName, currentTimeStr, localIp)
		insertRet := db.Create(exporterLock)
		dbErr, rowsAffected := insertRet.Error, insertRet.RowsAffected

		if dbErr != nil {
			return false, "插入数据失败"
		}
		if rowsAffected > 0 {
			return true, ""
		} else {
			return false, "插入数据没成功，可能是已经有节点插入成功了"
		}
	}
	if needUpdate {
		oldVersion := exporterLock.Version
		exporterLock.Version++
		exporterLock.UpdateTime, exporterLock.Ip = currentTimeStr, localIp

		rowsAffected := db.Model(&exporterLock).Where("name = ? AND version = ?",
			exporterLock.Name, oldVersion).Select([]string{"update_time", "ip", "version"}).
			Updates(exporterLock).RowsAffected
		log.Info("[updateExporterVersion] exporter name: %s, version: %d, env: %s, update ret: %d",
			exporterLock.Name, exporterLock.Version, exporterLock.Ip, rowsAffected)

		if rowsAffected == 1 {
			return true, ""
		} else {
			log.Error("[updateExporterTime] 更新节点更新时间失败: env: %s, exporter name: %s, time: %s",
				localIp, exporterName, exporterLock.UpdateTime)
			return false, "更新数据没成功，可能是已经有节点更新成功了"
		}
	} else {
		return false, "当前已经有exporter 在执行，不需要更新"
	}
}

/**
 * 已经占用成功的exporter 执行
 * 仅更新时间，保活作用
 */
func updateExporterTime(db *gorm.DB, exporterName string, localIp string) (bool, string) {
	currentTime := time.Now()
	currentTimeStr := currentTime.Format(timeutil.TIME_NORMAL_FORMAT)

	exporterLock := new(monitorcfgmodel.ExporterLock)
	exporterLock.InitExporterLock(exporterName, currentTimeStr, localIp)
	rowsAffected := db.Model(&exporterLock).Where("name = ? AND ip = ?", exporterName, localIp).
		Select([]string{"update_time", "ip"}).Updates(exporterLock).RowsAffected

	if rowsAffected == 1 {
		return true, ""
	} else {
		return false, "当前节点更新时间失败，可能代码有逻辑问题，请检查"
	}
}

// 公共方法: 返回一个合法的exporter 更新周期类型，默认为短周期更新
func getValidIntervalType(intervalType config.ExporterIntervalType) config.ExporterIntervalType {
	if config.EXPORTER_UPDATE_INTERVAL_SHORT != intervalType && config.EXPORTER_UPDATE_INTERVAL_MEDIUM == intervalType {
		intervalType = config.EXPORTER_UPDATE_INTERVAL_SHORT
	}
	return intervalType
}
