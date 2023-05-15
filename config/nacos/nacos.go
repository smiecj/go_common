package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	config "github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

const (
	nacosConfigSpace = "nacos"
)

type emptyLogger struct {
}

func (logger *emptyLogger) Info(args ...interface{})               {}
func (logger *emptyLogger) Warn(args ...interface{})               {}
func (logger *emptyLogger) Error(args ...interface{})              {}
func (logger *emptyLogger) Debug(args ...interface{})              {}
func (logger *emptyLogger) Infof(fmt string, args ...interface{})  {}
func (logger *emptyLogger) Warnf(fmt string, args ...interface{})  {}
func (logger *emptyLogger) Errorf(fmt string, args ...interface{}) {}
func (logger *emptyLogger) Debugf(fmt string, args ...interface{}) {}

type nacosConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	NamespaceId string `yaml:"namespace_id"`
}

// nacos config manager
type nacosConfigManager struct {
	nacosConfig nacosConfig
	client      config_client.IConfigClient
}

// nacos config init
func (manager *nacosConfigManager) init(conf nacosConfig) (err error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(manager.nacosConfig.Host, uint64(manager.nacosConfig.Port), constant.WithContextPath("/nacos")),
	}
	cc := *constant.NewClientConfig(
		constant.WithTimeoutMs(10*1000),
		constant.WithBeatInterval(2*1000),
		constant.WithCustomLogger(&emptyLogger{}),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUpdateCacheWhenEmpty(true),
		constant.WithUsername(manager.nacosConfig.User),
		constant.WithPassword(manager.nacosConfig.Password),
		constant.WithNamespaceId(manager.nacosConfig.NamespaceId),
	)

	client, getClientErr := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if nil != getClientErr {
		log.Warn("[NacosConfigManager] get client err: %s", getClientErr.Error())
		return getClientErr
	}

	manager.client = client
	return nil
}

// nacos config manager get config
func (manager *nacosConfigManager) Get(spaceName, key string) (ret interface{}, err error) {
	val, err := manager.client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  spaceName,
	})
	ret = val
	return
}

// not supported
func (manager *nacosConfigManager) GetAllSpaceName() (retArr []string, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *nacosConfigManager) GetSpace(spaceName string) (space config.Space, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// yaml config set config
func (manager *nacosConfigManager) Set(spaceName, key string, value interface{}) (err error) {
	valStr, ok := value.(string)
	if !ok {
		return errorcode.BuildErrorWithMsg(errorcode.ServiceError, "config set value not string")
	}
	_, publishConfErr := manager.client.PublishConfig(vo.ConfigParam{
		DataId:  key,
		Group:   spaceName,
		Content: valStr,
	})
	return publishConfErr
}

// nacos 接口不支持获取所有 group 、所有 key，因此无法支持 Unmarshal 方法
// https://nacos.io/en-us/docs/v2/guide/user/open-api.html
func (manager *nacosConfigManager) Unmarshal(spaceName string, obj interface{}) error {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *nacosConfigManager) Update() error {
	// manager.client.
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// get nacos config manager (server address config rely on local file)
func GetNacosConfigManager(yamlConfigManager config.Manager) (config.Manager, error) {
	conf := nacosConfig{}
	getConfigErr := yamlConfigManager.Unmarshal(nacosConfigSpace, &conf)
	if nil != getConfigErr {
		log.Warn("[GetNacosConfigManager] parse config err: %s", getConfigErr.Error())
		return nil, getConfigErr
	}
	if conf.Host == "mock" {
		return &nacosConfigManagerMock{}, nil
	}

	nacosConfigManager := new(nacosConfigManager)
	err := nacosConfigManager.init(conf)
	if nil != err {
		log.Error("[config.GetNacosConfigManager] init nacos config failed: %s", err.Error())
		return nil, err
	}
	return nacosConfigManager, nil
}
