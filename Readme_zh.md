# go common

提供 go 相关的公共库，供其他业务仓库使用，对一些常用的功能进行封装，方便使用

如：http 客户端、公共配置解析、公共数据库连接类等

# 已实现功能
## http 客户端

```
import "github.com/smiecj/go_common/http"
client := http.DefaultHTTPClient()
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

// 默认日志格式: [2023-01-06 10:00:00] INFO 日志内容
log.Info("msg: %s", msg)

// 前缀日志: [prefix] [2023-01-06 10:00:00] INFO 日志内容
prefixLogger := log.PrefixLogger("prefix")
prefixLogger.Info("after prefix")
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
// 配置
mysql:
  host: localhost
  port: 3306
  user: root
  password: pwd
  max_life_time: 300
  max_idle_time: 300  

// 初始化
configManager, _ := config.GetYamlConfigManager(config_path)
connector, err := GetMySQLConnector(configManager)

// 存入数据
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"), 
    InsertAddObjectArr([]object{obj1, obj2}))
// 存入一条数据
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"),
  InsertSetObject(object))

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

// 查询数据 - join
// 注意: join 条件中的 表名 最好带上 `, 避免拼装成SQL 的时候被加上 单引号 导致SQL 执行错误
searchRet, err = connector.Search(SearchSetSpace("db_name", "table_name"),
  SearchSetObjectArrType([]object{}),
  SearchSetCondition("ID", "=", "1"),
  SearchAddJoin("db_name", "left_table", "left_field", "right_table", "right_field"),
  SearchSetKeyArr([]string{"ID", "name"})

// 备份数据
backupRet, err := connector.Backup(BackupSetSourceSpace("db_name", "source_table_name"),
    BackupSetTargetSpace("db_name", "target_table_name"))

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
  host: smtp 服务端域名
  port: smtp 服务端端口 (default 587)
  sender: 发送人邮箱
  token: smtp token
  receiver: 收件人列表(逗号分隔)
```

发送逻辑:
```
configManager, err := config.GetYamlConfigManager("mail_conf.yml")
sender, err := NewSMTPMailSender(configManager)
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
manager := GetPrometheusMonitorManagerByConf(
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

## zookeeper client

### 配置
zk:
  address: zk_node_1:2181,zk_node_2:2181
  home: /

### 使用
```
import "github.com/smiecj/go_common/zk"

// 获取客户端
client, err := zk.GetZKonnector(configManager)

// 创建节点
// 默认: 创建永久节点
err = client.Create(SetPath("/persistent"))
// 创建临时节点
err = client.Create(SetPath("/ephemeral"), SetEphemeral(), SetTTL(time.Minute))

// 列举子节点
listNodes, err = client.List(SetPath("/parent"))

// 删除节点
err = client.Delete(SetPath("/parent/child"))

// 删除包括所有子节点
err = client.DeleteAll(SetPath("/parent"))
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
import tickerutil "github.com/smiecj/go_common/util/time/ticker"

// 每天早上 8点 调度（会在 8:00~9:00 的某一时间调度）
ticker := tickerutil.NewFixHourTicker(8, tickerutil.SetFunc(func() error {dosomething...}))
ticker.Start()
// 处理错误流
// 如果设置了不需要处理错误 (ticker.SetIsIgnoreError), 则可以忽略下面的逻辑
go func() {
  for e := range ticker.Error() {
    log.Warn("[ticker] error: %s", e.Error())
  }
}()
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

## log
支持创建 子 log，自定义前缀 (mysql connector 用于打印连接地址)

## bug 修复
- ticker: 可重复 start 并且不会报错，不能重复启动

## RPC 框架

## 告警收敛

## 告警设置默认接收人

## db
- mysql 数据操作，超过指定数据量自动分批（防止一次性删除/插入过多数据）