package db

import (
	"sync"

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
func (connector *localMemoryConnector) Insert(funcArr ...rdbInsertConfigFunc) (updateRet, error) {
	action := makeRDBInsertAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 本地存储仅会存入 key-value 格式的数据，并且会覆盖
	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		connector.storage[spaceName] = make(map[string]string)
	}
	for _, currentField := range action.fieldArr {
		for key, value := range currentField.keyValueMap {
			connector.storage[spaceName][key] = value
		}
	}

	return updateRet{AffectedRows: len(action.fieldArr)}, nil
}

// 本地存储: 更新数据
func (connector *localMemoryConnector) Update(funcArr ...rdbUpdateConfigFunc) (updateRet, error) {
	action := makeRDBUpdateAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 更新的逻辑和存储一样
	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		connector.storage[spaceName] = make(map[string]string)
	}
	for _, currentField := range action.fieldArr {
		for key, value := range currentField.keyValueMap {
			connector.storage[spaceName][key] = value
		}
	}

	return updateRet{AffectedRows: len(action.fieldArr)}, nil
}

// 本地存储: 删除数据
func (connector *localMemoryConnector) Delete(funcArr ...rdbDeleteConfigFunc) (updateRet, error) {
	action := makeRDBDeleteAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 删除: 忽略查询条件，直接删除整个space
	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		return updateRet{AffectedRows: 0}, nil
	} else {
		affectedRows := len(connector.storage[spaceName])
		delete(connector.storage, spaceName)
		return updateRet{AffectedRows: affectedRows}, nil
	}
}

// 本地存储: 查询数据
func (connector *localMemoryConnector) Search(funcArr ...rdbSearchConfigFunc) (searchRet, error) {
	action := makeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 查询: 按照 spaceName 查询
	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		return searchRet{Len: 0}, nil
	} else {
		currentField := BuildNewField()
		for key, value := range connector.storage[spaceName] {
			currentField.AddKeyValue(key, value)
		}
		return searchRet{Len: 1, FieldArr: []field{currentField}}, nil
	}
}

// 本地内存: 统计数据量
func (connector *localMemoryConnector) Count(funcArr ...rdbSearchConfigFunc) (searchRet, error) {
	action := makeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		return searchRet{Total: 0}, nil
	} else {
		return searchRet{Total: 1}, nil
	}
}

// 本地内存: 暂不需要实现 Distinct
func (connector *localMemoryConnector) Distinct(funcArr ...rdbSearchConfigFunc) (ret searchRet, err error) {
	return ret, errorcode.BuildError(errorcode.NotImplement, "[localMemoryConnector.Distinct] not implement")
}

// 实现本地内存连接器
func GetLocalMemoryConnector() RDBConnector {
	localMemoryConnectorOnce.Do(func() {
		localConnector := new(localMemoryConnector)
		localConnector.storage = make(map[string]map[string]string)
		localMemoryConnectorSingleton = localConnector
	})
	return localMemoryConnectorSingleton
}

// 后续: 初始化 连接器配置中，增加 id generator 配置
