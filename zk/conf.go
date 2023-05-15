package zk

import (
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	zkConfigDefaultSpace = "zk"
	addressSplitor       = ","
)

// mysql 连接配置
type zkConnectOption struct {
	Address string `yaml:"address"`
	Home    string `yaml:"home"`
}

func (option zkConnectOption) getAddressArr() []string {
	if strings.Contains(option.Address, addressSplitor) {
		return strings.Split(option.Address, addressSplitor)
	}
	return []string{option.Address}
}

type conf struct {
	path       string
	data       string
	mode       int32
	timeout    time.Duration
	permission []zk.ACL
}

type confFunc func(*conf)

func SetPath(path string) confFunc {
	return func(conf *conf) {
		conf.path = path
	}
}

func SetData(data string) confFunc {
	return func(conf *conf) {
		conf.data = data
	}
}

func SetEphemeral() confFunc {
	return func(conf *conf) {
		conf.mode = conf.mode | zk.FlagEphemeral
	}
}

func SetSequence() confFunc {
	return func(conf *conf) {
		conf.mode = conf.mode | zk.FlagSequence
	}
}

func SetTTL(timeout time.Duration) confFunc {
	return func(conf *conf) {
		conf.timeout = timeout
	}
}

func defaultConf() *conf {
	conf := new(conf)
	conf.permission = zk.WorldACL(zk.PermAll)
	return conf
}

func getConf(funcArr ...confFunc) *conf {
	conf := defaultConf()
	for _, currentFunc := range funcArr {
		currentFunc(conf)
	}
	return conf
}
