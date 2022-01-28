// package config config center
package config

// 配置中心定义
/*
 * 概念解释
 * manager: 配置中心最外层，负责管理配置、定期更新配置（后续实现）等
 * space: 配置空间概念，一般同一个配置中心下会包含多个子模块的配置，子模块通过这个“空间名”来区分
 * 示例: 一个网关服务名称为，读取的配置文件名为，该服务下有几个子模块:
 * register(负责注册网关信息), balancer(负责将请求分配给具体的下游服务节点), interface(请求接入层), requester(请求发送、结果处理)
 * 它们对应不同的space
 */
type Manager interface {
	Get(string, string) (interface{}, error)
	GetSpace(string) (space, error)
	GetAllSpaceName() ([]string, error)
	Set(string, string, interface{}) error
	Unmarshal(string, interface{}) error
	Update() error
}

// 配置空间
type space interface {
	Get(string) (interface{}, error)
	GetAllKey() ([]string, error)
	Set(string, interface{}) error
	Unmarshal(interface{}) error
}
