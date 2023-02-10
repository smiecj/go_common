package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testSpaceDB        = "db"
	testKeyMySQLHost   = "mysql_host"
	testKeyMySQLPort   = "mysql_port"
	testKeyMySQLDBArr  = "db_arr"
	testValueMySQLHost = "localhost"
	testValueMySQLPort = 3306
	testConfigPath     = "/tmp"
)

type dbConfig struct {
	MysqlHost string   `yaml:"mysql_host"`
	MysqlPort int      `yaml:"mysql_port"`
	DBArr     []string `yaml:"db_arr"`
}

func TestYamlConfig(t *testing.T) {
	// get config manager
	config, err := GetYamlConfigManager(fmt.Sprintf("%s%sexample.yaml", testConfigPath, string(os.PathSeparator)))
	require.Equal(t, nil, err)

	// get config
	host, err := config.Get(testSpaceDB, testKeyMySQLHost)
	require.Empty(t, err)
	port, err := config.Get(testSpaceDB, testKeyMySQLPort)
	require.Empty(t, err)
	require.Equal(t, testValueMySQLHost, host)
	require.Equal(t, testValueMySQLPort, port)
	spaceNameArr, err := config.GetAllSpaceName()
	require.Empty(t, err)
	require.NotEmpty(t, spaceNameArr)

	// unmarshal
	dbConfigObj := dbConfig{}
	_ = config.Unmarshal(testSpaceDB, &dbConfigObj)
	require.Equal(t, testValueMySQLHost, dbConfigObj.MysqlHost)
	require.Equal(t, 1, len(dbConfigObj.DBArr))
	require.NotEqual(t, "", dbConfigObj.DBArr[0])

	// get space & set config
	space, err := config.GetSpace(testSpaceDB)
	require.Empty(t, err)
	host, err = space.Get(testKeyMySQLHost)
	require.Empty(t, err)
	require.Equal(t, testValueMySQLHost, host)
	keyArr, err := space.GetAllKey()
	require.Empty(t, err)
	require.NotEmpty(t, keyArr)

	// unmarshal
	_ = space.Unmarshal(&dbConfigObj)
	require.Equal(t, testValueMySQLHost, dbConfigObj.MysqlHost)
	require.Equal(t, testValueMySQLPort, dbConfigObj.MysqlPort)
}
