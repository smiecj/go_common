package mysql

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/smiecj/go_common/config"
	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// distinct 用: 字段分隔符
	distinctSeparator = ";;;"

	// 配置中心中存放的 mysql 配置默认的space
	// 后续: 最好是可以由用户来控制 space 存放的位置
	mysqlConfigDefaultSpace = "mysql"

	// 默认配置: 最大空闲连接数
	defaultMaxIdleConn = 10
)

var (
	mysqlConnectorMap  map[MySQLConnectOption]RDBConnector
	mysqlConnectorLock sync.RWMutex
)

// mysql 连接配置
type MySQLConnectOption struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	IsSSL       bool   `yaml:"is_ssl"`
	MaxLifeTime int    `yaml:"max_life_time"`
	MaxIdleTime int    `yaml:"max_idle_time"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

// 对mysql 配置进行检查，不合理的配置配默认值
func (option *MySQLConnectOption) check() {
	if option.MaxLifeTime == 0 && option.MaxIdleTime == 0 {
		option.MaxLifeTime = 5 * 60
		option.MaxIdleTime = option.MaxLifeTime
	}

	if option.MaxIdleConn == 0 {
		option.MaxIdleConn = defaultMaxIdleConn
	}
}

// mysql 存储
type mysqlConnector struct {
	db *gorm.DB
}

// mysql: 插入数据
// 后续: 对批量插入场景，单次插入的数据量进行控制
func (connector *mysqlConnector) Insert(funcArr ...RDBInsertConfigFunc) (ret UpdateRet, err error) {
	// 后续: 考虑是否要适配，只传入表名，支持使用默认库名的场景
	action := MakeRDBInsertAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 判断是通过 object 插入 还是通过 直接指定key-value插入
	var dbRet *gorm.DB
	fieldArr := action.GetFieldArr()
	objectArr := action.GetObjectArr()
	if len(fieldArr) != 0 {
		keyValueMapArr := make([]map[string]interface{}, 0)
		for _, currentField := range fieldArr {
			currentKeyValueMap := make(map[string]interface{}, 0)
			for key, value := range currentField.GetMap() {
				currentKeyValueMap[key] = value
			}
			keyValueMapArr = append(keyValueMapArr, currentKeyValueMap)
		}
		dbRet = connector.db.Table(action.GetSpaceName()).Create(keyValueMapArr)
	} else if len(objectArr) != 0 {
		searchKeyArr := []string{}
		keyArr := action.GetKeyArr()
		if len(keyArr) != 0 {
			searchKeyArr = keyArr
		}
		// 注意数组类型需要转换一下，传入的 interface{} 数组无法被 gorm 识别（即数组需要保持原有的type）
		var toInsertArr interface{} = objectArr
		objectArrType := action.GetObjectArrType()
		if nil != objectArrType {
			slice := reflect.MakeSlice(objectArrType, 0, 0)
			for _, currentObj := range objectArr {
				slice = reflect.Append(slice, reflect.ValueOf(currentObj))
			}
			toInsertArr = slice.Interface()
		}
		dbRet = connector.db.Table(action.GetSpaceName()).Select(searchKeyArr).Create(toInsertArr)
	} else {
		log.Warn("[mysqlConnector.Insert] To insert data is empty")
		return ret, nil
	}

	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Insert] Insert failed: table: %s, reason: %s", action.GetSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Insert] Insert success: %s, insert rows: %d", action.GetSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 更新数据
func (connector *mysqlConnector) Update(funcArr ...RDBUpdateConfigFunc) (ret UpdateRet, err error) {
	action := MakeRDBUpdateAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 根据查询条件 更新指定数据，只更新一种取值
	var dbRet *gorm.DB
	fieldArr := action.GetFieldArr()
	objectArr := action.GetObjectArr()
	keyArr := action.GetKeyArr()
	condition := action.GetCondition()
	if len(fieldArr) != 0 {
		keyValueMap := make(map[string]interface{}, 0)
		currentField := fieldArr[0]
		for key, value := range currentField.GetMap() {
			keyValueMap[key] = value
		}
		dbRet = connector.db.Table(action.GetSpaceName()).Where(condition.WhereArr.ToSQL()).Updates(keyValueMap)
	} else if len(objectArr) != 0 {
		searchKeyArr := []string{}
		if len(keyArr) != 0 {
			searchKeyArr = keyArr
		}
		dbRet = connector.db.Table(action.GetSpaceName()).Where(condition.WhereArr.ToSQL()).Select(searchKeyArr).Updates(objectArr[0])
	} else {
		log.Warn("[mysqlConnector.Update] To update data is empty")
		return ret, nil
	}

	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Update] Update failed, table: %s, reason: %s", action.GetSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Update] Update success: %s, update rows: %d", action.GetSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 删除数据
func (connector *mysqlConnector) Delete(funcArr ...RDBDeleteConfigFunc) (ret UpdateRet, err error) {
	action := MakeRDBDeleteAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	condition := action.GetCondition()
	limitCondition := ""
	if condition.Limit != 0 {
		limitCondition = fmt.Sprintf("limit %d", condition.Limit)
	}

	dbRet := connector.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s %s",
		action.GetSpaceName(), condition.WhereArr.ToSQL(), limitCondition))
	ret.AffectedRows, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Delete] Delete failed, table: %s, reason: %s", action.GetSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Delete] Delete success: %s, update rows: %d", action.GetSpaceName(), ret.AffectedRows)
	}
	return
}

// mysql: 查询数据
func (connector *mysqlConnector) Search(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	condition := action.GetCondition()

	// 统计 count
	var count int64
	keyArr := action.GetKeyArr()
	dbRet := connector.db.Table(action.GetSpaceName()).
		Select(keyArr).Where(condition.WhereArr.ToSQL()).Count(&count)
	if nil != dbRet.Error {
		log.Error("[mysqlConnector.Count] Count failed, table: %s, reason: %s", action.GetSpaceName(), dbRet.Error.Error())
		return ret, dbRet.Error
	} else {
		ret.Total = int(count)
	}

	// order condition
	var orderStr string
	if "" != condition.Order.Field {
		orderStr = fmt.Sprintf("%s %s", condition.Order.Field, condition.Order.Sc)
	}

	objectArrType := action.GetObjectArrType()
	if nil != objectArrType {
		objectReflectArr := reflect.MakeSlice(objectArrType, 0, 0).Interface()
		dbRet = connector.db.Table(action.GetSpaceName()).
			Select(keyArr).Where(condition.WhereArr.ToSQL()).Order(orderStr).
			Offset(condition.Page.No * condition.Page.Limit).Limit(condition.Page.Limit).
			Find(&objectReflectArr)
		ret.ObjectArr = objectReflectArr
	} else {
		// 非导入到 object 情况，存在 value 在转换的时候不准确的问题，需要测试
		keyValueMapArr := make([]map[string]interface{}, 0)
		dbRet = connector.db.Table(action.GetSpaceName()).
			Select(keyArr).Where(condition.WhereArr.ToSQL()).Order(orderStr).
			Offset(condition.Page.No * condition.Page.Limit).Limit(condition.Page.Limit).
			Find(&keyValueMapArr)
		for _, currentKeyValueMap := range keyValueMapArr {
			currentField := BuildNewField()
			for key, value := range currentKeyValueMap {
				valuestr := fmt.Sprintf("%v", value)
				currentField.AddKeyValue(key, valuestr)
			}
			ret.AddField(currentField)
		}
	}

	ret.Len, err = int(dbRet.RowsAffected), dbRet.Error

	if nil != err {
		log.Error("[mysqlConnector.Select] Select failed, table: %s, reason: %s", action.GetSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Select] Select success: %s, search rows: %d", action.GetSpaceName(), ret.Len)
	}
	return
}

// mysql: 统计数据量
func (connector *mysqlConnector) Count(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	condition := action.GetCondition()
	var count int64
	dbRet := connector.db.Table(action.GetSpaceName()).Where(condition.WhereArr.ToSQL()).Count(&count)

	ret.Total, err = int(count), dbRet.Error
	ret.Len = ret.Total

	if nil != err {
		log.Error("[mysqlConnector.Count] Count failed, table: %s, reason: %s", action.GetSpaceName(), err.Error())
	} else {
		log.Info("[mysqlConnector.Count] Count success: %s, total: %d", action.GetSpaceName(), ret.Total)
	}
	return
}

// mysql: distinct
func (connector *mysqlConnector) Distinct(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}
	// distinct 必须指定需要查询的列名
	keyArr := action.GetKeyArr()
	if len(keyArr) == 0 {
		return ret, errorcode.BuildErrorWithMsg(errorcode.DBParamInvalid, "[mysqlConnector.Distinct] distinct must set key array")
	}

	fieldValueArr := make([]string, 0)
	// 查询包含多个字段，通过 SQL concat 关键字进行拼接
	var distinctColumn string
	if len(keyArr) == 1 {
		distinctColumn = keyArr[0]
	} else {
		distinctColumn = "CONCAT("
		for index := 0; index < len(keyArr); index++ {
			if index != 0 {
				distinctColumn += fmt.Sprintf(", '%s', ", distinctSeparator)
			}
			distinctColumn += fmt.Sprintf("%s", keyArr[index])
		}
		distinctColumn += ")"
	}

	dbRet := connector.db.Table(action.GetSpaceName()).Select(keyArr).Where(action.GetCondition().WhereArr.ToSQL()).
		Distinct().Pluck(distinctColumn, &fieldValueArr)
	err = dbRet.Error
	if nil != dbRet.Error {
		log.Error("[mysqlConnector.Distinct] Distinct %s get field value failed: %s", action.GetSpaceName(), err.Error())
		return ret, err
	}

	// 多种取值组合最后会放到 多个 field 中
	for _, currentValue := range fieldValueArr {
		currentField := BuildNewField()
		valueSplitArr := strings.Split(currentValue, distinctSeparator)
		for index := 0; index < len(keyArr); index++ {
			currentField.AddKeyValue(keyArr[index], valueSplitArr[index])
		}
		ret.AddField(currentField)
	}
	// distinct 只计算 len，不计算 total
	ret.Len = len(fieldValueArr)
	return
}

// 获取 mysql 连接器
func GetMySQLConnector(configManager config.Manager) (RDBConnector, error) {
	var connector RDBConnector
	mysqlConnectorLock.RLock()

	option := MySQLConnectOption{}
	configManager.Unmarshal(mysqlConfigDefaultSpace, &option)
	option.check()

	if nil == mysqlConnectorMap {
		mysqlConnectorMap = make(map[MySQLConnectOption]RDBConnector)
	}

	connector = mysqlConnectorMap[option]
	mysqlConnectorLock.RUnlock()

	if nil != connector {
		return connector, nil
	}

	mysqlConnectorLock.Lock()
	defer mysqlConnectorLock.Unlock()

	// useAffectedRows 等配置提示无效，后续需要确认原因
	extendParam := ""
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4%s",
		option.User, option.Password, option.Host, option.Port, option.Database, extendParam)
	// gorm 日志默认不打印
	db, err := gorm.Open(mysql.Open(connectStr), &gorm.Config{Logger: logger.Discard})

	// mysql 连接能成功创建，并执行 SQL, 才算是创建成功
	if nil != err {
		log.Error("[GetMySQLConnector] Get mysql connector failed, please check config: %s, err: %s", connectStr, err.Error())
		return nil, errorcode.BuildErrorWithMsg(errorcode.DBConnectFailed, err.Error())
	}

	connDB, _ := db.DB()
	connDB.SetMaxIdleConns(option.MaxIdleConn)
	connDB.SetConnMaxIdleTime(time.Second * time.Duration(option.MaxIdleTime))
	connDB.SetConnMaxLifetime(time.Second * time.Duration(option.MaxLifeTime))

	err = db.Exec("SELECT 1;").Error
	if nil != err {
		log.Error("[GetMySQLConnector] Exec mysql check sql failed, please check config: %s, err: %s", connectStr, err.Error())
		return nil, errorcode.BuildErrorWithMsg(errorcode.DBConnectFailed, err.Error())
	}
	mysqlConnector := new(mysqlConnector)
	mysqlConnector.db = db
	mysqlConnectorMap[option] = mysqlConnector
	return mysqlConnector, nil
}
