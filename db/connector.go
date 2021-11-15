package db

import (
	"bytes"
	"fmt"
	"reflect"
)

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
	ObjectArr interface{}
	FieldArr  []field
	Page      int
	Len       int
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

// 公共方法: 生成一个新的 field
func BuildNewField() field {
	currentField := field{keyValueMap: make(map[string]string)}
	return currentField
}

// field 中添加单个元素
func (field *field) AddKeyValue(key, value string) {
	field.keyValueMap[key] = value
}

// field 中 批量添加元素
func (field *field) AddMap(keyValueMap map[string]string) {
	for key, value := range keyValueMap {
		field.keyValueMap[key] = value
	}
}

// String
func (field *field) String() string {
	buf := new(bytes.Buffer)
	for key, value := range field.keyValueMap {
		if 0 != buf.Len() {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%s->%s", key, value))
	}
	return buf.String()
}

// 库表属性定义（包括表字段）
type rdbField struct {
	space
	objectArr []interface{}
	// 插入的数组类型
	objectArrType reflect.Type
	fieldArr      []field
	// 需要更新/插入的字段列表，一般在 mysql connector 中使用
	keyArr []string
}

// 添加字段和对应值
func (rdbField *rdbField) addField(field field) {
	rdbField.fieldArr = append(rdbField.fieldArr, field)
}

// 批量添加字段和对应值
func (rdbField *rdbField) addFieldArr(fieldArr []field) {
	rdbField.fieldArr = append(rdbField.fieldArr, fieldArr...)
}

// 添加一整个结构体
func (rdbField *rdbField) addObject(object interface{}) {
	rdbField.objectArr = append(rdbField.objectArr, object)
}

// 批量添加结构体
func (rdbField *rdbField) addObjectArr(objectArr []interface{}) {
	rdbField.objectArr = append(rdbField.objectArr, objectArr...)
}

// 批量添加表字段
func (rdbField *rdbField) addKeyArr(keyArr []string) {
	rdbField.keyArr = append(rdbField.keyArr, keyArr...)
}

// 设置插入的数组类型
func (rdbField *rdbField) setObjectArrType(t reflect.Type) {
	rdbField.objectArrType = t
}

// DB connect config
// 插入配置
type rdbInsertAction struct {
	rdbField
}

// 创建一个插入设置
func makeRDBInsertAction() *rdbInsertAction {
	action := new(rdbInsertAction)
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
func InsertAddField(field field) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addField(field)
	}
}

// 添加表数据: 批量添加 key-value 格式
func InsertAddFieldArr(field []field) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addFieldArr(field)
	}
}

// 添加表数据: 一整个结构体，需要能通过 json 工具类进行解析
func InsertAddObject(object interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addObject(object)
	}
}

// 添加表数据: 批量添加结构体
func InsertAddObjectArr(objectArr []interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addObjectArr(objectArr)
	}
}

// 设置表字段: 设置需要插入的结构体数组类型
func InsertSetObjectArrType(arr interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.setObjectArrType(reflect.TypeOf(arr))
	}
}

// 添加表字段: 设置插入数据需要涉及的表字段列表
func InsertAddKeyArr(keyArr []string) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addKeyArr(keyArr)
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
func UpdateAddField(field field) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.rdbField.addField(field)
	}
}

// 添加表数据: 一整个结构体
func UpdateAddObject(object interface{}) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.rdbField.addObject(object)
	}
}

// 设置更新条件
func UpdateSetCondition(args ...string) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.condition.WhereArr = buildWhereConditionArr(args...)
	}
}

// 添加表字段: 设置更新数据需要涉及的表字段列表
func UpdateAddKeyArr(keyArr []string) func(*rdbUpdateAction) {
	return func(action *rdbUpdateAction) {
		action.rdbField.addKeyArr(keyArr)
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
func DeleteSetCondition(args ...string) func(*rdbDeleteAction) {
	return func(action *rdbDeleteAction) {
		action.condition.WhereArr = buildWhereConditionArr(args...)
	}
}

// 设置删除数据量
func DeleteSetLimit(limit int) func(*rdbDeleteAction) {
	return func(action *rdbDeleteAction) {
		action.condition.Limit = limit
	}
}

// 查询配置
type rdbSearchAction struct {
	space
	keyArr        []string
	object        interface{}  // 用于 format 对象的类型
	objectArrType reflect.Type // 用于生成最后的对象数组的类型
	condition     SearchCondition
}

// 创建一个查询配置
func makeRDBSearchAction() *rdbSearchAction {
	action := new(rdbSearchAction)
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
func SearchSetKeyArr(keyArr []string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.keyArr = append(action.keyArr, keyArr...)
	}
}

// 设置需要查询的结构体
func SearchSetObject(object interface{}) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.object = object
	}
}

// 设置需要返回的结构体数组类型
func SearchSetObjectArrType(arr interface{}) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.objectArrType = reflect.TypeOf(arr)
	}
}

// 设置查询条件
func SearchSetCondition(args ...string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.WhereArr = buildWhereConditionArr(args...)
	}
}

// 设置分页条件
func SearchSetPageCondition(start, limit int) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.Page.No, action.condition.Page.Limit = start/limit, limit
	}
}

// 后续: 设置排序条件
// 最好通过 表名 + 字段名 + 顺序 的方式来设置
