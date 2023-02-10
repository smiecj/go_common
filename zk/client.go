package zk

import (
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

const (
	root          = "/"
	pathSeparator = "/"
)

var (
	zkClientMap  map[zkConnectOption]Client
	zkClientLock sync.RWMutex
)

type Client interface {
	List(...confFunc) ([]string, error)
	Create(...confFunc) error
	Delete(...confFunc) error
	DeleteAll(...confFunc) error
	Status() string
}

type zkClient struct {
	option     zkConnectOption
	connection *zk.Conn
}

type emptyLogger struct{}

func (emptyLogger) Printf(format string, a ...interface{}) {}

func (client *zkClient) init() error {
	c, _, err := zk.Connect(client.option.getAddressArr(), time.Second)
	client.connection = c

	if nil != err {
		log.Error("[zkClient.init] zk connect %s failed: %s", client.option.Address, err.Error())
		return errorcode.BuildError(errorcode.ZKConnectFailed)
	}

	c.SetLogger(emptyLogger{})

	return nil
}

func (client *zkClient) List(funcArr ...confFunc) ([]string, error) {
	conf := client.getConf(funcArr...)
	listNodes, _, err := client.connection.Children(conf.path)
	if nil != err {
		log.Error("[zkClient.List] zk list failed: %s", err.Error())
		return nil, errorcode.BuildError(errorcode.ZKListFailed)
	}
	return listNodes, nil
}

func (client *zkClient) Create(funcArr ...confFunc) error {
	conf := client.getConf(funcArr...)
	// path string, data []byte, flags int32, acl []ACL
	_, err := client.connection.Create(conf.path, []byte(conf.data), conf.mode, conf.permission)
	if nil != err {
		log.Error("[zkClient.Create] zk create node %s failed: %s", conf.path, err.Error())
		return errorcode.BuildError(errorcode.ZKCreateFailed)
	}
	return nil
}

func (client *zkClient) Delete(funcArr ...confFunc) error {
	conf := client.getConf(funcArr...)
	// path string, data []byte, flags int32, acl []ACL
	err := client.connection.Delete(conf.path, -1)
	if nil != err {
		log.Error("[zkClient.Delete] zk delete node %s failed: %s", conf.path, err.Error())
		return errorcode.BuildError(errorcode.ZKCreateFailed)
	}
	return nil
}

// deleteall: refer: https://github.com/go-zookeeper/zk/issues/52
func (client *zkClient) DeleteAll(funcArr ...confFunc) error {
	conf := client.getConf(funcArr...)
	// to prevent delete all data, not allow delete root
	if conf.path == root {
		return errorcode.BuildError(errorcode.ZKDeleteRootFailed)
	}
	allChildArr := []string{}
	currentChildArr := []string{conf.path}
	for len(currentChildArr) != 0 {
		tempChildArr := []string{}
		for _, currentChild := range currentChildArr {
			childArr, _, err := client.connection.Children(currentChild)
			if nil != err {
				log.Error("[zkClient.DeleteAll] zk get child of %s failed: %s", currentChild, err.Error())
				return errorcode.BuildErrorWithMsg(errorcode.ZKDeleteFailed, err.Error())
			}
			for index := 0; index < len(childArr); index++ {
				childArr[index] = currentChild + pathSeparator + childArr[index]
			}
			allChildArr = append(allChildArr, currentChild)
			tempChildArr = append(tempChildArr, childArr...)
		}
		currentChildArr = tempChildArr
	}
	log.Info("[test] all child: %v", allChildArr)
	for index := len(allChildArr) - 1; index >= 0; index-- {
		currentChild := allChildArr[index]
		err := client.connection.Delete(currentChild, -1)
		if nil != err {
			log.Error("[zkClient.DeleteAll] zk delete node %s failed: %s", currentChild, err.Error())
			return errorcode.BuildErrorWithMsg(errorcode.ZKCreateFailed, err.Error())
		}
	}
	return nil
}

func (client *zkClient) Status() string {
	return "todo"
}

func (client *zkClient) getConf(funcArr ...confFunc) *conf {
	conf := defaultConf()
	for _, currentFunc := range funcArr {
		currentFunc(conf)
	}
	return conf
}

// 获取zk连接
func GetZKonnector(configManager config.Manager) (Client, error) {
	var client Client
	zkClientLock.RLock()

	option := zkConnectOption{}
	configManager.Unmarshal(zkConfigDefaultSpace, &option)

	if nil == zkClientMap {
		zkClientMap = make(map[zkConnectOption]Client)
	}

	client = zkClientMap[option]
	zkClientLock.RUnlock()

	if nil != client {
		return client, nil
	}

	zkClientLock.Lock()
	defer zkClientLock.Unlock()

	zkClient := new(zkClient)
	zkClient.option = option
	zkClient.init()
	zkClientMap[option] = zkClient
	return zkClient, nil
}
