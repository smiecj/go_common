package monitor

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/util/log"
	"github.com/smiecj/go_common/util/net"
)

const (
	configSpaceMonitor = "monitor"
)

var (
	prometheusMonitorManagerInitOnce sync.Once
	prometheusMonitorManagerInstance Manager
)

// prometheus monitor manager 初始化配置
type prometheusMonitorManagerConf struct {
	Port int `yaml:"port"`
}

func (conf prometheusMonitorManagerConf) portUsed() bool {
	// 检查端口是否被占用，如果被占用直接退出
	if net.CheckLocalPortIsUsed(conf.Port) {
		log.Error("[prometheusMonitorManagerConf.checkPortNotUse] monitor port: %d is used", conf.Port)
		return true
	}
	return false
}

// prometheus monitor manager 实现
type prometheusMonitorManager struct {
	conf *monitorConf
}

func (manager *prometheusMonitorManager) init(conf prometheusMonitorManagerConf) {
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

func GetPrometheusMonitorManager(configManager config.Manager) Manager {
	conf := prometheusMonitorManagerConf{}
	configManager.Unmarshal(configSpaceMonitor, &conf)
	return getPrometheusMonitorManager(conf)
}

func GetPrometheusMonitorManagerByConf(conf prometheusMonitorManagerConf) Manager {
	return getPrometheusMonitorManager(conf)
}

// 获取 prometheus 监控管理器，单例模式
func getPrometheusMonitorManager(conf prometheusMonitorManagerConf) Manager {
	prometheusMonitorManagerInitOnce.Do(func() {
		if conf.portUsed() {
			return
		}

		manager := new(prometheusMonitorManager)
		manager.conf = newMonitorConf(monitorDescToPrometheusMetrics)
		manager.init(conf)
		prometheusMonitorManagerInstance = manager
	})
	return prometheusMonitorManagerInstance
}
