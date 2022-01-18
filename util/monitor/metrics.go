package monitor

import (
	"sync"

	"github.com/smiecj/go_common/errorcode"
)

const (
	Gauge MetricsType = iota
	Counter
	Summary
	Histogram
)

// 指标类型
// openmetrics: gauge, counter, summary, histogram
type MetricsType int

// 指标标签
type MetricsLabel map[string]string

// 标签列表
type LabelKey []string

// 监控指标描述
type monitorMetricsDesc struct {
	// _type: 指标类型，参考 MetricsType
	_type MetricsType

	// name: 指标名称，一般不能重复
	name string

	// description: 额外描述信息
	description string

	// labelKeyArr: 维度列表
	labelKeyArr LabelKey
}

// 创建一个监控指标基本描述/配置
func NewMonitorMetrics(t MetricsType, name string, desc string, labelKeyArr LabelKey) monitorMetricsDesc {
	return monitorMetricsDesc{
		_type:       t,
		name:        name,
		description: desc,
		labelKeyArr: LabelKey(append([]string{}, labelKeyArr...)),
	}
}

type confToMetricsFunc func(monitorMetricsDesc) (metrics, error)

// 监控基本配置
// 主要功能: 维护底层需要监控的指标
type monitorConf struct {
	// 监控指标配置
	monitorMetricsMap map[string]monitorMetricsDesc

	// 具体指标
	metricsMap map[string]metrics

	// metrics 配置和实际指标之间的 转换方法，需要在 monitor manager 的具体实现中提供
	confToMetricsFunc confToMetricsFunc

	// 读写锁
	lock sync.RWMutex
}

// 新建 & 初始化监控基本配置
func newMonitorConf(confToMetricsFunc confToMetricsFunc) *monitorConf {
	return &monitorConf{
		monitorMetricsMap: make(map[string]monitorMetricsDesc),
		metricsMap:        make(map[string]metrics),
		confToMetricsFunc: confToMetricsFunc,
		lock:              sync.RWMutex{},
	}
}

// 监控配置: 添加一个指标 （未做去重）
func (conf *monitorConf) AddMetrics(metricsDesc monitorMetricsDesc) error {
	// 考虑到调用 AddMetrics 方法的次数不会太多，不再进行读写锁的区分
	conf.lock.Lock()
	defer conf.lock.Unlock()

	// 判断是否已经有对应指标，有的话不再重复添加
	if nil != conf.metricsMap[metricsDesc.name] {
		return errorcode.BuildError(errorcode.MonitorMetricsExists)
	}

	metrics, err := conf.confToMetricsFunc(metricsDesc)
	if nil != err {
		return err
	}

	conf.monitorMetricsMap[metricsDesc.name] = metricsDesc
	conf.metricsMap[metricsDesc.name] = metrics
	return nil
}

// 监控配置: 删除一个指标
func (conf *monitorConf) RemoveMetricsByName(metricsName string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	if nil == conf.metricsMap[metricsName] {
		return errorcode.BuildError(errorcode.MonitorMetricsNotExists)
	}

	delete(conf.monitorMetricsMap, metricsName)
	delete(conf.metricsMap, metricsName)
	return nil
}

// 监控配置: 获取一个指标
func (conf *monitorConf) GetMetrics(metricsName string) (metrics, error) {
	conf.lock.RLock()
	defer conf.lock.RUnlock()

	if nil == conf.metricsMap[metricsName] {
		return nil, errorcode.BuildError(errorcode.MonitorMetricsNotExists)
	}

	return conf.metricsMap[metricsName], nil
}

// metrics 类的统一抽象，方便 manager 返回该对象
// 具体监控指标
type metrics interface {
	MetricsName() string
}

// openmetrics 监控对象抽象 - Gauge
type gauge interface {
	Set(float64)
}

// openmetrics 监控对象抽象 - Counter
type counter interface {
	Add(float64)
	Inc()
}

// 后续: 实现 summary 和 histogram
