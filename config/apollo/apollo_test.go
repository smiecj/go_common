package config

import (
	"flag"
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

var (
	configPath = flag.String("config", "conf_local.yaml", "config path")
)

func TestApolloConfig(t *testing.T) {
	yamlConfigManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(*configPath))
	require.Nil(t, err)
	apolloConfigManager, err := GetApolloConfigManager(yamlConfigManager)
	require.Nil(t, err)

	testValue, err := apolloConfigManager.Get("test_group", "test_key")
	require.Nil(t, err)
	require.Equal(t, "test_value", testValue)
}
