[中文](https://github.com/smiecj/go_common/blob/master/Readme_zh.md)

# go common

go common library, supply some common library for other repo to use

e.g. http client, config manager and db connector etc.

# implemented features
## http client

```
import "github.com/smiecj/go_common/http"
client := http.DefaultHTTPClient()
client.Do(Get(), Url("http://..."))
```

## logger

```
import "github.com/smiecj/go_common/util/log"

// default log format: [2023-01-06 10:00:00] INFO log content
log.Info("msg: %s", msg)

// prefix logger: [prefix] [2023-01-06 10:00:00] INFO log content
prefixLogger := log.PrefixLogger("prefix")
prefixLogger.Info("after prefix")
```

## file writer

```
import "github.com/smiecj/go_common/util/file"
file.Write("/tmp/test.log", "content", file.ModeCreate)
```

## errorcode
```
import log "github.com/smiecj/go_common/errorcode"
return errorcode.ServiceError

// return self define error code
return errorcode.BuildErrorWithMsg(errorcode.NetHandleFailed, "connect server failed")
// or: msg same as code
return errorcode.BuildError(errorcode.NetHandleFailed)
```

## DB connector
### local memory connector
```
// store
localConnector := GetLocalMemoryConnector()
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// query
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
```

### local file connector
```
// get connector
localConnector, err := GetLocalFileConnector(file_store_folder)

// insert field
field := db.BuildNewField()
field.AddKeyValue("key", "value")
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddField(field))

// insert object
insertRet, err := localConnector.Insert(InsertSetSpace(dbName, tableName), InsertAddObject(obj))


// query
// search field
searchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
// search object
searchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName), SearchSetObject(testStruct{}), SearchSetObjectArrType([]*testStruct{}))
// notice: local file connector use reflect lib to unmarshal object everyline, so please use pointer array when set object array type
actualObjectArr := searchRet.ObjectArr.([]*testStruct{})
```

### mysql (gorm)
```
// config
mysql:
  host: localhost
  port: 3306
  user: root
  password: pwd
  max_life_time: 300
  max_idle_time: 300  

// init
configManager, _ := config.GetYamlConfigManager(config_path)
connector, err := GetMySQLConnector(configManager)

// store
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"), 
    InsertAddObjectArr([]object{obj1, obj2}))
// store one record
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"),
  InsertSetObject(object))

// update
UpdateRet, err := connector.Update(UpdateSetSpace("db_name", "table_name"),
		UpdateSetCondition("ID", "=", "1"),
		UpdateAddObject(object{Name: "ToUpdateName"}), UpdateAddKeyArr([]string{"name"}))
// notice: gorm will update all fields by default value, so it's better to set keyArr when update

// query - select
SearchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetObjectArrType([]object{}), SearchSetPageCondition(0, 10))
objectArr := SearchRet.ObjectArr.([]object{})
for _, currentObject := range objectArr {
	log.Info("current object: %v", currentObject)
}

// query - count
SearchRet, err := connector.Search(SearchSetSpace("db_name", "table_name"), SearchSetCondition("ID", "=", "1"))
log.Info("count: %d", SearchRet.Total)

// query - distinct
SearchRet, err := connector.Distinct(SearchSetSpace("db_name", "table_name"),
    SearchSetCondition("ID", "=", "1"), SearchSetKeyArr([]string{"ID", "name"}))
for _, currentField := range SearchRet.FieldArr {
	for columnName, value := range currentField.GetMap() {
		log.Info("distinct result: %s -> %s", columnName, value)
	}
}

// query - join
// notice: table name in join condition should contains '`' to avoid append ' character
searchRet, err = connector.Search(SearchSetSpace("db_name", "table_name"),
  SearchSetObjectArrType([]object{}),
  SearchSetCondition("ID", "=", "1"),
  SearchAddJoin("db_name", "left_table", "left_field", "right_table", "right_field"),
  SearchSetKeyArr([]string{"ID", "name"})

// backup
backupRet, err := connector.Backup(BackupSetSourceSpace("db_name", "source_table_name"),
    BackupSetTargetSpace("db_name", "target_table_name"))

// delete
deleteRet, err := connector.Delete(DeleteSetSpace("db_name", "table_name"),
		DeleteSetCondition("ID", "=", "1"))
```

### impala
refer: github.com/bippio/go-impala

```
connector, err := GetImpalaConnector(ImpalaConnectOption{Host: "impala_host", Port: 21050})

// count
ret, err := connector.Count(db.SearchSetSpace(db_name, table_name))
```

## config manager
### config file format
db: -- space
  mysql_host: localhost -- key: value
  mysql_port: 3306
  db_arr: 
    - school

### yaml config manager
```
// get config manager
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

## mail sender
mail_conf.yml:
```
mail:
  host: smtp_server_host
  port: smtp_server_port (default 587)
  sender: sender_email
  token: smtp_token
  receiver: receiver_email(split by comma)
```

send mail
```
configManager, err := config.GetYamlConfigManager("mail_conf.yml")
sender, err := NewSMTPMailSender(configManager)
// AddReceiver is not necessary, defaultly use receiver in config file (mail_conf.yml)
err = sender.Send(AddReceiver("receiver mail account"), SetTitle("test_title"), SetContent("test_content"), SetNickName("nickname"))
```

## alerter
```
# get mail alerter, you have to get mail sender first
alerter := GetMailAlerter(sender, testDefaultReceiver)

# send alert
alerter.Alert(SetAlertTitleAndMsg(testAlertTitle, testAlertMsg))
```

## monitor
monitor manager supply metrics management. The actual implement include prometheus. User is no need to concern about how metrics upload to monitor service

Current support metrics: gauge, counter

```
// get prometheus monitor manager
manager := GetPrometheusMonitorManagerByConf(
	PrometheusMonitorManagerConf{
		Port: http_server_port,
})

// add metrics
metricsDesc := NewMonitorMetrics(Gauge, "test_gauge", "test gauge", LabelKey{"name"})
err := manager.AddMetrics(metricsDesc)

// get metrics
metrics, err := manager.GetMetrics(currentTestCase.metricsName)

// set metrics value, need transform to actual metrics type first ( e.g. metrics to prometheusGauge )
metrics.(*PrometheusGauge).With(MetricsLabel{"name": "smiecj"}).Set(10)
```

## zookeeper client

### config
zk:
  address: zk_node_1:2181,zk_node_2:2181
  home: /

### usage
```
import "github.com/smiecj/go_common/zk"

// get zookeeper client
client, err := zk.GetZKonnector(configManager)

// create node
// default: create persistent node
err = client.Create(SetPath("/persistent"))
// create ephemeral node
err = client.Create(SetPath("/ephemeral"), SetEphemeral(), SetTTL(time.Minute))

// list children
listNodes, err = client.List(SetPath("/parent"))

// delete node
err = client.Delete(SetPath("/parent/child"))

// delete node include children
err = client.DeleteAll(SetPath("/parent"))
```

## other common tool

### time - format
```
import "github.com/smiecj/go_common/util/time"

// get current timestamp (format: 2006-01-02 15:04:05)
time.GetCurrentTimestamp()
```

### time - fixed hour ticker
```
import tickerutil "github.com/smiecj/go_common/util/time/ticker"

// tick on 8 am everyday (between 8:00~9:00)
ticker := tickerutil.NewFixHourTicker(8, tickerutil.SetFunc(func() error {dosomething...}))
ticker.Start()
// handle error chan
// if you set ignore error true (ticker.SetIsIgnoreError), then you don't need to consume error chan
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

// get local ipv4
ip, err := net.GetLocalIPV4()

// check local port is used
isUsed := net.CheckLocalPortIsUsed(22)
```
