package db

import (
	"bytes"
	"fmt"
	"reflect"
)

const (
	// 变更操作 最大 limit
	maxModifyLimit = 100000
)

// Relational data connector
type RDBConnector interface {
	Insert(...RDBInsertConfigFunc) (UpdateRet, error)
	Update(...RDBUpdateConfigFunc) (UpdateRet, error)
	Delete(...RDBDeleteConfigFunc) (UpdateRet, error)
	Backup(...RDBBackupConfigFunc) (UpdateRet, error)
	Search(...RDBSearchConfigFunc) (SearchRet, error)
	Count(...RDBSearchConfigFunc) (SearchRet, error)
	Distinct(...RDBSearchConfigFunc) (SearchRet, error)
	Stat() (DBStat, error)
	Close() error
}

// db connect status
// refer: golang/src/database/sql/sql.go DBStats
type DBStat struct {
	// pool status
	OpenConnections int
	InUse           int
	Idle            int
}

// 更新类型动作结果
type UpdateRet struct {
	AffectedRows int
}

// 查询类型动作结果
type SearchRet struct {
	ObjectArr interface{}
	FieldArr  []field
	Page      int
	Len       int
	Total     int
}

// SearchRet: 添加字段和对应值
func (SearchRet *SearchRet) AddField(field field) {
	SearchRet.FieldArr = append(SearchRet.FieldArr, field)
}

// 库表空间定义
type space struct {
	db    string
	table string
}

// 获取库名.表名的格式
func (space *space) GetSpaceName() string {
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

// 获取内部 map 对象
func (field *field) GetMap() map[string]string {
	return field.keyValueMap
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
	// 插入/更新单条数据
	object interface{}
	// 批量插入/更新数据
	objectArr []interface{}
	// 插入的数组类型
	objectArrType reflect.Type
	fieldArr      []field
	// 需要更新/插入的字段列表，一般在 mysql connector 中使用
	keyArr []string
}

// 设置单个 object
func (rdbField *rdbField) setObject(object interface{}) {
	rdbField.object = object
}

// 获取单个 object
func (rdbField *rdbField) GetObject() interface{} {
	return rdbField.object
}

// 获取所有的 field 字段
func (rdbField *rdbField) GetFieldArr() []field {
	return rdbField.fieldArr
}

// 获取 object arr
func (rdbField *rdbField) GetObjectArr() []interface{} {
	return rdbField.objectArr
}

// 获取数组类型
func (rdbField *rdbField) GetObjectArrType() reflect.Type {
	return rdbField.objectArrType
}

// 获取需要更新的字段列表
func (rdbField *rdbField) GetKeyArr() []string {
	return rdbField.keyArr
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
func MakeRDBInsertAction() *rdbInsertAction {
	action := new(rdbInsertAction)
	return action
}

// 插入数据配置方法定义
type RDBInsertConfigFunc func(*rdbInsertAction)

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

// 添加表数据: 添加单个结构体
func InsertSetObject(object interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.setObject(object)
	}
}

// 添加表数据: 一整个结构体，需要能通过 json 工具类进行解析
func InsertAddObject(object interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		action.rdbField.addObject(object)
	}
}

// 添加表数据: 批量添加结构体
func InsertAddObjectArr(objectArr interface{}) func(*rdbInsertAction) {
	return func(action *rdbInsertAction) {
		t := reflect.TypeOf(objectArr)
		if t.Kind() != reflect.Slice {
			return
		}
		objectArrVal := reflect.ValueOf(objectArr)
		toAddObjectArr := make([]interface{}, objectArrVal.Len())
		for index := 0; index < objectArrVal.Len(); index++ {
			toAddObjectArr[index] = objectArrVal.Index(index).Interface()
		}
		action.rdbField.addObjectArr(toAddObjectArr)
		action.setObjectArrType(reflect.TypeOf(objectArr))
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
	condition updateCondition
}

// 创建一个更新设置
func MakeRDBUpdateAction() *rdbUpdateAction {
	action := new(rdbUpdateAction)
	return action
}

// 获取更新条件
func (action *rdbUpdateAction) GetCondition() updateCondition {
	return action.condition
}

// 插入数据配置方法定义
type RDBUpdateConfigFunc func(*rdbUpdateAction)

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
	condition updateCondition
}

// 获取删除条件
func (action *rdbDeleteAction) GetCondition() updateCondition {
	return action.condition
}

// 创建一个删除配置
// DB 保护: 默认最多查询 1kw 条数据
func MakeRDBDeleteAction() *rdbDeleteAction {
	action := new(rdbDeleteAction)
	action.condition.Limit = maxModifyLimit
	return action
}

// 删除数据配置方法定义
type RDBDeleteConfigFunc func(*rdbDeleteAction)

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

// 获取需要查询的字段列表
func (action *rdbSearchAction) GetKeyArr() []string {
	return action.keyArr
}

// 获取查询条件
func (action *rdbSearchAction) GetCondition() SearchCondition {
	return action.condition
}

// 获取用于format 的对象类型
func (action *rdbSearchAction) GetObject() interface{} {
	return action.object
}

// 获取查询条件
func (action *rdbSearchAction) GetObjectArrType() reflect.Type {
	return action.objectArrType
}

// 创建一个查询配置
// DB 保护: 默认最多查询 1kw 条数据
func MakeRDBSearchAction() *rdbSearchAction {
	action := new(rdbSearchAction)
	action.condition.Page.Limit = maxModifyLimit
	return action
}

// 查询数据配置方法定义
type RDBSearchConfigFunc func(*rdbSearchAction)

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

// 设置 limit 条件
func SearchSetLimit(limit int) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.Page.Limit = limit
	}
}

// 设置分页条件
func SearchSetPageCondition(start, limit int) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.Page.No, action.condition.Page.Limit = start/limit, limit
	}
}

// 设置排序条件 (仅排序字段 默认 asc)
func SearchSetOrderField(field string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.Order.Field = field
	}
}

// 设置排序条件（排序字段 + 排序顺序，asc or desc）
func SearchSetOrderFieldAndAsc(field, asc string) func(*rdbSearchAction) {
	return func(action *rdbSearchAction) {
		action.condition.Order.Field, action.condition.Order.Sc = field, asc
	}
}

// 设置 join 条件, 默认 join 表 为 查询条件中的右侧表
func SearchAddJoin(fieldArr ...string) RDBSearchConfigFunc {
	return func(rsa *rdbSearchAction) {
		// at least contain: db, left table, left field, right table, right field
		// can also provide: join method (default left join)
		if len(fieldArr) < 5 {
			return
		}
		toAppendJoinCondition := joinCondition{}
		db, leftTable, leftField, rightTable, rightField := fieldArr[0], fieldArr[1], fieldArr[2], fieldArr[3], fieldArr[4]
		toAppendJoinCondition.joinMethod = LeftJoin
		if len(fieldArr) == 6 {
			toAppendJoinCondition.joinMethod = JoinMethod(fieldArr[6])
		}
		toAppendJoinCondition.space = space{db: db, table: rightTable}
		// condition 的组装，一般左右条件都是字段名，暂不考虑更复杂的情况（多条件，或是值判断）

		// condition: 去掉多余空格，并按空格拆分成多个条件
		toAppendJoinCondition.condition = []string{
			// 加上 ` 表标识，防止当成字段
			fmt.Sprintf("`%s`.%s", leftTable, leftField),
			"=",
			fmt.Sprintf("`%s`.%s", rightTable, rightField),
		}

		rsa.condition.Join = append(rsa.condition.Join, toAppendJoinCondition)
	}
}

// 备份配置
type rdbBackupAction struct {
	sourceSpace *space
	condition   updateCondition
	targetSpace *space
}

// 获取更新条件
func (action *rdbBackupAction) GetCondition() updateCondition {
	return action.condition
}

// 获取源表名称
func (action *rdbBackupAction) GetSourceSpaceName() string {
	return action.sourceSpace.GetSpaceName()
}

// 获取目标表名称
func (action *rdbBackupAction) GetTargetSpaceName() string {
	return action.targetSpace.GetSpaceName()
}

// 创建一个备份配置
// DB 保护: 默认最多备份 1kw 条数据
func MakeRDBBackupAction() *rdbBackupAction {
	action := new(rdbBackupAction)
	action.condition.Limit = maxModifyLimit
	return action
}

// 备份数据配置方法定义
type RDBBackupConfigFunc func(*rdbBackupAction)

// 设置备份数据源表空间
func BackupSetSourceSpace(db, table string) func(*rdbBackupAction) {
	return func(action *rdbBackupAction) {
		action.sourceSpace = &space{db: db, table: table}
	}
}

// 设置备份条件
func BackupSetCondition(args ...string) func(*rdbBackupAction) {
	return func(action *rdbBackupAction) {
		action.condition.WhereArr = buildWhereConditionArr(args...)
	}
}

// 设置备份数据量
func BackupSetLimit(limit int) func(*rdbBackupAction) {
	return func(action *rdbBackupAction) {
		action.condition.Limit = limit
	}
}

// 设置备份目标表
func BackupSetTargetSpace(db, table string) func(*rdbBackupAction) {
	return func(action *rdbBackupAction) {
		action.targetSpace = &space{db: db, table: table}
	}
}

// 后续: 设置排序条件
// 最好通过 表名 + 字段名 + 顺序 的方式来设置
