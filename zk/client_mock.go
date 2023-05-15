package zk

import (
	"strings"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

type mockZKConnector struct {
	rootNode map[string]interface{}
}

func (c *mockZKConnector) List(funcArr ...confFunc) ([]string, error) {
	conf := getConf(funcArr...)
	pathMap, err := c.getPathMap(conf.path)
	if nil != err {
		return nil, errorcode.BuildError(errorcode.ZKListFailed)
	}

	retArr := []string{}
	for key := range pathMap {
		retArr = append(retArr, key)
	}
	return retArr, nil
}

func (c *mockZKConnector) Create(funcArr ...confFunc) error {
	conf := getConf(funcArr...)
	pathSplitArr := strings.Split(conf.path, "/")
	if len(pathSplitArr) == 0 {
		return nil
	}

	parentMap, err := c.getPathMapWithPrefix(conf.path, 1)
	if nil != err {
		log.Warn("[test] create: get root path failed")
		return errorcode.BuildError(errorcode.ZKCreateFailed)
	}

	// node has created
	lastNode := pathSplitArr[len(pathSplitArr)-1]
	if _, ok := parentMap[lastNode]; ok {
		return errorcode.BuildError(errorcode.ZKCreateFailed)
	}

	newNodeMap := make(map[string]interface{})
	parentMap[lastNode] = newNodeMap
	return nil
}

func (c *mockZKConnector) Delete(funcArr ...confFunc) error {
	conf := getConf(funcArr...)
	pathSplitArr := strings.Split(conf.path, "/")
	if len(pathSplitArr) == 0 {
		return nil
	}

	// get parent path map
	parentMap, err := c.getPathMapWithPrefix(conf.path, 1)
	if nil != err {
		return errorcode.BuildError(errorcode.ZKDeleteFailed)
	}

	// if node with child node, will delete failed
	lastNode := pathSplitArr[len(pathSplitArr)-1]
	if parentMap[lastNode] != nil && len(parentMap[lastNode].(map[string]interface{})) > 0 {
		return errorcode.BuildError(errorcode.ZKDeleteFailed)
	} else {
		delete(parentMap, lastNode)
	}
	return nil
}

func (c *mockZKConnector) DeleteAll(funcArr ...confFunc) error {
	conf := getConf(funcArr...)
	pathSplitArr := strings.Split(conf.path, "/")
	if len(pathSplitArr) == 0 {
		return nil
	}

	// get parent path map
	parentMap, err := c.getPathMapWithPrefix(conf.path, 1)
	if nil != err {
		return errorcode.BuildError(errorcode.ZKCreateFailed)
	}

	// delete directly
	lastNode := pathSplitArr[len(pathSplitArr)-1]
	delete(parentMap, lastNode)
	return nil
}

func (client *mockZKConnector) Status() string {
	return "fake"
}

func (c *mockZKConnector) getPathMapWithPrefix(path string, prefixLevel int) (map[string]interface{}, error) {
	// check path
	if len(path) == 0 || path[0] != '/' {
		return nil, errorcode.BuildError(errorcode.ServiceError)
	}

	pathSplitArr := strings.Split(path, "/")[1:]
	if prefixLevel < 0 || prefixLevel > len(pathSplitArr) {
		return nil, errorcode.BuildError(errorcode.ServiceError)
	}

	return c.getPathMap("/" + strings.Join(pathSplitArr[:len(pathSplitArr)-prefixLevel], "/"))
}

func (c *mockZKConnector) getPathMap(path string) (map[string]interface{}, error) {
	// check path
	if path == "/" {
		return c.rootNode, nil
	}
	if path == "" || path[0] != '/' {
		return nil, errorcode.BuildError(errorcode.ServiceError)
	}

	// split path
	pathSplitArr := strings.Split(path, "/")[1:]
	if len(pathSplitArr) == 0 {
		return nil, errorcode.BuildError(errorcode.ServiceError)
	}

	// find latest node
	currentRootMap := c.rootNode
	for _, currentPath := range pathSplitArr {
		if nextMap, ok := currentRootMap[currentPath]; ok {
			currentRootMap = nextMap.(map[string]interface{})
		} else {
			return nil, errorcode.BuildError(errorcode.ZKCreateFailed)
		}
	}
	return currentRootMap, nil
}

func getZKonnectorMock() (Client, error) {
	// default with zk nodes
	initNodeMap := make(map[string]interface{})
	initNodeMap["zookeeper"] = make(map[string]interface{})
	// initNodeMap["zookeeper"]["config"] = make(map[string]interface{})
	return &mockZKConnector{
		rootNode: initNodeMap,
	}, nil
}
