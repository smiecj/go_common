package yaml

import (
	"testing"

	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

const (
	testSpaceMySQL     = "mysql"
	testKeyMySQLHost   = "host"
	testKeyMySQLPort   = "port"
	testValueMySQLHost = "localhost"
	testValueMySQLPort = 3306
	exampleConfigFile  = "conf_example.yaml"
)

type mysqlConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func TestYamlConfig(t *testing.T) {
	// get config manager
	config, err := GetYamlConfigManager(file.FindFilePath(exampleConfigFile))
	require.Equal(t, nil, err)

	// get config
	host, err := config.Get(testSpaceMySQL, testKeyMySQLHost)
	require.Empty(t, err)
	port, err := config.Get(testSpaceMySQL, testKeyMySQLPort)
	require.Empty(t, err)
	require.Equal(t, testValueMySQLHost, host)
	require.Equal(t, testValueMySQLPort, port)
	spaceNameArr, err := config.GetAllSpaceName()
	require.Empty(t, err)
	require.NotEmpty(t, spaceNameArr)

	// unmarshal
	mysqlConfigObj := mysqlConfig{}
	_ = config.Unmarshal(testSpaceMySQL, &mysqlConfigObj)
	require.Equal(t, testValueMySQLHost, mysqlConfigObj.Host)
	require.Equal(t, testValueMySQLPort, mysqlConfigObj.Port)

	// get space & set config
	mysqlSpace, err := config.GetSpace(testSpaceMySQL)
	require.Empty(t, err)
	host, err = mysqlSpace.Get(testKeyMySQLHost)
	require.Empty(t, err)
	require.Equal(t, testValueMySQLHost, host)
	keyArr, err := mysqlSpace.GetAllKey()
	require.Empty(t, err)
	require.NotEmpty(t, keyArr)

	// unmarshal
	_ = mysqlSpace.Unmarshal(&mysqlConfigObj)
	require.Equal(t, testValueMySQLHost, mysqlConfigObj.Host)
	require.Equal(t, testValueMySQLPort, mysqlConfigObj.Port)
}
