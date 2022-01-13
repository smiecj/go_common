package local

import (
	"testing"

	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	User    string `json:"user"`
	Country string `json:"country"`
}

const (
	testDBName    = "test_db"
	testTableName = "test_table"
)

var (
	testKeyValueMap = map[string]string{
		"user":    "smiecj",
		"country": "China",
	}
	testObj = testStruct{User: "smiecj", Country: "China"}
)

func TestLocalMemoryConnector(t *testing.T) {
	localConnector, _ := GetLocalMemoryConnector()
	field := BuildNewField()
	field.AddMap(testKeyValueMap)
	UpdateRet, _ := localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddField(field))
	log.Info("[TestLocalConnector] update affected rows: %d", UpdateRet.AffectedRows)
	SearchRet, _ := localConnector.Search(SearchSetSpace(testDBName, testTableName))
	log.Info("[TestLocalConnector] search rows: %d", SearchRet.Len)
}

func TestLocalFileConnector(t *testing.T) {
	localConnector, err := GetLocalFileConnector("/tmp/golang")
	require.Empty(t, err)

	// input field
	field := BuildNewField()
	field.AddMap(testKeyValueMap)
	insertRet, err := localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddField(field))
	require.Equal(t, nil, err)
	require.Less(t, 0, insertRet.AffectedRows)

	SearchRet, err := localConnector.Search(SearchSetSpace(testDBName, testTableName))
	require.Equal(t, nil, err)
	require.Less(t, 0, SearchRet.Len)
	for index, currentField := range SearchRet.FieldArr {
		log.Info("[TestLocalFileConnector] search ret: index: %d, field: %s", index, currentField)
	}

	// input struct
	insertRet, err = localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddObject(testObj))
	require.Equal(t, nil, err)
	require.Less(t, 0, insertRet.AffectedRows)

	SearchRet, err = localConnector.Search(SearchSetSpace(testDBName, testTableName),
		SearchSetObject(testStruct{}), SearchSetObjectArrType([]*testStruct{}))
	require.Equal(t, nil, err)
	require.Less(t, 0, SearchRet.Len)
	testStructArr := SearchRet.ObjectArr.([]*testStruct)
	for index, currentStruct := range testStructArr {
		log.Info("[TestLocalFileConnector] search ret: index: %d, object: %s", index, currentStruct)
	}
}
