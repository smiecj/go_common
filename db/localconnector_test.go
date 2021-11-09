package db

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
)

const (
	dbName    = "test_db"
	tableName = "test_table"
)

func TestLocalConnector(t *testing.T) {
	localConnector := GetLocalMemoryConnector()
	dataMap := map[string]string{
		"user":    "smiecj",
		"country": "China",
	}
	updateRet, _ := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(dataMap))
	log.Info("[TestLocalConnector] update affected rows: %d", updateRet.AffectedRows)
	searchRet, _ := localConnector.Search(SearchSetSpace(dbName, tableName))
	log.Info("[TestLocalConnector] search rows: %d", searchRet.Total)
}
