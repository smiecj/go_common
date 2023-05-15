package config

import (
	config "github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
)

type nacosConfigManagerMock struct{}

// nacos config manager get config
func (manager *nacosConfigManagerMock) Get(spaceName, key string) (ret interface{}, err error) {
	return "test_value", nil
}

// not supported
func (manager *nacosConfigManagerMock) GetAllSpaceName() (retArr []string, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// not supported
func (manager *nacosConfigManagerMock) GetSpace(spaceName string) (space config.Space, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

// yaml config set config
func (manager *nacosConfigManagerMock) Set(spaceName, key string, value interface{}) (err error) {
	return nil
}

func (manager *nacosConfigManagerMock) Unmarshal(spaceName string, obj interface{}) error {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *nacosConfigManagerMock) Update() error {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}
