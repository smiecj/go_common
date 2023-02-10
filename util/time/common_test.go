package time

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

func TestTimeCommand(t *testing.T) {
	log.Info(GetThisWeekLastDate())
	log.Info("%d", GetCurrentDateZeroTimestmapMill())

	timeStr, err := GetTimestampByUnixtimeStr("1438398103")
	require.Nil(t, err)
	log.Info("%s", timeStr)
	log.Info("%s", GetTimestampByUnixMill(1438398103463))
}
