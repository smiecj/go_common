// package impala impala 数据库连接器，可进行基本的count查询
// 后续实现 对 struct 进行赋值（需要了解 struct tag 的机制）
package impala

import (
	"context"
	"database/sql"
	"strconv"
	"sync"

	"database/sql/driver"

	api "github.com/bippio/go-impala"
	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

var (
	impalaConnectorMap  map[ImpalaConnectOption]RDBConnector
	impalaConnectorLock sync.RWMutex
)

// impala 连接配置
// 后续 可支持传入用户 进行校验操作
type ImpalaConnectOption struct {
	Host string
	Port int
}

// impala 连接器
type impalaConnector struct {
	innerConnector driver.Connector
}

// todo: 实现 impala 基本连接 和 count 查询
// 插入: 暂不实现
func (connector *impalaConnector) Insert(funcArr ...RDBInsertConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Insert] not implement")
}

// 更新: 暂不实现
func (connector *impalaConnector) Update(funcArr ...RDBUpdateConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Update] not implement")
}

// 删除: 暂不实现
func (connector *impalaConnector) Delete(funcArr ...RDBDeleteConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Delete] not implement")
}

// 备份: 暂不实现
func (connector *impalaConnector) Backup(funcArr ...RDBBackupConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Backup] not implement")
}

// 查询: 暂不实现
func (connector *impalaConnector) Search(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Search] not implement")
}

// 计数
func (connector *impalaConnector) Count(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	db := sql.OpenDB(connector.innerConnector)
	defer db.Close()

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM "+action.GetSpaceName())
	if err != nil {
		return ret, errorcode.BuildErrorWithMsg(errorcode.DBExecFailed, err.Error())
	}
	defer rows.Close()

	var dataCount int
	if rows.Next() {
		rows.Scan(&dataCount)
	}
	ret.Total = dataCount
	return
}

// distinct: 暂不实现
func (connector *impalaConnector) Distinct(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[impalaConnector.Distinct] not implement")
}

// 获取 impala 连接器
func GetImpalaConnector(option ImpalaConnectOption) (RDBConnector, error) {
	var connector RDBConnector
	impalaConnectorLock.RLock()
	if nil == impalaConnectorMap {
		impalaConnectorMap = make(map[ImpalaConnectOption]RDBConnector)
	}
	connector = impalaConnectorMap[option]
	impalaConnectorLock.RUnlock()

	if nil != connector {
		return connector, nil
	}

	impalaConnectorLock.Lock()
	defer impalaConnectorLock.Unlock()

	opts := api.DefaultOptions
	opts.Host = option.Host
	opts.Port = strconv.Itoa(option.Port)
	opts.QueryTimeout = 5

	innerConnector := api.NewConnector(&opts)
	db := sql.OpenDB(innerConnector)
	defer db.Close()

	ctx := context.Background()
	_, err := db.QueryContext(ctx, "SHOW DATABASES")
	if nil != err {
		log.Error("[GetImpalaConnector] Get Impala connector failed: %s", err.Error())
		return nil, errorcode.BuildErrorWithMsg(errorcode.DBConnectFailed, err.Error())
	}

	ret := new(impalaConnector)
	ret.innerConnector = innerConnector
	impalaConnectorMap[option] = ret
	return ret, nil
}
