// package config config center
package config

// 配置中心定义
/*
 * 概念解释
 * manager: 配置中心最外层，负责管理配置、定期更新配置（后续实现）等
 * space: 配置空间概念，一般同一个配置中心下会包含多个子模块的配置，子模块通过“空间名”来区分
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
