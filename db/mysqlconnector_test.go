package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testMySQLHost     = "localhost"
	testMySQLPort     = 23306
	testMySQLUser     = "root"
	testMySQLPassword = "root123"

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

type testStudent struct {
	Name  string `gorm:"column:name"`
	Grade int    `gorm:"column:grade"`
}

func TestMySQLConnector(t *testing.T) {
	// todo: 增删改查测试
	connector := GetMySQLConnector(
		MySQLConnectOption{Host: "localhost", Port: 23306, User: testMySQLUser, Password: testMySQLPassword})

	ret, err := connector.Insert(InsertSetSpace(testMySQLDBName, testTableName), InsertAddObjectArr(testStudentArr))
	require.Equal(t, nil, err)
	require.Equal(t, len(testStudentArr), ret.AffectedRows)
}
