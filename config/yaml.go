package config

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"

	"gopkg.in/yaml.v3"
)

var (
	yamlConfigMapLock sync.RWMutex
	yamlConfigMap     map[string]Manager
)

// yaml config manager
type yamlManager struct {
	filePath string
	spaceMap map[string]space
}

// yaml config space
type yamlSpace struct {
	configMap map[string]interface{}
	mapLock   sync.RWMutex
}

// yaml space get config
func (space *yamlSpace) Get(key string) (interface{}, error) {
	space.mapLock.RLock()
	defer space.mapLock.RUnlock()

	ret, ok := space.configMap[key]
	if !ok {
		return "", errorcode.BuildError(errorcode.ConfigNotExist)
	}
	return ret, nil
}

func (space *yamlSpace) GetAllKey() (retArr []string, err error) {
	space.mapLock.RLock()
	defer space.mapLock.RUnlock()

	for key := range space.configMap {
		retArr = append(retArr, key)
	}
	return
}

// yaml space set config
func (space *yamlSpace) Set(key string, value interface{}) error {
	space.mapLock.Lock()
	defer space.mapLock.Unlock()

	space.configMap[key] = value
	return nil
}

// yaml space unmarshal
func (space *yamlSpace) Unmarshal(obj interface{}) error {
	space.mapLock.RLock()
	defer space.mapLock.RUnlock()

	configContent, _ := yaml.Marshal(space.configMap)
	return yaml.Unmarshal(configContent, obj)
}

// yaml config init
func (config *yamlManager) init() (err error) {
	file, err := os.Open(config.filePath)
	if nil != err {
		return
	}
	defer file.Close()

	fileContentBytes, err := ioutil.ReadAll(file)
	if nil != err {
		return
	}

	fullConfigMap := make(map[string]map[string]interface{})
	err = yaml.Unmarshal(fileContentBytes, &fullConfigMap)
	if nil != err {
		return
	}

	// 对每一层配置 都初始化一个 space
	config.spaceMap = make(map[string]space)
	for spaceName, configMap := range fullConfigMap {
		currentSpace := yamlSpace{configMap: configMap}
		config.spaceMap[spaceName] = &currentSpace
	}
	return
}

// yaml config manager get config
func (config *yamlManager) Get(spaceName, key string) (ret interface{}, err error) {
	space, err := config.getSpace(spaceName)
	if nil != err {
		return "", err
	}
	return space.Get(key)
}

// yaml get all space name
func (config *yamlManager) GetAllSpaceName() (retArr []string, err error) {
	// space map 目前只有读操作，不需要加锁
	for key := range config.spaceMap {
		retArr = append(retArr, key)
	}
	return
}

// yaml config get space
func (config *yamlManager) GetSpace(spaceName string) (space space, err error) {
	space, err = config.getSpace(spaceName)
	if nil != err {
		return
	}
	return space, nil
}

// yaml config set config
func (config *yamlManager) Set(spaceName, key string, value interface{}) (err error) {
	space, err := config.getSpace(spaceName)
	if nil != err {
		return err
	}
	return space.Set(key, value)
}

// yaml config unmarshal
func (config *yamlManager) Unmarshal(spaceName string, obj interface{}) error {
	space, err := config.getSpace(spaceName)
	if nil != err {
		return err
	}
	return space.Unmarshal(obj)
}

// common method: get space
func (config *yamlManager) getSpace(spaceName string) (space, error) {
	space := config.spaceMap[spaceName]
	if nil == space {
		return nil, errorcode.BuildError(errorcode.SpaceNotExist)
	}
	return space, nil
}

// update: yaml config not implement (current no need implement)
func (config *yamlManager) Update() error {
	return errorcode.BuildError(errorcode.NotImplement)
}

// 获取 yaml 配置中心单例
func GetYamlConfigManager(filePath string) (Manager, error) {
	var manager Manager
	yamlConfigMapLock.RLock()
	if nil == yamlConfigMap {
		yamlConfigMap = make(map[string]Manager)
	}
	manager = yamlConfigMap[filePath]
	yamlConfigMapLock.RUnlock()

	if nil != manager {
		return manager, nil
	}

	yamlConfigMapLock.Lock()
	defer yamlConfigMapLock.Unlock()

	// 必须要初始化成功
	yamlConfig := new(yamlManager)
	yamlConfig.filePath = filePath
	err := yamlConfig.init()
	if nil != err {
		log.Error("[config.GetYamlConfigManager] init yaml config failed: %s", err.Error())
		return nil, err
	}
	yamlConfigMap[filePath] = yamlConfig
	return yamlConfig, nil
}
