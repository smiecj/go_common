package clusterlock

import (
	"context"
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/db/mysql"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

const (
	testLockName    = "test_lock"
	testEnvName     = "test_env"
	localConfigFile = "conf_local.yaml"
)

// 测试占锁功能
func TestLock(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(localConfigFile))
	require.Empty(t, err)
	connector, err := mysql.GetMySQLConnector(configManager)
	require.Empty(t, err)

	errorChan := make(chan error)
	lockManager := GetLockManager(connector, errorChan)
	lockManager.envName = testEnvName
	dataChan := make(chan int)
	putDataToChan := func(ctx context.Context) error {
		dataChan <- 1
		return nil
	}
	// 占锁
	lockManager.Lock(testLockName, IntervalShort, putDataToChan)
	<-dataChan
}
