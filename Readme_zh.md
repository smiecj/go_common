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

## DB 连接器
### 本地内存
```
// 存入数据
localConnector := GetLocalMemoryConnector()
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// 查询数据
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
```

### 文件
```
// 存入数据
localConnector, err := GetLocalFileConnector(file_store_folder)

// insert field
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// insert object
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddObject(obj))


// 查询数据
// search field
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
// search object
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName), SearchSetObject(testStruct{}), SearchSetObjectArrType([]*testStruct{}))
// 注意 因为 通过 reflect 包 生成新对象 （调用 interface{} 方法）返回的是指针， 所以 SearchSetObjectArrType 一般需要设置指针数组，否则会转换失败
```

### mysql (gorm)
```
connector, err := GetMySQLConnector(MySQLConnectOption{Host: host, Port: port, User: user, Password: password})

// 存入数据
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"), 
    InsertAddObjectArr([]objectArr{obj1, obj2}), InsertSetObjectArrType([]object{}))
// 备注: 设置插入的数据数组格式的时候，直接插入一个大小为0 的数组即可，connector 内部逻辑会赋予 reflect.Type 格式
// 为什么需要 objectArrType: 和gorm的机制有关系，[]interface{} 类型无法正常判断数组内成员的 gorm tag

// 更新数据
UpdateRet, err := connector.Update(UpdateSetSpace("db_name", "table_name"),
		UpdateSetCondition("ID", "=", "1"),
		UpdateAddObject(object{Name: "ToUpdateName"}), UpdateAddKeyArr([]string{"name"}))
// 注意: gorm 默认会修改所有的字段，最好是通过 UpdateAddKeyArr 设置需要修改的字段列表

// 查询数据 - select
SearchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetObjectArrType([]object{}), SearchSetPageCondition(0, 10))
objectArr := SearchRet.ObjectArr.([]object{})
for _, currentObject := range objectArr {
	log.Info("current object: %v", currentObject)
}

// 查询数据 - count
SearchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"), SearchSetCondition("ID", "=", "1"))
log.Info("count: %d", SearchRet.Total)

// 查询数据 - distinct
SearchRet, err := connector.Distinct(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetKeyArr([]string{"ID", "name"}))
for _, currentField := range SearchRet.FieldArr {
	for columnName, value := range currentField.GetMap() {
		log.Info("distinct result: %s -> %s", columnName, value)
	}
}

// 删除数据
deleteRet, err := connector.Delete(DeleteSetSpace("db_name", "table_name"),
		DeleteSetCondition("ID", "=", "1"))
```

### impala
引用: github.com/bippio/go-impala

```
connector, err := GetImpalaConnector(ImpalaConnectOption{Host: "impala_host", Port: 21050})

// count
ret, err := connector.Count(db.SearchSetSpace(db_name, table_name))
```

## 自定义配置 yaml 文件解析
### config file format
db: -- space
  mysql_host: localhost -- key: value
  mysql_port: 3306
  db_arr: 
    - school

### yaml 配置解析
```
// 从本地 yaml 配置文件 获取 配置管理
config, err := config.GetYamlConfigManager(config_file_name)

// get config value by space name and key
value, err := config.get(space_name, key)

// or get by space
configSpace, err := config.getSpace(space_name)
value := configSpace.get(key)

// or transform to object
type dbConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
configSpace.Unmarshal(&dbConfig)
log.Info(dbConfig.Host)
```

## mail sender 邮件发送器
mail_conf.yml 配置文件内容:
```
mail:
  token: smtp_token
  sender: sender_email
  receiver: receiver_email(split by comma)
```

发送逻辑:
```
configManager, err := config.GetYamlConfigManager("mail_conf.yml")
sender, err := NewQQMailSender(configManager)
// 可以不设置 AddReceiver，默认收件人 使用 配置中定义的收件人
err = sender.Send(AddReceiver("receiver mail account"), SetTitle("test_title"), SetContent("test_content"), SetNickName("nickname"))
```

## alerter 告警发送器
alerter 发送告警的功能依赖 sender，本身还会带简单的告警收敛的功能，比如10s 内一样的告警内容不会重复发送

```
// 获取邮件告警发送器，需要先设置 sender，获取 sender 的方式参考 mail sender
alerter := GetMailAlerter(sender, testDefaultReceiver)

// 发送告警
alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg))
```

## monitor
监控功能封装，提供基本的监控指标 配置方法，上层实现包括 prometheus，使用者不需要关注监控的具体实现
当前支持指标: gauge, counter

```
// 获取 Prometheus 监控，单例模式
manager := GetPrometheusMonitorManager(
	PrometheusMonitorManagerConf{
		Port: 开放端口,
})

// 添加一个指标
metricsDesc := NewMonitorMetrics(Gauge, "test_gauge", "test gauge", LabelKey{"name"})
err := manager.AddMetrics(metricsDesc)

// 获取指标
metrics, err := manager.GetMetrics(currentTestCase.metricsName)

// 设置指标值（需要用户侧在外面自行强转成具体的类型，比如 PrometheusGauge ）
metrics.(*PrometheusGauge).With(MetricsLabel{"name": "smiecj"}).Set(10)
```

## 其他公共工具方法

### time - 格式化相关工具方法
```
import "github.com/smiecj/go_common/util/time"

// 获取当前时间戳 (格式: 2006-01-02 15:04:05)
time.GetCurrentTimestamp()
```

### time - 指定小时调度
```
import "github.com/smiecj/go_common/util/time"

// 每天早上8点调度
ticker := time.NewFixHourTicker(8, time.SetFunc(func() error {dosomething...}))
ticker.Start()
// ticker.Stop()
```

### net
```
import "github.com/smiecj/go_common/util/net"

// 获取本机 ipv4 ip
ip, err := net.GetLocalIPV4()

// 获取本机是否占用了指定端口
isUsed := net.CheckLocalPortIsUsed(22)
```

# 待实现功能
## RPC 框架

## 告警收敛

## 告警设置默认接收人