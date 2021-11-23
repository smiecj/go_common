# go common

提供 go 相关的公共库，供其他业务仓库使用，对一些常用的功能进行封装，方便使用

如：http 客户端、公共配置解析、公共数据库连接类等

# 已实现功能
## http 客户端

```
import "github.com/smiecj/go_common/http"
client := GetHTTPClient()
client.Do(Get(), Url("http://..."))
```

## file writer

```
import "github.com/smiecj/go_common/util/file"
file.Write("/tmp/test.log", "content", file.ModeCreate)
```

## logger

```
import "github.com/smiecj/go_common/util/log"
log.Info("msg: %s", msg)
```

## 错误码
```
import log "github.com/smiecj/go_common/errorcode"
return errorcode.ServiceError

// 返回自定义错误
return errorcode.BuildErrorWithMsg(errorcode.NetHandleFailed, "connect server failed")
// or: msg same as code
return errorcode.BuildError(errorcode.NetHandleFailed)
```

## RDB
### 本地内存
```
// 存入数据
localConnector := GetLocalMemoryConnector()
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// 查询数据
searchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
```

### 文件
```
// 存入数据
localConnector := GetLocalFileConnector(file_store_folder)

// insert field
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// insert object
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddObject(obj))


// 查询数据
// search field
searchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
// search object
searchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName), SearchSetObject(testStruct{}), SearchSetObjectArrType([]*testStruct{}))
// 注意 因为 通过 reflect 包 生成新对象 （调用 interface{} 方法）返回的是指针， 所以 SearchSetObjectArrType 一般需要设置指针数组，否则会转换失败
```

### mysql (gorm)
```
connector := GetMySQLConnector(MySQLConnectOption{Host: host, Port: port, User: user, Password: password})

// 存入数据
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"), 
    InsertAddObjectArr([]objectArr{obj1, obj2}), InsertSetObjectArrType([]object{}))
// 备注: 设置插入的数据数组格式的时候，直接插入一个大小为0 的数组即可，connector 内部逻辑会赋予 reflect.Type 格式
// 为什么需要 objectArrType: 和gorm的机制有关系，[]interface{} 类型无法正常判断数组内成员的 gorm tag

// 更新数据
updateRet, err := connector.Update(UpdateSetSpace("db_name", "table_name"),
		UpdateSetCondition("ID", "=", "1"),
		UpdateAddObject(object{Name: "ToUpdateName"}), UpdateAddKeyArr([]string{"name"}))
// 注意: gorm 默认会修改所有的字段，最好是通过 UpdateAddKeyArr 设置需要修改的字段列表

// 查询数据 - select
searchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetObjectArrType([]object{}), SearchSetPageCondition(0, 10))
objectArr := searchRet.ObjectArr.([]object{})
for _, currentObject := range objectArr {
	log.Info("current object: %v", currentObject)
}

// 查询数据 - count
searchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"), SearchSetCondition("ID", "=", "1"))
log.Info("count: %d", searchRet.Total)

// 查询数据 - distinct
searchRet, err := connector.Distinct(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetKeyArr([]string{"ID", "name"}))
for _, currentField := range searchRet.FieldArr {
	for columnName, value := range currentField.GetMap() {
		log.Info("distinct result: %s -> %s", columnName, value)
	}
}

// 删除数据
deleteRet, err := connector.Delete(DeleteSetSpace("db_name", "table_name"),
		DeleteSetCondition("ID", "=", "1"))
```

## 自定义配置 yaml 文件解析
### 获取 配置管理器
```
// 本地yaml配置文件
config := config.GetYamlConfig(config_file_name)
```

### 获取具体配置
```
// get config value by space name and key
value, err := config.get(space_name, key)

// or get by space
configSpace, err := config.getSpace(space_name)
value := configSpace.get(key)

// or transform to object
type dbConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
configSpace.Unmarshal(&dbConfig)
log.Info(dbConfig.Host)
```

## mail sender
```
sender := NewQQMailSender(MailSenderConf{
		Token:  "qq mail token",
		Sender: "sender qq mail account",
	})
err := sender.Send(AddReceiver("receiver mail account"), SetTitle("test_title"), SetContent("test_content"), SetNickName("nickname"))
```

# 待实现功能
## RPC 框架
