package impala

import (
	"flag"
	"testing"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/db"
	"github.com/stretchr/testify/require"
)

const (
	testDBName    = "db_name"
	testTableName = "table_name"
)

var (
	configPath = flag.String("config", "/tmp/conf.yaml", "config path")
)

func TestImpalaConnector(t *testing.T) {
	configManager, err := config.GetYamlConfigManager(*configPath)
	require.Empty(t, err)
	connector, err := GetImpalaConnector(configManager)
	require.Empty(t, err)

	ret, err := connector.Count(db.SearchSetSpace(testDBName, testTableName))
	require.Empty(t, err)
	require.Less(t, 0, ret.Total)
}
