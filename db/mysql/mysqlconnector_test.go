package mysql

import (
	"testing"

	"github.com/smiecj/go_common/config"
	. "github.com/smiecj/go_common/db"
	"github.com/stretchr/testify/require"
)

const (
	testMySQLDBName    = "temp"
	testMySQLTableName = "test_student"
)

var (
	testStudentArr = []interface{}{
		testStudent{Name: "xiaoming", Grade: 1},
		testStudent{Name: "xiaohong", Grade: 2},
		testStudent{Name: "xiaolin", Grade: 3},
	}
	testStudentSingle = testStudent{Name: "xiaozhang", Grade: 2}
)

// 测试mysql 操作的结构体
/*
CREATE TABLE `temp`.`test_student` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL COMMENT '学生名',
  `grade` int(1) DEFAULT 1 COMMENT '年级',
   PRIMARY KEY (`id`),
  KEY `name_index` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
INSERT INTO temp.test_student(name, grade) VALUES('xiaoming', 1), ('xiaohong', 2), ('xiaolin', 3);
*/
type testStudent struct {
	Name  string `gorm:"column:name"`
	Grade int    `gorm:"column:grade"`
}

type studentSlice []testStudent

func TestMySQLConnector(t *testing.T) {
	configManager, err := config.GetYamlConfigManager("/tmp/conf.yaml")
	require.Empty(t, err)
	connector, err := GetMySQLConnector(configManager)
	require.Empty(t, err)

	// insert
	var testStudentSlice studentSlice
	// insert batch
	insertRet, err := connector.Insert(InsertSetSpace(testMySQLDBName, testMySQLTableName),
		InsertAddObjectArr(testStudentArr), InsertSetObjectArrType(testStudentSlice))
	require.Equal(t, nil, err)
	require.Equal(t, len(testStudentArr), insertRet.AffectedRows)
	// insert single
	insertRet, err = connector.Insert(InsertSetSpace(testMySQLDBName, testMySQLTableName),
		InsertSetObject(testStudentSingle))
	require.Equal(t, nil, err)
	require.Equal(t, 1, insertRet.AffectedRows)

	// search
	searchRet, err := connector.Search(SearchSetSpace(testMySQLDBName, testMySQLTableName),
		SearchSetObjectArrType(testStudentSlice), SearchSetPageCondition(0, 10))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, searchRet.Len)
	studentArr := searchRet.ObjectArr.(studentSlice)
	require.GreaterOrEqual(t, len(studentArr), 3)
	require.Equal(t, len(studentArr), searchRet.Len)
	require.GreaterOrEqual(t, searchRet.Total, searchRet.Len)

	// search min/max
	searchRet, err = connector.Search(SearchSetSpace(testMySQLDBName, testMySQLTableName),
		SearchSetKeyArr([]string{"max(grade)"}))
	require.Empty(t, err)
	require.Equal(t, 1, len(searchRet.FieldArr))
	require.Equal(t, 1, len(searchRet.FieldArr[0].GetMap()))

	// distinct
	searchRet, err = connector.Distinct(SearchSetSpace(testMySQLDBName, testMySQLTableName),
		SearchSetKeyArr([]string{"name", "grade"}))
	require.Equal(t, nil, err)
	require.GreaterOrEqual(t, searchRet.Len, 3)

	// update
	UpdateRet, err := connector.Update(UpdateSetSpace(testMySQLDBName, testMySQLTableName),
		UpdateSetCondition("name", "=", "xiaoming"),
		UpdateAddObject(testStudent{Grade: 2}), UpdateAddKeyArr([]string{"grade"}))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, UpdateRet.AffectedRows)

	// delete
	deleteRet, err := connector.Delete(DeleteSetSpace(testMySQLDBName, testMySQLTableName),
		DeleteSetCondition("name", "in", "('xiaoming', 'xiaohong', 'xiaolin', 'xiaozhang')"), DeleteSetLimit(insertRet.AffectedRows))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, deleteRet.AffectedRows)
}
