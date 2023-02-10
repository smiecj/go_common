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
	testEnvName  = "test_env"
)

// 测试占锁功能
func TestLock(t *testing.T) {
	configManager, err := config.GetYamlConfigManager("/tmp/conf.yaml")
	require.Empty(t, err)
	connector, err := mysql.GetMySQLConnector(configManager)
	require.Empty(t, err)

	errorChan := make(chan error)
	lockManager := GetLockManager(connector, errorChan)
	// 测试用: 修改环境名，确认是否能成功占锁（正常情况下不能马上抢占）
	lockManager.envName = testEnvName
	dataChan := make(chan int)
	putDataToChan := func(ctx context.Context) error {
		dataChan <- 1
		return nil
	}
	lockManager.Lock(testLockName, IntervalShort, putDataToChan)
	<-dataChan

	// 后续: 可以补充更完整的测试用例，包括全部场景
}
