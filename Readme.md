# go common

提供 go 相关的公共库，其他业务仓库使用
如：http 客户端、公共配置解析、公共数据库基类等

# 已实现功能
## http 客户端

```
import "github.com/smiecj/go_common/http"
client := GetHTTPClient()
client.Do(Get(), Url("http://..."))
contentBytes := client.DoGetRequest(url, nil)
```

## file writer

```
import "github.com/smiecj/go_common/util/file"
file.Write("/tmp/test.log", "content", file.ModeCreate)
```

## logger

```
import log "github.com/smiecj/go_common/util/log"
log.Info("msg: %s", msg)
```

## 错误码
```
import log "github.com/smiecj/go_common/errorcode"
return errorcode.ServiceError
```

## DB 操作
### 待补充