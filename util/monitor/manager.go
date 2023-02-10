// package monitor 监控管理器封装
package monitor

// 监控管理器定义，提供基本定义，指标增删查改接口，本身不关心监控指标如何上报
type Manager interface {
	GetMetrics(string) (metrics, error)
	AddMetrics(desc monitorMetricsDesc) error
	RemoveMetrics(string) error
}
