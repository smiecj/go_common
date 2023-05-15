package alert

import (
	"testing"

	yamlconfig "github.com/smiecj/go_common/config/yaml"
	"github.com/smiecj/go_common/util/file"
	"github.com/smiecj/go_common/util/mail"
	"github.com/stretchr/testify/require"
)

const (
	testAlertTitle = "alert title"
	testAlertMsg   = "alert msg"

	localConfigFile = "conf_local.yaml"
)

func TestSendMail(t *testing.T) {
	configManager, err := yamlconfig.GetYamlConfigManager(file.FindFilePath(localConfigFile))
	require.Empty(t, err)

	sender, err := mail.NewSMTPMailSender(configManager)
	require.Empty(t, err)

	alerter := GetMailAlerter(sender)
	err = alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg))
	require.Empty(t, err)
}
