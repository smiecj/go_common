package zk

import (
	"testing"
	"time"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/util/file"
	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

const (
	localConfigFile = "conf_local.yaml"
	homePath        = "/"
	testPath        = "/smiecj"
)

var (
	testChildPathArr = []string{"/test", "/child", "/nephew"}
)

func TestZKConnect(t *testing.T) {
	configManager, err := config.GetYamlConfigManager(file.FindFilePath(localConfigFile))
	require.Nil(t, err)
	client, err := GetZKonnector(configManager)
	require.Nil(t, err)
	listNodes, _ := client.List(SetPath(homePath))
	require.NotEmpty(t, listNodes)
	log.Info("[test] zk nodes from root: %v", listNodes)

	// write new node
	err = client.Create(SetPath(testPath), SetEphemeral(), SetTTL(time.Minute))
	require.Nil(t, err)

	// check node
	listNodes, _ = client.List(SetPath(homePath))
	hasNode := false
	for _, currentNode := range listNodes {
		if "/"+currentNode == testPath {
			hasNode = true
			break
		}
	}
	require.True(t, hasNode)

	// delete node
	err = client.Delete(SetPath(testPath))
	require.Nil(t, err)

	// delete all
	currentPath := ""
	for _, currentChild := range testChildPathArr {
		currentPath = currentPath + currentChild
		client.Create(SetPath(currentPath))
	}
	err = client.DeleteAll(SetPath(testChildPathArr[0]))
	require.Nil(t, err)
}
