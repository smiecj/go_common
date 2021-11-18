// package config config center
package config

// 配置中心定义
type Manager interface {
	Get(string, string) (string, error)
	GetSpace(string) (space, error)
	Set(string, string, string) error
	Unmarshal(string, interface{}) error
}

// 配置属性
type space interface {
	Get(string) (string, error)
	Set(string, string) error
	Unmarshal(interface{}) error
}
