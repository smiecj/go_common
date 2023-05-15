package zk

import (
	"testing"
	"time"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/util/file"
	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

const (
	localConfigFile = "conf_local.yaml"
	homePath        = "/"
	rootPath        = "/root"
)

var (
	testChildPathArr = []string{"/parent", "/child", "/nephew"}
)

func TestZKConnect(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(localConfigFile))
	require.Nil(t, err)

	client, err := GetZKonnector(configManager)
	require.Nil(t, err)
	listNodes, _ := client.List(SetPath(homePath))
	require.NotEmpty(t, listNodes)
	log.Info("[test] zk nodes from root: %v", listNodes)

	// write new node
	err = client.Create(SetPath(rootPath), SetEphemeral(), SetTTL(time.Minute))
	require.Nil(t, err)

	// check node
	listNodes, _ = client.List(SetPath(homePath))
	hasNode := false
	for _, currentNode := range listNodes {
		if "/"+currentNode == rootPath {
			hasNode = true
			break
		}
	}
	require.True(t, hasNode)

	// delete node
	err = client.Delete(SetPath(rootPath))
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
