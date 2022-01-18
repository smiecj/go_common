package monitor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 测试 prometheus guage
func TestPrometheusGauge(t *testing.T) {
	gaugeMetrics := NewMonitorMetrics(Gauge, "test_gauge", "guage desc", LabelKey{"name"})
	gauge, err := monitorDescToPrometheusMetrics(gaugeMetrics)
	require.Empty(t, err)
	gauge.(*PrometheusGauge).With(MetricsLabel{"name": "smiecj"}).Set(10)
}

// 测试 prometheus counter
func TestPrometheusCounter(t *testing.T) {
	counterMetrics := NewMonitorMetrics(Counter, "test_counter", "guage desc", LabelKey{"name"})
	counter, err := monitorDescToPrometheusMetrics(counterMetrics)
	require.Empty(t, err)
	counter.(*PrometheusCounter).With(MetricsLabel{"name": "smiecj"}).Inc()
}
