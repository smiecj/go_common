package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smiecj/go_common/errorcode"
)

// prometheus gauge 定义
type PrometheusGauge struct {
	// 指标名称，仅用于适配 metrics 的 MetricsName 接口
	name string
	// 当前具体指标的标签（维度），仅用于With 方法调用时作为中间参数透传
	label MetricsLabel
	// 具体 prometheus 监控接口定义
	gaugeVec *prometheus.GaugeVec
}

func (gauge *PrometheusGauge) MetricsName() string {
	return gauge.name
}

// 设置标签
func (gauge *PrometheusGauge) With(label MetricsLabel) *PrometheusGauge {
	// 需要创建一个新的对象，保证并发调用不冲突
	newGauge := PrometheusGauge{
		label:    label,
		gaugeVec: gauge.gaugeVec,
	}
	return &newGauge
}

// 设置指标值
func (gauge *PrometheusGauge) Set(val float64) {
	gauge.gaugeVec.With(prometheus.Labels(gauge.label)).Set(val)
}

// 新建一个 prometheus gauge
func newPrometheusGauge(metrics monitorMetricsDesc) *PrometheusGauge {
	return &PrometheusGauge{
		name: metrics.name,
		gaugeVec: promauto.NewGaugeVec(
			prometheus.GaugeOpts{Name: metrics.name, Help: metrics.description},
			metrics.labelKeyArr),
	}
}

// prometheus counter 定义
type PrometheusCounter struct {
	name       string
	label      MetricsLabel
	counterVec *prometheus.CounterVec
}

func (counter *PrometheusCounter) MetricsName() string {
	return counter.name
}

// counter: 设置要记录的 label 列表
func (counter *PrometheusCounter) With(label MetricsLabel) *PrometheusCounter {
	// 需要创建一个新的对象，保证并发调用不冲突
	newCounter := PrometheusCounter{
		label:      label,
		counterVec: counter.counterVec,
	}
	return &newCounter
}

// counter: 递增指定的值（可以是负数）
func (counter *PrometheusCounter) Add(val float64) {
	counter.counterVec.With(prometheus.Labels(counter.label)).Add(val)
}

// counter: 加一
func (counter *PrometheusCounter) Inc() {
	counter.counterVec.With(prometheus.Labels(counter.label)).Inc()
}

// 新建一个 prometheus counter
func newPrometheusCounter(metrics monitorMetricsDesc) *PrometheusCounter {
	return &PrometheusCounter{
		name: metrics.name,
		counterVec: promauto.NewCounterVec(
			prometheus.CounterOpts{Name: metrics.name, Help: metrics.description},
			metrics.labelKeyArr),
	}
}

// prometheus 指标配置和具体指标之间的转换方法，由 monitor conf 初始化的时候使用
func monitorDescToPrometheusMetrics(desc monitorMetricsDesc) (metrics, error) {
	// 配置基本参数判断
	if desc.name == "" {
		return nil, errorcode.BuildErrorWithMsg(
			errorcode.MonitorMetricsTransformFailed,
			"metrics name empty",
		)
	}

	switch desc._type {
	case Gauge:
		return newPrometheusGauge(desc), nil
	case Counter:
		return newPrometheusCounter(desc), nil
	default:
		return nil, errorcode.BuildErrorWithMsg(
			errorcode.MonitorMetricsTransformFailed,
			"metrics type not supported",
		)
	}
}
