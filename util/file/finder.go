package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

// 返回指定文件的绝对路径，从调用的位置一直往上查找（一般用于配置文件查找）
func FindFilePath(fileName string) string {
	// 通过 caller 查找上一个调用方的路径，并往前遍历到项目根路径
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return fileName
	}
	filePathSplitArr := strings.Split(filePath, string(os.PathSeparator))
	filePathSplitArr = filePathSplitArr[:len(filePathSplitArr)-1]
	for len(filePathSplitArr) > 0 {
		folderPath := strings.Join(filePathSplitArr, string(os.PathSeparator))
		// find current level all file, and if this level match, return it
		fileArr, _ := ioutil.ReadDir(folderPath)
		for _, currentFile := range fileArr {
			if !currentFile.IsDir() && currentFile.Name() == fileName {
				return fmt.Sprintf("%s%s%s", folderPath, string(os.PathSeparator), fileName)
			}
		}
		filePathSplitArr = filePathSplitArr[:len(filePathSplitArr)-1]
	}
	return fileName
}
