package mysql

import (
	"flag"
	"fmt"
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

const (
	dbTemp          = "temp"
	tableClass      = "test_class"
	tableStudent    = "test_student"
	tableStudentBak = "test_student_bak"
)

var (
	testStudentArr = []testStudent{
		{Name: "xiaoming", ClassId: 1},
		{Name: "xiaohong", ClassId: 2},
		{Name: "xiaolin", ClassId: 2},
	}
	testStudentSingle = testStudent{Name: "xiaozhang", ClassId: 3}

	anotherSchoolStudentName = "xiaobai"

	configPath = flag.String("config", "conf_local.yaml", "config path")
)

// test mysql struct
type testStudent struct {
	Name    string `gorm:"column:name"`
	ClassId int    `gorm:"column:class_id"`
}

type studentSlice []testStudent

func (slice *studentSlice) getFields() []string {
	return []string{"name", "class_id"}
}

type testStudentWithClass struct {
	Name      string `gorm:"column:name"`
	ClassId   int    `gorm:"column:class_id"`
	ClassName string `gorm:"column:class_name"` // 测试 join
}

type testStudentWithClassSlice []testStudentWithClass

func (slice *testStudentWithClassSlice) getFields() []string {
	return []string{"test_student.name", "test_student.class_id", "test_class.name AS class_name"}
}

// mysql db 连接器完整测试
func TestMySQLConnector(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(*configPath))
	require.Empty(t, err)
	connector, err := GetMySQLConnector(configManager)
	require.Empty(t, err)

	// insert
	var testStudentSlice studentSlice
	// insert batch
	insertRet, err := connector.Insert(InsertSetSpace(dbTemp, tableStudent),
		InsertAddObjectArr(testStudentArr),
		InsertAddKeyArr(testStudentSlice.getFields()))
	require.Equal(t, nil, err)
	require.Equal(t, len(testStudentArr), insertRet.AffectedRows)
	// insert single
	insertRet, err = connector.Insert(InsertSetSpace(dbTemp, tableStudent),
		InsertSetObject(testStudentSingle), InsertAddKeyArr(testStudentSlice.getFields()))
	require.Equal(t, nil, err)
	require.Equal(t, 1, insertRet.AffectedRows)

	// search
	searchRet, err := connector.Search(SearchSetSpace(dbTemp, tableStudent),
		SearchSetObjectArrType(testStudentSlice), SearchSetPageCondition(0, 10),
		SearchSetKeyArr(testStudentSlice.getFields()))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, searchRet.Len)
	searchStudentArrRet := searchRet.ObjectArr.(studentSlice)
	require.GreaterOrEqual(t, len(searchStudentArrRet), 3)
	require.Equal(t, len(searchStudentArrRet), searchRet.Len)
	require.GreaterOrEqual(t, searchRet.Total, searchRet.Len)
	require.NotEmpty(t, searchStudentArrRet[0].Name)
	// search not exist result
	searchRet, err = connector.Search(SearchSetSpace(dbTemp, tableStudent),
		SearchSetObjectArrType(testStudentSlice), SearchSetPageCondition(0, 10),
		SearchSetCondition("name", "=", anotherSchoolStudentName),
		SearchSetKeyArr(testStudentSlice.getFields()))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 0, searchRet.Len)
	_, isConvertSuccess := searchRet.ObjectArr.(studentSlice)
	require.True(t, isConvertSuccess)

	// search min/max
	searchRet, err = connector.Search(SearchSetSpace(dbTemp, tableStudent),
		SearchSetKeyArr([]string{"max(class_id)"}))
	require.Empty(t, err)
	require.Equal(t, 1, len(searchRet.FieldArr))
	require.Equal(t, 1, len(searchRet.FieldArr[0].GetMap()))

	// search with join
	var testStudentClassSlice testStudentWithClassSlice
	searchRet, err = connector.Search(SearchSetSpace(dbTemp, tableStudent),
		SearchSetObjectArrType(testStudentClassSlice),
		SearchSetCondition(fmt.Sprintf("%s.%s", tableStudent, "name"), "=", testStudentSingle.Name),
		SearchAddJoin(dbTemp, tableStudent, "class_id", tableClass, "id"),
		SearchSetKeyArr(testStudentClassSlice.getFields()))
	require.Empty(t, err)
	require.LessOrEqual(t, 1, searchRet.Len)

	// search with group by
	// select class_id, count(*) from test_class left join test_student on ... group by test_class.id
	searchRet, err = connector.Search(SearchSetSpace(dbTemp, tableStudent),
		SearchSetObjectArrType(testStudentClassSlice),
		SearchSetCondition(fmt.Sprintf("%s.%s", tableStudent, "name"), "=", testStudentSingle.Name),
		SearchAddJoin(dbTemp, tableStudent, "class_id", tableClass, "id"),
		SearchSetKeyArr(testStudentClassSlice.getFields()))

	// distinct
	searchRet, err = connector.Distinct(SearchSetSpace(dbTemp, tableStudent),
		SearchSetKeyArr([]string{"name", "class_id"}))
	require.Equal(t, nil, err)
	require.GreaterOrEqual(t, searchRet.Len, 3)

	// update
	UpdateRet, err := connector.Update(UpdateSetSpace(dbTemp, tableStudent),
		UpdateSetCondition("name", "=", "xiaoming"),
		UpdateAddObject(testStudent{ClassId: 2}), UpdateAddKeyArr([]string{"class_id"}))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, UpdateRet.AffectedRows)

	// backup
	backupRet, err := connector.Backup(BackupSetSourceSpace(dbTemp, tableStudent),
		BackupSetTargetSpace(dbTemp, tableStudentBak),
		BackupSetCondition("name", "in", "('xiaoming', 'xiaohong', 'xiaolin', 'xiaozhang')"))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, backupRet.AffectedRows)

	// delete all
	deleteRet, err := connector.Delete(DeleteSetSpace(dbTemp, tableStudent),
		DeleteSetCondition("name", "in", "('xiaoming', 'xiaohong', 'xiaolin', 'xiaozhang')"))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, deleteRet.AffectedRows)
	deleteRet, err = connector.Delete(DeleteSetSpace(dbTemp, tableStudentBak),
		DeleteSetLimit(backupRet.AffectedRows))
	require.Equal(t, nil, err)
	require.Equal(t, backupRet.AffectedRows, deleteRet.AffectedRows)

	// close
	err = connector.Close()
	require.Nil(t, err)
}
