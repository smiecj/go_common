package local

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	. "github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/json"
	"github.com/smiecj/go_common/util/log"
)

const (
	fileFormatObject   = "# object"
	fileFormatKeyValue = "# key-value"
	keyValueSplitor    = " --- "
	lineSeparator      = "\n"
)

var (
	fileConnectorMap  map[string]RDBConnector
	fileConnectorLock sync.RWMutex
)

// 本地文件存储
type localFileConnector struct {
	localFolderPath string
}

// 插入数据
// 文件名: db.table
// 文件格式:
/*
# object / key-value （数据格式）
object1 string format
object2 string format

or:
key1 --- key2 --- key3
value1 --- value2 --- value3 (object1)
key1 --- key2 --- key3
value1 --- value2 --- value3 (object2)
*/
func (connector *localFileConnector) Insert(funcArr ...RDBInsertConfigFunc) (ret UpdateRet, err error) {
	action := MakeRDBInsertAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	fileAbsolutePath := connector.getFileAbsolutePath(action.GetSpaceName())
	file, err := os.OpenFile(fileAbsolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		log.Error("[localFileConnector.Insert] write file failed, file name: %s, err: %s", fileAbsolutePath, err.Error())
		return
	}
	defer file.Close()

	fieldArr := action.GetFieldArr()
	objectArr := action.GetObjectArr()
	switch {
	case 0 != len(fieldArr):
		file.WriteString(fileFormatKeyValue + lineSeparator)
		for _, currentField := range fieldArr {
			keyStrAppender := new(bytes.Buffer)
			valueStrAppender := new(bytes.Buffer)
			for key, value := range currentField.GetMap() {
				if keyStrAppender.Len() != 0 {
					keyStrAppender.WriteString(keyValueSplitor)
					valueStrAppender.WriteString(keyValueSplitor)
				}
				// 注意: 内容中如果有换行符要替换掉
				keyStrAppender.WriteString(connector.cleanInvalidChar(key))
				valueStrAppender.WriteString(connector.cleanInvalidChar(value))
			}

			file.WriteString(keyStrAppender.String())
			file.WriteString(lineSeparator)
			file.WriteString(valueStrAppender.String())
			file.WriteString(lineSeparator)
		}

		ret.AffectedRows = len(fieldArr)
	case 0 != len(objectArr):
		file.WriteString(fileFormatObject + lineSeparator)
		for _, currentObject := range objectArr {
			objectBytes, _ := json.Marshal(currentObject)
			file.Write(objectBytes)
			file.WriteString(lineSeparator)
		}
		ret.AffectedRows = len(objectArr)
	}

	log.Info("[localFileConnector.Insert] Write file success, rows: %d", ret.AffectedRows)
	return
}

// 更新数据
// 文件不支持更新，只能覆盖
func (connector *localFileConnector) Update(funcArr ...RDBUpdateConfigFunc) (ret UpdateRet, err error) {
	err = fmt.Errorf("File storage not support update")
	return
}

// 删除数据
// 将会直接删除整个文件
func (connector *localFileConnector) Delete(funcArr ...RDBDeleteConfigFunc) (ret UpdateRet, err error) {
	action := MakeRDBDeleteAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	fileAbsolutePath := connector.getFileAbsolutePath(action.GetSpaceName())
	err = os.Remove(fileAbsolutePath)
	if nil != err {
		log.Error("[localFileConnector.Deletel] delete file failed, file name: %s, err: %s", fileAbsolutePath, err.Error())
	}
	// 后续: 获取文件内容中对应的数据条数
	return
}

// 备份数据
func (connector *localFileConnector) Backup(funcArr ...RDBBackupConfigFunc) (ret UpdateRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[localFileConnector.Backup] not implement")
}

// 查询数据
func (connector *localFileConnector) Search(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	fileAbsolutePath := connector.getFileAbsolutePath(action.GetSpaceName())
	// 文件不存在，返回失败结果
	if _, err := os.Stat(fileAbsolutePath); errors.Is(err, os.ErrNotExist) {
		return ret, fmt.Errorf("File is not exists")
	}

	file, _ := os.Open(fileAbsolutePath)
	defer file.Close()

	reader := bufio.NewReader(file)
	firstLine, _, _ := reader.ReadLine()
	if string(firstLine) == fileFormatKeyValue {
		for {
			keyBytes, _, readErr := reader.ReadLine()
			valueBytes, _, _ := reader.ReadLine()
			if nil != readErr {
				// read finish
				// 文件读取现在没做 limit，所以 total = len
				ret.Total = ret.Len
				return
			}

			keyArr := strings.Split(string(keyBytes), keyValueSplitor)
			valueArr := strings.Split(string(valueBytes), keyValueSplitor)

			currentField := BuildNewField()
			for index := 0; index < len(keyArr); index++ {
				currentField.AddKeyValue(keyArr[index], valueArr[index])
			}
			ret.AddField(currentField)
			ret.Len++
		}
	} else if string(firstLine) == fileFormatObject {
		// 查询条件中 对象为空 或者是 结果数组类型为空，则直接返回错误信息
		object, objectArrType := action.GetObject(), action.GetObjectArrType()
		if nil == object || nil == objectArrType {
			return ret, fmt.Errorf("Search base object is empty, please use 'SearchSetObject' to set object struct")
		}
		// reflect
		objValue := reflect.New(reflect.TypeOf(object))
		objectReflectArr := reflect.MakeSlice(objectArrType, 0, 0)
		ret.ObjectArr = make([]interface{}, 0)

		for {
			objectBytes, _, readErr := reader.ReadLine()
			if nil != readErr {
				// read finish
				ret.ObjectArr = objectReflectArr.Interface()
				ret.Total = ret.Len
				return
			}

			currentObj := objValue.Interface()
			err = json.Unmarshal(objectBytes, currentObj)
			if nil != err {
				log.Error("[localFileConnector.Search] object unmarshal failed, object: %s, err: %s", string(objectBytes), err.Error())
				return
			}

			objectReflectArr = reflect.Append(objectReflectArr, reflect.ValueOf(currentObj))
			ret.Len++
		}
	} else {
		return ret, fmt.Errorf("File format is not valid")
	}
}

// 统计数据量
func (connector *localFileConnector) Count(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	// 注意: 文件统计 暂时不支持按指定条件过滤，直接统计所有行数
	action := MakeRDBSearchAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	fileAbsolutePath := connector.getFileAbsolutePath(action.GetSpaceName())
	// 文件不存在，直接返回 （总数为0）
	if _, err = os.Stat(fileAbsolutePath); errors.Is(err, os.ErrNotExist) {
		return ret, nil
	}

	file, _ := os.Open(fileAbsolutePath)
	defer file.Close()

	reader := bufio.NewReader(file)
	for _, _, readErr := reader.ReadLine(); nil != readErr; {
		ret.Total++
	}
	// 去掉开头 表示文件格式的那一行
	ret.Total--
	return
}

// distinct
// file connector 暂不需要实现
func (connector *localFileConnector) Distinct(funcArr ...RDBSearchConfigFunc) (ret SearchRet, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.NotImplement, "[localFileConnector.Distinct] not implement")
}

// close
func (connector *localFileConnector) Close() error {
	return nil
}

// stat
func (connector *localFileConnector) Stat() (ret DBStat, err error) {
	return ret, errorcode.BuildErrorWithMsg(errorcode.DBStatFailed, err.Error())
}

// 公共方法: 获取需要操作的文件的绝对路径
func (connector *localFileConnector) getFileAbsolutePath(spaceName string) string {
	return fmt.Sprintf("%s%s%s", connector.localFolderPath, string(os.PathSeparator), spaceName)
}

// 公共方法: 清理不合法字符
func (connector *localFileConnector) cleanInvalidChar(toWriteStr string) (retStr string) {
	retStr = strings.ReplaceAll(toWriteStr, "\n", "")
	retStr = strings.ReplaceAll(retStr, keyValueSplitor, "")
	return
}

// 获取 根据目录路径匹配的单例
func GetLocalFileConnector(folderPath string) (RDBConnector, error) {
	var connector RDBConnector
	fileConnectorLock.RLock()
	if nil == fileConnectorMap {
		fileConnectorMap = make(map[string]RDBConnector)
	}
	connector = fileConnectorMap[folderPath]
	fileConnectorLock.RUnlock()

	if nil != connector {
		return connector, nil
	}

	fileConnectorLock.Lock()
	defer fileConnectorLock.Unlock()

	// 目录能成功创建，才能正常创建 connector
	err := os.MkdirAll(folderPath, os.ModeDir)
	if nil != err {
		log.Error("[GetLocalFileConnector] Get local connector failed, folder create failed: %s", folderPath)
		return nil, errorcode.BuildErrorWithMsg(errorcode.DBConnectFailed, err.Error())
	}
	fileConnector := new(localFileConnector)
	fileConnector.localFolderPath = folderPath
	fileConnectorMap[folderPath] = fileConnector
	return fileConnector, nil
}
