package db

import (
	"testing"

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
	localConnector := GetLocalMemoryConnector()
	field := BuildNewField()
	field.AddMap(testKeyValueMap)
	updateRet, _ := localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddField(field))
	log.Info("[TestLocalConnector] update affected rows: %d", updateRet.AffectedRows)
	searchRet, _ := localConnector.Search(SearchSetSpace(testDBName, testTableName))
	log.Info("[TestLocalConnector] search rows: %d", searchRet.Total)
}

func TestLocalFileConnector(t *testing.T) {
	localConnector := GetLocalFileConnector("/tmp/golang")

	// input field
	field := BuildNewField()
	field.AddMap(testKeyValueMap)
	insertRet, err := localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddField(field))
	require.Equal(t, nil, err)
	require.Less(t, 0, insertRet.AffectedRows)

	searchRet, err := localConnector.Search(SearchSetSpace(testDBName, testTableName))
	require.Equal(t, nil, err)
	require.Less(t, 0, searchRet.Total)
	for index, currentField := range searchRet.FieldArr {
		log.Info("[TestLocalFileConnector] search ret: index: %d, field: %s", index, currentField)
	}

	// input struct
	insertRet, err = localConnector.Insert(InsertSetSpace(testDBName, testTableName), InsertAddObject(testObj))
	require.Equal(t, nil, err)
	require.Less(t, 0, insertRet.AffectedRows)

	searchRet, err = localConnector.Search(SearchSetSpace(testDBName, testTableName), SetSearchObject(testStruct{}))
	require.Equal(t, nil, err)
	require.Less(t, 0, searchRet.Total)
	for index, currentObject := range searchRet.ObjectArr {
		log.Info("[TestLocalFileConnector] search ret: index: %d, object: %s", index, currentObject)
	}
}
