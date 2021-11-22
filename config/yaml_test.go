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
	testValueMySQLHost = "localhost"
	testValueMySQLPort = 3306
	testConfigPath     = "/tmp"
)

type dbConfig struct {
	MysqlHost string `json:"mysql_host"`
	MysqlPort int    `json:"mysql_port"`
}

func TestYamlConfig(t *testing.T) {
	// get config manager
	config, err := GetYamlConfig(fmt.Sprintf("%s%sexample.yaml", testConfigPath, string(os.PathSeparator)))
	require.Equal(t, nil, err)

	// get config
	host, err := config.Get(testSpaceDB, testKeyMySQLHost)
	require.Equal(t, nil, err)
	port, err := config.Get(testSpaceDB, testKeyMySQLPort)
	require.Equal(t, nil, err)
	require.Equal(t, testValueMySQLHost, host)
	require.Equal(t, testValueMySQLPort, port)

	// unmarshal
	dbConfigObj := dbConfig{}
	_ = config.Unmarshal(testSpaceDB, &dbConfigObj)
	require.Equal(t, testValueMySQLHost, dbConfigObj.MysqlHost)

	// get space & set config
	space, err := config.GetSpace(testSpaceDB)
	require.Equal(t, nil, err)
	host, err = space.Get(testKeyMySQLHost)
	require.Equal(t, nil, err)
	require.Equal(t, testValueMySQLHost, host)
	_ = space.Unmarshal(&dbConfigObj)
	require.Equal(t, testValueMySQLHost, dbConfigObj.MysqlHost)
	require.Equal(t, testValueMySQLPort, dbConfigObj.MysqlPort)
}
