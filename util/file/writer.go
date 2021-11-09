package file

import (
	"os"

	"github.com/smiecj/go_common/util/log"
)

const (
	ModeCreate = iota
	ModeAppend
)

const (
	defaultFileMode = 0644
)

// 覆盖/追加指定文件和内容
// 后续: 实现一个更通用的 fileappender ，可自动插入空行
func Write(fileAbsolutePath string, content []byte, fileMode int) int {
	var mode int
	switch fileMode {
	case ModeCreate:
		mode = os.O_CREATE | os.O_WRONLY
	case ModeAppend:
		mode = os.O_APPEND | os.O_WRONLY
	default:
		mode = os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(fileAbsolutePath, mode, defaultFileMode)
	if nil != err {
		log.Error("File init error: fileName: %s, err: %s", fileAbsolutePath, err.Error())
		return 0
	}
	defer f.Close()

	writeSize, err := f.Write(content)
	if nil != err {
		log.Error("File write error: fileName: %s, err: %s", fileAbsolutePath, err.Error())
	}
	return writeSize
}
