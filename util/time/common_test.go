package time

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

func TestTimeCommand(t *testing.T) {
	log.Info(ThisWeekLastDate())
	log.Info("%d", CurrentDateZeroTimestmapMill())

	timeStr, err := TimestampByUnixtimeStr("1438398103")
	require.Nil(t, err)
	log.Info("%s", timeStr)
	log.Info("%s", TimestampByUnixMill(1438398103463))

	timeStr, _ = TimestampByGoFormat("29 July 2014")
	log.Info("%s", timeStr)
	timeStr, _ = TimestampByGoFormat("3 August 2017")
	log.Info("%s", timeStr)
}
