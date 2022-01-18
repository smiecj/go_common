[中文](https://github.com/smiecj/go_common/blob/master/Readme_zh.md)

# go common

go common library, supply some common library for other repo to use

e.g. http client, config manager and db connector etc.

# implemented features
## http client

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
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName))
// search object
SearchRet, err := localConnector.Search(SearchSetSpace(dbName, tableName), SearchSetObject(testStruct{}), SearchSetObjectArrType([]*testStruct{}))
// notice: local file connector use reflect lib to unmarshal object everyline, so please use pointer array when set object array type
```

### mysql (gorm)
```
connector, err := GetMySQLConnector(MySQLConnectOption{Host: host, Port: port, User: user, Password: password})

// store
insertRet, err := connector.Insert(InsertSetSpace("db_name", "table_name"), 
    InsertAddObjectArr([]objectArr{obj1, obj2}), InsertSetObjectArrType([]object{}))
// notice: you can set an empty slice when call InsertSetObjectArrType, SearchSetObjectArrType, when call reflect lib, the array will automatically init
// notice: objectArrType is needed because []interface{} cannot be recognized in gorm

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

// deletr
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
config, err := config.GetYamlConfig(config_file_name)

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
```
sender := NewQQMailSender(MailSenderConf{
		Token:  "qq mail token",
		Sender: "sender qq mail account",
	})
err := sender.Send(AddReceiver("receiver mail account"), SetTitle("test_title"), SetContent("test_content"), SetNickName("nickname"))
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
manager := GetPrometheusMonitorManager(
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

# todo
## RPC interface

## alert convergence

## alert set default receiver