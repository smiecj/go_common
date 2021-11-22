// package config config center
package config

// 配置中心定义
type Manager interface {
	Get(string, string) (interface{}, error)
	GetSpace(string) (space, error)
	Set(string, string, interface{}) error
	Unmarshal(string, interface{}) error
}

// 配置属性
type space interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
	Unmarshal(interface{}) error
}
