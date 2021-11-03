package db

// todo: 设计Connector 接口，包含增删改查等功能
/*
 * Store()
 * Write()
 * Update()
 * Delete()
 */

/*
// 添加学生信息
func (controller *StudentDataController) AddStudent(student *model.Student) int64 {
	createRet := controller.dbSchool.Table(dbutil.StudentTable).Select([]string{"name", "sex"}).Create(student)
	if nil != createRet.Error {
		log.Printf("[AddStudent] 添加学生信息失败，请检查: %s", createRet.Error.Error())
		return 0
	}
	return createRet.RowsAffected
}

func (controller *StudentDataController) UpdateStudentGrade(updateStudent *model.Student) int64 {
	currentStudent := new(model.Student)
	err := controller.dbSchool.Where("id = ?", updateStudent.Id).Select("id").First(&currentStudent).Error
	if nil != err {
		log.Printf("[UpdateStudentGrade] 查询DB错误，请检查: %s", err.Error())
		return 0
	}
	// 这里采用指定字段更新的方式，减少更新成本
	rowsAffected := controller.dbSchool.Model(&currentStudent).Select([]string{"version"}).
		Updates(updateStudent).RowsAffected
	return rowsAffected
}

func (controller *StudentDataController) GetStudentSimpleInfoByName(studentName string) *model.Student {
	currentStudent := new(model.Student)
	err := controller.dbSchool.Where("name = ?", studentName).Select("name, sex").Offset(0).Limit(10).
		First(&currentStudent).Error
	if nil != err {
		log.Printf("[GetStudentSimpleInfoByName] 查询DB错误，请检查: %s", err.Error())
	}
	return currentStudent
}
*/

// Relational data connector
type RDBConnector interface {
	Add(RDBAddConfig) (int, string)
	AddBatch(RDBAddBatchConfig) (int, []string)
	Update(RDBUpdateConfig) int
	Delete(RDBDeleteConfig) int
	Search(RDBSearchConfig) int
}

// DB connect config
// 插入配置
type RDBAddConfig struct {
	TableName string
	Object    interface{}
	Fields    []string
}

// 分批插入配置
type RDBAddBatchConfig struct {
	TableName string
	ObjectArr []interface{}
	Fields    []string
}

// 更新配置
type RDBUpdateConfig struct {
	TableName string
	Object    interface{}
	Fields    []string
	condition UpdateCondition
}

// 删除配置
type RDBDeleteConfig struct {
	TableName string
	Object    interface{}
	condition UpdateCondition
}

// 查询配置
// 删除配置
type RDBSearchConfig struct {
	TableName string
	Object    interface{}
	Fields    []string
	condition SearchCondition
}

// todo: 实现本地内存连接器
func GetLocalMemoryConnector() *RDBConnector {
	return nil
}

// 后续: 初始化 连接器配置中，增加 id generator 配置
