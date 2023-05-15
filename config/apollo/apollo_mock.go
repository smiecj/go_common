package config

import (
	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
)

// apollo config manager
type apolloConfigManagerMock struct{}

func (manager *apolloConfigManagerMock) Get(spaceName, key string) (ret interface{}, err error) {
	return "test_value", nil
}

func (manager *apolloConfigManagerMock) GetAllSpaceName() (retArr []string, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *apolloConfigManagerMock) GetSpace(spaceName string) (space config.Space, err error) {
	return nil, errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *apolloConfigManagerMock) Set(spaceName, key string, value interface{}) (err error) {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *apolloConfigManagerMock) Unmarshal(spaceName string, obj interface{}) error {
	return errorcode.BuildError(errorcode.ConfigMethodNotSupport)
}

func (manager *apolloConfigManagerMock) Update() error {
	return nil
}
