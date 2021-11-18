package db

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// distinct 用: 字段分隔符
	distinctSeparator = ";;;"
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
		// 注意数组类型需要转换一下，传入的 interface{} 数组无法被 gorm 识别
		var toInsertArr interface{} = action.objectArr
		if nil != action.objectArrType {
			slice := reflect.MakeSlice(action.objectArrType, 0, 0)
			for _, currentObj := range action.objectArr {
				slice = reflect.Append(slice, reflect.ValueOf(currentObj))
			}
			toInsertArr = slice.Interface()
		}
		dbRet = connector.db.Table(action.getSpaceName()).Select(searchKeyArr).Create(toInsertArr)
	} else {
		return ret, errorcode.BuildErrorWithMsg(errorcode.DBParamInvalid, "Insert failed: to insert data is empty")
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
		dbRet = connector.db.Table(action.getSpaceName()).Where(action.condition.WhereArr.toSQL()).Select(searchKeyArr).Updates(action.objectArr[0])
	} else {
		return ret, errorcode.BuildErrorWithMsg(errorcode.DBParamInvalid, "Insert failed: to insert data is empty")
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

	dbRet := connector.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s %s",
		action.getSpaceName(), action.condition.WhereArr.toSQL(), limitCondition))
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

	// 统计 count
	var count int64
	dbRet := connector.db.Table(action.getSpaceName()).
		Select(action.keyArr).Where(action.condition.WhereArr.toSQL()).Count(&count)
	if nil != dbRet.Error {
		log.Error("[mysqlConnector.Count] Count failed, table: %s, reason: %s", action.getSpaceName(), err.Error())
		return ret, dbRet.Error
	} else {
		ret.Total = int(count)
	}

	if nil != action.objectArrType {
		objectReflectArr := reflect.MakeSlice(action.objectArrType, 0, 0).Interface()
		dbRet = connector.db.Table(action.getSpaceName()).
			Select(action.keyArr).Where(action.condition.WhereArr.toSQL()).
			Offset(action.condition.Page.No * action.condition.Page.Limit).Limit(action.condition.Page.Limit).
			Find(&objectReflectArr)
		ret.ObjectArr = objectReflectArr
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
			ret.addField(currentField)
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

// mysql: 统计数据量
func (connector *mysqlConnector) Count(funcArr ...rdbSearchConfigFunc) (ret searchRet, err error) {
	action := makeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	var count int64
	dbRet := connector.db.Table(action.getSpaceName()).Where(action.condition.WhereArr.toSQL()).Count(&count)

	ret.Total, err = int(count), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Count] Count failed, table: %s, reason: %s", action.getSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Count] Count success: %s, total: %d", action.getSpaceName(), ret.Total)
	}
	return
}

// mysql: distinct
func (connector *mysqlConnector) Distinct(funcArr ...rdbSearchConfigFunc) (ret searchRet, err error) {
	action := makeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}
	// distinct 必须指定需要查询的列名
	if len(action.keyArr) == 0 {
		return ret, errorcode.BuildErrorWithMsg(errorcode.DBParamInvalid, "[mysqlConnector.Distinct] distinct must set key array")
	}

	fieldValueArr := make([]string, 0)
	// 查询包含多个字段，通过 SQL concat 关键字进行拼接
	var distinctColumn string
	if 1 == len(action.keyArr) {
		distinctColumn = action.keyArr[0]
	} else {
		distinctColumn = "CONCAT("
		for index := 0; index < len(action.keyArr); index++ {
			if index != 0 {
				distinctColumn += fmt.Sprintf(", '%s', ", distinctSeparator)
			}
			distinctColumn += fmt.Sprintf("%s", action.keyArr[index])
		}
		distinctColumn += ")"
	}

	dbRet := connector.db.Table(action.getSpaceName()).Select(action.keyArr).Where(action.condition.WhereArr.toSQL()).
		Distinct().Pluck(distinctColumn, &fieldValueArr)
	err = dbRet.Error
	if nil != dbRet.Error {
		log.Error("[mysqlConnector.Distinct] Distinct %s get field value failed: %s", action.getSpaceName(), err.Error())
		return ret, err
	}

	// 多种取值组合最后会放到 多个 field 中
	for _, currentValue := range fieldValueArr {
		currentField := BuildNewField()
		valueSplitArr := strings.Split(currentValue, distinctSeparator)
		for index := 0; index < len(action.keyArr); index++ {
			currentField.AddKeyValue(action.keyArr[index], valueSplitArr[index])
		}
		ret.addField(currentField)
	}
	// distinct 只计算 len，不计算 total
	ret.Len = len(fieldValueArr)
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

	// useAffectedRows 等配置提示无效，后续需要确认原因
	extendParam := ""
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4%s",
		option.User, option.Password, option.Host, option.Port, extendParam)
	// gorm 日志默认不打印
	db, err := gorm.Open(mysql.Open(connectStr), &gorm.Config{Logger: logger.Discard})

	// mysql 连接能成功创建，并执行 SQL, 才算是创建成功
	if nil != err {
		log.Error("[GetMySQLConnector] Get mysql connector failed, please check config: %s, err: %s", connectStr, err.Error())
		return nil
	}
	err = db.Exec("SELECT 1;").Error
	if nil != err {
		log.Error("[GetMySQLConnector] Exec mysql check sql failed, please check config: %s, err: %s", connectStr, err.Error())
		return nil
	}
	mysqlConnector := new(mysqlConnector)
	mysqlConnector.db = db
	mysqlConnectorMap[option] = mysqlConnector
	return mysqlConnector
}
