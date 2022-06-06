package impala

import (
	"testing"

	"github.com/smiecj/go_common/db"
	"github.com/stretchr/testify/require"
)

const (
	testImpalaHost = "impala_host"
	testImpalaPort = 21050

	testDBName    = "db_name"
	testTableName = "table_name"
)

func TestImpalaConnector(t *testing.T) {
	connector, err := GetImpalaConnector(ImpalaConnectOption{Host: testImpalaHost, Port: testImpalaPort})
	require.Empty(t, err)

	ret, err := connector.Count(db.SearchSetSpace(testDBName, testTableName))
	require.Empty(t, err)
	require.Less(t, 0, ret.Total)
}
