package impala

import (
	"flag"
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

const (
	testDBName    = "db_name"
	testTableName = "table_name"
)

var (
	configPath = flag.String("config", "conf_local.yaml", "config path")
)

func TestImpalaConnector(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(*configPath))
	require.Empty(t, err)
	connector, err := GetImpalaConnector(configManager)
	require.Empty(t, err)

	_, err = connector.Count(db.SearchSetSpace(testDBName, testTableName))
	require.Empty(t, err)
}
