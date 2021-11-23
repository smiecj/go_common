package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

// 寻找指定文件的路径，不能超过当前项目路径（一般是配置文件）
func FindFilePath(fileName string) (string, bool) {
	// 通过 caller 查找上一个调用方的路径，并往前遍历到项目根路径
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return "", false
	}
	filePathSplitArr := strings.Split(filePath, string(os.PathSeparator))
	filePathSplitArr = filePathSplitArr[:len(filePathSplitArr)-1]
	for len(filePathSplitArr) > 0 {
		folderPath := strings.Join(filePathSplitArr[:len(filePathSplitArr)-1], string(os.PathSeparator))

		// find current level all file, and if this level match, return it
		fileArr, _ := ioutil.ReadDir(folderPath)
		for _, currentFile := range fileArr {
			if !currentFile.IsDir() && currentFile.Name() == fileName {
				return fmt.Sprintf("%s%s%s", folderPath, string(os.PathSeparator), fileName), true
			}
		}
		filePathSplitArr = filePathSplitArr[:len(filePathSplitArr)-1]
	}
	return "", false
}
