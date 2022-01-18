package monitor

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	prometheusMonitorManagerInitOnce sync.Once
	prometheusMonitorManagerInstance Manager
)

// prometheus monitor manager 初始化配置
type PrometheusMonitorManagerConf struct {
	Port int
}

// prometheus monitor manager 实现
type prometheusMonitorManager struct {
	conf *monitorConf
}

func (manager *prometheusMonitorManager) init(conf PrometheusMonitorManagerConf) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	}()
}

func (manager *prometheusMonitorManager) GetMetrics(name string) (metrics, error) {
	return manager.conf.GetMetrics(name)
}

func (manager *prometheusMonitorManager) AddMetrics(desc monitorMetricsDesc) error {
	return manager.conf.AddMetrics(desc)
}

func (manager *prometheusMonitorManager) RemoveMetrics(name string) error {
	return manager.conf.RemoveMetricsByName(name)
}

// 获取 prometheus 监控管理器，单例模式
func GetPrometheusMonitorManager(conf PrometheusMonitorManagerConf) Manager {
	prometheusMonitorManagerInitOnce.Do(func() {
		manager := new(prometheusMonitorManager)
		manager.conf = newMonitorConf(monitorDescToPrometheusMetrics)
		manager.init(conf)
		prometheusMonitorManagerInstance = manager
	})
	return prometheusMonitorManagerInstance
}
