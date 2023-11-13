package impala

import (
	. "github.com/smiecj/go_common/db"
)

// impala 连接器
type mockImpalaConnector struct{}

func (connector *mockImpalaConnector) Insert(funcArr ...RDBInsertConfigFunc) (ret UpdateRet, err error) {
	return UpdateRet{}, nil
}

func (connector *mockImpalaConnector) Update(funcArr ...RDBUpdateConfigFunc) (ret UpdateRet, err error) {
	return UpdateRet{}, nil
}

func (connector *mockImpalaConnector) Delete(funcArr ...RDBDeleteConfigFunc) (ret UpdateRet, err error) {
	return UpdateRet{}, nil
}

func (connector *mockImpalaConnector) Backup(funcArr ...RDBBackupConfigFunc) (ret UpdateRet, err error) {
	return UpdateRet{}, nil
}

func (connector *mockImpalaConnector) Search(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return SearchRet{}, nil
}

func (connector *mockImpalaConnector) ExecSearch(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return SearchRet{}, nil
}

func (connector *mockImpalaConnector) Count(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return SearchRet{}, nil
}

func (connector *mockImpalaConnector) Distinct(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return SearchRet{}, nil
}

func (connector *mockImpalaConnector) Close() error {
	return nil
}

func (connector *mockImpalaConnector) Stat() (ret DBStat, err error) {
	return DBStat{}, nil
}
