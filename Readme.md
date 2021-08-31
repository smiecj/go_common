# go common
提供 go 相关的公共库，其他业务仓库使用
如：http 客户端、公共配置解析、公共数据库基类等

# 已实现功能
## http 客户端
使用方式
```
contentBytes := client.DoGetRequest(url, nil)
```