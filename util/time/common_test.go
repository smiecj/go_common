package time

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
)

func TestTimeCommand(t *testing.T) {
	log.Info(GetThisWeekLastDate())
	log.Info("%d", GetCurrentDateZeroTimestmapMill())
}
