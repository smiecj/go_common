package monitor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type managerTestCase struct {
	metricsName         string
	desc                monitorMetricsDesc
	setMetricsValueFunc func(metrics metrics)
}

var (
	// table driven test
	testCaseArr = []managerTestCase{
		{
			metricsName:         "test_gauge",
			desc:                NewMonitorMetrics(Gauge, "test_gauge", "test gauge", LabelKey{"name"}),
			setMetricsValueFunc: func(metrics metrics) { metrics.(*PrometheusGauge).With(MetricsLabel{"name": "smiecj"}).Set(10) },
		},
		{
			metricsName:         "test_counter",
			desc:                NewMonitorMetrics(Counter, "test_counter", "test counter", LabelKey{"name"}),
			setMetricsValueFunc: func(metrics metrics) { metrics.(*PrometheusCounter).With(MetricsLabel{"name": "smiecj"}).Add(10) },
		},
	}
)

// 测试 prometheus monitor manager 基本功能
// 添加指标、查询指标
func TestPrometheusMonitorManager(t *testing.T) {
	manager := GetPrometheusMonitorManagerByConf(prometheusMonitorManagerConf{
		Port: 2112,
	})
	for _, currentTestCase := range testCaseArr {
		err := manager.AddMetrics(currentTestCase.desc)
		require.Empty(t, err)

		metrics, err := manager.GetMetrics(currentTestCase.metricsName)
		require.Empty(t, err)
		require.NotEmpty(t, metrics)

		currentTestCase.setMetricsValueFunc(metrics)
	}
}
