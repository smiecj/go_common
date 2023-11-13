package local

import (
	"sync"

	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/errorcode"
)

var (
	localMemoryConnectorSingleton RDBConnector
	localMemoryConnectorOnce      sync.Once
)

// 本地内存存储
type localMemoryConnector struct {
	// db_name -> key: table name; id: uuid
	storage map[string]map[string]string
}

func (connector *localMemoryConnector) init() {
	connector.storage = make(map[string]map[string]string)
}

// 本地存储: 插入数据
func (connector *localMemoryConnector) Insert(funcArr ...RDBInsertConfigFunc) (UpdateRet, error) {
	action := MakeRDBInsertAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 本地存储仅会存入 key-value 格式的数据，并且会覆盖
	spaceName := action.GetSpaceName()
	fieldArr := action.GetFieldArr()
	if nil == connector.storage[spaceName] {
		connector.storage[spaceName] = make(map[string]string)
	}
	for _, currentField := range fieldArr {
		for key, value := range currentField.GetMap() {
			connector.storage[spaceName][key] = value
		}
	}

	return UpdateRet{AffectedRows: len(fieldArr)}, nil
}

// 本地存储: 更新数据
func (connector *localMemoryConnector) Update(funcArr ...RDBUpdateConfigFunc) (UpdateRet, error) {
	action := MakeRDBUpdateAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 更新的逻辑和存储一样
	spaceName := action.GetSpaceName()
	fieldArr := action.GetFieldArr()
	if nil == connector.storage[spaceName] {
		connector.storage[spaceName] = make(map[string]string)
	}
	for _, currentField := range fieldArr {
		for key, value := range currentField.GetMap() {
			connector.storage[spaceName][key] = value
		}
	}

	return UpdateRet{AffectedRows: len(fieldArr)}, nil
}

// 本地存储: 删除数据
func (connector *localMemoryConnector) Delete(funcArr ...RDBDeleteConfigFunc) (UpdateRet, error) {
	action := MakeRDBDeleteAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 删除: 忽略查询条件，直接删除整个space
	spaceName := action.GetSpaceName()
	if nil == connector.storage[spaceName] {
		return UpdateRet{AffectedRows: 0}, nil
	} else {
		affectedRows := len(connector.storage[spaceName])
		delete(connector.storage, spaceName)
		return UpdateRet{AffectedRows: affectedRows}, nil
	}
}

// 备份数据
func (connector *localMemoryConnector) Backup(funcArr ...RDBBackupConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[localMemoryConnector.Backup] not implement")
}

// 本地存储: 查询数据
func (connector *localMemoryConnector) Search(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 查询: 按照 spaceName 查询
	spaceName := action.GetSpaceName()
	if nil != connector.storage[spaceName] {
		currentField := BuildNewField()
		for key, value := range connector.storage[spaceName] {
			currentField.AddKeyValue(key, value)
		}
		ret.Len = 1
		ret.AddField(currentField)
	}
	return
}

// 本地内存: 统计数据量
func (connector *localMemoryConnector) Count(funcArr ...RDBSearchConfigFunc) (SearchRet, error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	spaceName := action.GetSpaceName()
	if nil == connector.storage[spaceName] {
		return SearchRet{Total: 0}, nil
	} else {
		return SearchRet{Total: 1}, nil
	}
}

// 本地内存: 暂不需要实现 Distinct
func (connector *localMemoryConnector) Distinct(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[localMemoryConnector.Distinct] not implement")
}

// close
func (connector *localMemoryConnector) Close() error {
	connector.storage = nil
	return nil
}

// stat
func (connector *localMemoryConnector) Stat() (ret DBStat, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.DBStatFailed, err.Error())
}

// 实现本地内存连接器
func GetLocalMemoryConnector() (RDBConnector, error) {
	localMemoryConnectorOnce.Do(func() {
		localConnector := new(localMemoryConnector)
		localConnector.storage = make(map[string]map[string]string)
		localMemoryConnectorSingleton = localConnector
	})
	return localMemoryConnectorSingleton, nil
}

func (connector *localMemoryConnector) ExecSearch(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[localMemoryConnector.ExecSearch] not implement")
}

// 后续: 初始化 连接器配置中，增加 id generator 配置
