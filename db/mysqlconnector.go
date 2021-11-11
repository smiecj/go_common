package db

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlConnectorMap  map[MySQLConnectOption]RDBConnector
	mysqlConnectorLock sync.RWMutex
)

// mysql 连接配置
type MySQLConnectOption struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	IsSSL    bool
}

// mysql 存储
type mysqlConnector struct {
	db *gorm.DB
}

// mysql: 插入数据
// 后续: 对批量插入场景，单次插入的数据量进行控制
func (connector *mysqlConnector) Insert(funcArr ...rdbInsertConfigFunc) (ret updateRet, err error) {
	action := makeRDBInsertAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 判断是通过 object 插入 还是通过 直接指定key-value插入
	var dbRet *gorm.DB
	if len(action.fieldArr) != 0 {
		keyValueMapArr := make([]map[string]interface{}, 0)
		for _, currentField := range action.fieldArr {
			currentKeyValueMap := make(map[string]interface{}, 0)
			for key, value := range currentField.keyValueMap {
				currentKeyValueMap[key] = value
			}
			keyValueMapArr = append(keyValueMapArr, currentKeyValueMap)
		}
		dbRet = connector.db.Table(action.getSpaceName()).Create(keyValueMapArr)
	} else if len(action.objectArr) != 0 {
		searchKeyArr := []string{}
		if len(action.keyArr) != 0 {
			searchKeyArr = action.keyArr
		}
		dbRet = connector.db.Table(action.getSpaceName()).Select(searchKeyArr).Create(action.objectArr)
	} else {
		return ret, errorcode.BuildError(errorcode.DBParamInvalid, "Insert failed: to insert data is empty")
	}

	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Insert] Insert failed: table: %s, reason: %s", action.getSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Insert] Insert success: %s, insert rows: %d", action.getSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 更新数据
func (connector *mysqlConnector) Update(funcArr ...rdbUpdateConfigFunc) (ret updateRet, err error) {
	action := makeRDBUpdateAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 根据查询条件 更新指定数据，只更新一种取值
	var dbRet *gorm.DB
	if len(action.fieldArr) != 0 {
		keyValueMap := make(map[string]interface{}, 0)
		currentField := action.fieldArr[0]
		for key, value := range currentField.keyValueMap {
			keyValueMap[key] = value
		}
		dbRet = connector.db.Table(action.getSpaceName()).Where(action.condition.WhereArr.toSQL()).Updates(keyValueMap)
	} else if len(action.objectArr) != 0 {
		searchKeyArr := []string{}
		if len(action.keyArr) != 0 {
			searchKeyArr = action.keyArr
		}
		dbRet = connector.db.Table(action.getSpaceName()).Select(searchKeyArr).Updates(action.objectArr[0])
	} else {
		return ret, errorcode.BuildError(errorcode.DBParamInvalid, "Insert failed: to insert data is empty")
	}

	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Update] Update failed, table: %s, reason: %s", action.getSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Update] Update success: %s, update rows: %d", action.getSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 删除数据
func (connector *mysqlConnector) Delete(funcArr ...rdbDeleteConfigFunc) (ret updateRet, err error) {
	action := makeRDBDeleteAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}
	limitCondition := ""
	if action.condition.Limit != 0 {
		limitCondition = fmt.Sprintf("limit %d", action.condition.Limit)
	}

	dbRet := connector.db.Exec("DELETE FROM %s WHERE %s %s",
		action.getSpaceName(), action.condition.WhereArr.toSQL(), limitCondition)
	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Delete] Delete failed, table: %s, reason: %s", action.getSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Delete] Delete success: %s, update rows: %d", action.getSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 查询数据
func (connector *mysqlConnector) Search(funcArr ...rdbSearchConfigFunc) (ret searchRet, err error) {
	action := makeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	var dbRet *gorm.DB

	if nil != action.object {
		objectArr := reflect.MakeSlice(reflect.TypeOf(action.object), 0, 0)
		dbRet = connector.db.Table(action.getSpaceName()).
			Select(action.keyArr).Where(action.condition.WhereArr.toSQL()).
			Offset(action.condition.Page.No * action.condition.Page.Limit).Limit(action.condition.Page.Limit).
			Find(&objectArr)
		ret.ObjectArr = objectArr.Interface().([]interface{})
	} else {
		// 非导入到 object 情况，存在 value 在转换的时候不准确的问题，需要测试
		keyValueMapArr := make([]map[string]interface{}, 0)
		dbRet = connector.db.Table(action.getSpaceName()).
			Select(action.keyArr).Where(action.condition.WhereArr.toSQL()).
			Offset(action.condition.Page.No * action.condition.Page.Limit).Limit(action.condition.Page.Limit).
			Find(&keyValueMapArr)
		for _, currentKeyValueMap := range keyValueMapArr {
			currentField := BuildNewField()
			for key, value := range currentKeyValueMap {
				valuestr := fmt.Sprintf("%v", value)
				currentField.AddKeyValue(key, valuestr)
			}
			ret.FieldArr = append(ret.FieldArr, currentField)
		}
	}

	ret.Len, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Select] Select failed, table: %s, reason: %s", action.getSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Select] Select success: %s, search rows: %d", action.getSpaceName(), ret.Len)
	}
	return
}

// 获取 mysql 连接器
func GetMySQLConnector(option MySQLConnectOption) RDBConnector {
	var connector RDBConnector
	mysqlConnectorLock.RLock()
	if nil == fileConnectorMap {
		mysqlConnectorMap = make(map[MySQLConnectOption]RDBConnector)
	}
	connector = mysqlConnectorMap[option]
	mysqlConnectorLock.RUnlock()

	if nil != connector {
		return connector
	}

	fileConnectorLock.Lock()
	defer fileConnectorLock.Unlock()

	extendParam := ""
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4%s",
		option.User, option.Password, option.Host, option.Port, option.Database, extendParam)
	db, err := gorm.Open(mysql.Open(connectStr), &gorm.Config{})

	// mysql 连接能成功创建，并执行 SQL, 才算是创建成功
	if nil != err {
		log.Error("[GetMySQLConnector] Get mysql connector failed, please check config: %s", connectStr)
		return nil
	}
	err = db.Exec("SELECT 1;").Error
	if nil != err {
		log.Error("[GetMySQLConnector] Exec mysql check sql failed, please check config: %s", connectStr)
		return nil
	}
	mysqlConnector := new(mysqlConnector)
	mysqlConnector.db = db
	mysqlConnectorMap[option] = mysqlConnector
	return mysqlConnector
}
