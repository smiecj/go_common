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

func TestNacosConfig(t *testing.T) {
	yamlConfigManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(*configPath))
	require.Nil(t, err)
	nacosConfigManager, err := GetNacosConfigManager(yamlConfigManager)
	require.Nil(t, err)

	testValue, err := nacosConfigManager.Get("test_group", "test_key")
	require.Nil(t, err)
	require.Equal(t, "test_value", testValue)
}
