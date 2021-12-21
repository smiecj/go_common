package impala

import (
	"testing"

	"github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

const (
	testImpalaHost = "impala_host"
	testImpalaPort = 21050

	testDBName    = "db_name"
	testTableName = "table_name"
)

func TestImpalaConnector(t *testing.T) {
	connector := GetImpalaConnector(ImpalaConnectOption{Host: testImpalaHost, Port: testImpalaPort})
	require.NotEmpty(t, connector)

	ret, err := connector.Count(db.SearchSetSpace(testDBName, testTableName))
	require.Empty(t, err)
	require.Less(t, 0, ret.Total)
	log.Info("[test] total: %d", ret.Total)
}
