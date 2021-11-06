package db

import "fmt"

// Relational data connector
type RDBConnector interface {
	Insert(...rdbInsertConfigFunc) (updateRet, error)
	Update(...rdbUpdateConfigFunc) (updateRet, error)
	Delete(...rdbDeleteConfigFunc) (updateRet, error)
	Search(...rdbSearchConfigFunc) (searchRet, error)
}

// 更新类型动作结果
type updateRet struct {
	AffectedRows int
}

// 查询类型动作结果
type searchRet struct {
	ObjectArr []interface{}
	FieldArr  []field
	Page      int
	Total     int
}

// 库表空间定义
type space struct {
	db    string
	table string
}

// 获取库名.表名的格式
func (space *space) getSpaceName() string {
	return fmt.Sprintf("%s.%s", space.db, space.table)
}

// 数据字段定义
type field struct {
	keyValueMap map[string]string
}

// 库表属性定义（包括表字段）
type rdbField struct {
	space
	objectArr []interface{}
	fieldArr  []field
}

// 添加字段和对应值
func (rdbField *rdbField) addField(keyValueMap map[string]string) {
	newKeyValueMap := make(map[string]string, 0)
	for key, value := range keyValueMap {
		newKeyValueMap[key] = value
	}
	rdbField.fieldArr = append(rdbField.fieldArr, field{keyValueMap: keyValueMap})
}

// 添加一整个结构体
func (rdbField *rdbField) addObject(object interface{}) {
	rdbField.objectArr = append(rdbField.objectArr, object)
}

// DB connect config
// 插入配置
type rdbInsertAction struct {
	rdbField
}

// 创建一个插入设置
func makeRDBAddAction() *rdbInsertAction {
	action := new(rdbInsertAction)
	action.objectArr = make([]interface{}, 0)
	action.fieldArr = make([]field, 0)
	return action
}

// 插入数据配置方法定义
type rdbInsertConfigFunc func(*rdbInsertAction)

// 设置表空间
func InsertSetSpace(db, table string) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.space = space{db: db, table: table}
	}
}

// 添加表数据: key-value 格式
func InsertAddField(keyValueMap map[string]string) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addField(keyValueMap)
	}
}

// 添加表数据: 一整个结构体
func InsertAddObject(object interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addObject(object)
	}
}

// 更新配置
type rdbUpdateAction struct {
	rdbField
	condition UpdateCondition
}

// 创建一个更新设置
func makeRDBUpdateAction() *rdbUpdateAction {
	action := new(rdbUpdateAction)
	action.objectArr = make([]interface{}, 0)
	action.fieldArr = make([]field, 0)
	return action
}

// 插入数据配置方法定义
type rdbUpdateConfigFunc func(*rdbUpdateAction)

// 设置表空间
func UpdateSetSpace(db, table string) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.space = space{db: db, table: table}
	}
}

// 添加表数据: key-value 格式
func UpdateAddField(keyValueMap map[string]string) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.rdbField.addField(keyValueMap)
	}
}

// 添加表数据: 一整个结构体
func UpdateAddObject(object interface{}) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.rdbField.addObject(object)
	}
}

// 设置更新条件
func UpdateSetCondition(condition UpdateCondition) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.condition = condition
	}
}

// 删除配置
type rdbDeleteAction struct {
	space
	condition UpdateCondition
}

// 创建一个删除配置
func makeRDBDeleteAction() *rdbDeleteAction {
	action := new(rdbDeleteAction)
	return action
}

// 删除数据配置方法定义
type rdbDeleteConfigFunc func(*rdbDeleteAction)

// 设置表空间
func DeleteSetSpace(db, table string) func(*rdbDeleteAction) {
	return func(action *rdbDeleteAction) {
		action.space = space{db: db, table: table}
	}
}

// 设置删除条件
func DeleteSetCondition(condition UpdateCondition) func(*rdbDeleteAction) {
	return func(action *rdbDeleteAction) {
		action.condition = condition
	}
}

// 查询配置
type rdbSearchAction struct {
	space
	fields    []string
	condition SearchCondition
}

// 创建一个查询配置
func makeRDBSearchAction() *rdbSearchAction {
	action := new(rdbSearchAction)
	action.fields = make([]string, 0)
	return action
}

// 查询数据配置方法定义
type rdbSearchConfigFunc func(*rdbSearchAction)

// 设置表空间
func SearchSetSpace(db, table string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.space = space{db: db, table: table}
	}
}

// 设置需要查询的字段
func SetSearchFields(fields []string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.fields = fields
	}
}

// 设置查询条件
func SetSearchCondition(condition SearchCondition) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition = condition
	}
}
