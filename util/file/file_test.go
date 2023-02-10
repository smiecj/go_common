package file

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

// find file path
func TestFindFilePath(t *testing.T) {
	filePath := FindFilePath("Makefile")
	isMatch, _ := regexp.MatchString("go_common", filePath)
	require.Equal(t, true, isMatch)
}
