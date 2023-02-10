package alert

import (
	"testing"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/util/mail"
	"github.com/stretchr/testify/require"
)

const (
	testAlertTitle = "alert title"
	testAlertMsg   = "alert msg"
)

func TestSendMail(t *testing.T) {
	configManager, err := config.GetYamlConfigManager("/tmp/mailconf.yml")
	require.Empty(t, err)

	sender, err := mail.NewSMTPMailSender(configManager)
	require.Empty(t, err)

	alerter := GetMailAlerter(sender)
	err = alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg))
	require.Empty(t, err)
}
