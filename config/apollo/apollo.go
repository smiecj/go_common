package config

import (
	"fmt"

	"github.com/smiecj/agollo/v4"
	apolloconfig "github.com/smiecj/agollo/v4/env/config"
	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

const (
	apolloConfigSpace = "apollo"
)

type apolloConfig struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	AppId            string `yaml:"app_id"`
	Cluster          string `yaml:"cluster"`
	Secret           string `yaml:"secret"`
	DefaultNamespace string `yaml:"default_namespace"`
}

// apollo config manager
type apolloConfigManager struct {
	apolloConfig apolloConfig
	client       agollo.Client
}

// apollo config init
func (manager *apolloConfigManager) init(conf apolloConfig) (err error) {
	c := &apolloconfig.AppConfig{
		AppID:         manager.apolloConfig.AppId,
		Cluster:       manager.apolloConfig.Cluster,
		IP:            fmt.Sprintf("http://%s:%d", manager.apolloConfig.Host, manager.apolloConfig.Port),
		Secret:        manager.apolloConfig.Secret,
		NamespaceName: manager.apolloConfig.DefaultNamespace,
	}

	client, getClientErr := agollo.StartWithConfig(func() (*apolloconfig.AppConfig, error) {
		return c, nil
	})

	if nil != getClientErr {
		log.Warn("[ApolloConfigManager] get client err: %s", getClientErr.Error())
		return getClientErr
	}

	manager.client = client
	return nil
}

// apollo config manager get config
func (manager *apolloConfigManager) Get(spaceName, key string) (ret interface{}, err error) {
	// conf := manager.client.GetConfigCache(spaceName)
	conf := manager.client.GetConfigAndInit(spaceName)
	// conf := manager.client.GetConfig(spaceName)
	if nil == conf {
		return nil, errorcode.BuildError(errorcode.ConfigNotExist)
	}
	return conf.GetValue(key), nil
}

// not supported
func (manager *apolloConfigManager) GetAllSpaceName() (retArr []string, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *apolloConfigManager) GetSpace(spaceName string) (space config.Space, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *apolloConfigManager) Set(spaceName, key string, value interface{}) (err error) {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *apolloConfigManager) Unmarshal(spaceName string, obj interface{}) error {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *apolloConfigManager) Update() error {
	manager.client.GetApolloConfigCache().Clear()
	return nil
}

// get apollo config manager
func GetApolloConfigManager(yamlConfigManager config.Manager) (config.Manager, error) {
	conf := apolloConfig{}
	getConfigErr := yamlConfigManager.Unmarshal(apolloConfigSpace, &conf)
	if nil != getConfigErr {
		log.Warn("[GetApolloConfigManager] parse config err: %s", getConfigErr.Error())
		return nil, getConfigErr
	}
	if conf.Host == "mock" {
		return &apolloConfigManagerMock{}, nil
	}

	apolloConfigManager := new(apolloConfigManager)
	err := apolloConfigManager.init(conf)
	if nil != err {
		log.Error("[config.GetApolloConfigManager] init apollo config failed: %s", err.Error())
		return nil, err
	}
	return apolloConfigManager, nil
}
