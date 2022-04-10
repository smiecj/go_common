package clusterlock

import (
	"context"
	"testing"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/db/mysql"
	"github.com/stretchr/testify/require"
)

const (
	testLockName = "test_lock"
)

func TestLock(t *testing.T) {
	configManager, err := config.GetYamlConfigManager("/tmp/conf.yaml")
	require.Empty(t, err)
	connector, err := mysql.GetMySQLConnector(configManager)
	require.Empty(t, err)

	errorChan := make(chan error)
	lockManager := GetLockManager(connector, errorChan)
	dataChan := make(chan int)
	putDataToChan := func(ctx context.Context) error {
		dataChan <- 1
		return nil
	}
	lockManager.Lock(testLockName, IntervalShort, putDataToChan)
	<-dataChan
}
