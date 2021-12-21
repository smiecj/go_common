package db

import (
	"testing"

	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

const (
	testMySQLHost     = "test_host"
	testMySQLPort     = 3306
	testMySQLUser     = "test_user"
	testMySQLPassword = "test_password"

	testMySQLDBName    = "temp"
	testMySQLTableName = "test_student"
)

var (
	testStudentArr = []interface{}{
		testStudent{Name: "xiaoming", Grade: 1},
		testStudent{Name: "xiaohong", Grade: 2},
		testStudent{Name: "xiaolin", Grade: 3},
	}
)

// 测试mysql 操作的结构体
/*
CREATE TABLE temp.test_student (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL COMMENT '学生名',
  `grade` int(1) DEFAULT 1 COMMENT '年级',
   PRIMARY KEY (`id`),
  KEY `name_index` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/
type testStudent struct {
	Name  string `gorm:"column:name"`
	Grade int    `gorm:"column:grade"`
}

type studentSlice []testStudent

func TestMySQLConnector(t *testing.T) {
	connector := GetMySQLConnector(
		MySQLConnectOption{Host: testMySQLHost, Port: testMySQLPort, Database: testMySQLDBName, User: testMySQLUser, Password: testMySQLPassword})

	// 插入
	var testStudentSlice studentSlice
	insertRet, err := connector.Insert(InsertSetSpace(testMySQLDBName, testMySQLTableName),
		InsertAddObjectArr(testStudentArr), InsertSetObjectArrType(testStudentSlice))
	require.Equal(t, nil, err)
	require.Equal(t, len(testStudentArr), insertRet.AffectedRows)

	// 查询
	SearchRet, err := connector.Search(SearchSetSpace(testMySQLDBName, testMySQLTableName),
		SearchSetObjectArrType(testStudentSlice), SearchSetPageCondition(0, 10))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, SearchRet.Len)
	studentArr := SearchRet.ObjectArr.(studentSlice)
	log.Info("[TestMySQLConnector] object arr len: %d, total: %d", len(studentArr), SearchRet.Total)
	for _, currentStudent := range studentArr {
		log.Info("[TestMySQLConnector] current student: %v", currentStudent)
	}

	// distinct
	SearchRet, err = connector.Distinct(SearchSetSpace(testMySQLDBName, testMySQLTableName),
		SearchSetKeyArr([]string{"name", "grade"}))
	require.Equal(t, nil, err)
	log.Info("[TestMySQLConnector] distinct len: %d", SearchRet.Len)
	for _, currentField := range SearchRet.FieldArr {
		log.Info("[TestMySQLConnector] Distinct name result: %s", currentField)
	}

	// 更新
	UpdateRet, err := connector.Update(UpdateSetSpace(testMySQLDBName, testMySQLTableName),
		UpdateSetCondition("name", "=", "xiaoming"),
		UpdateAddObject(testStudent{Grade: 2}), UpdateAddKeyArr([]string{"grade"}))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, UpdateRet.AffectedRows)

	// 删除
	deleteRet, err := connector.Delete(DeleteSetSpace(testMySQLDBName, testMySQLTableName),
		DeleteSetCondition("name", "=", "xiaoming"))
	require.Equal(t, nil, err)
	require.LessOrEqual(t, 1, deleteRet.AffectedRows)
}
