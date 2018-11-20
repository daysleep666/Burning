启动服务器

go run server/server.go port

启动客户端读屏

go run clientreader/client.go ip:port

启动 客户端写屏

go run clientwrite/client.go ip:port

例子:
go run clientreader/client.go 10.0.0.66:1234
go run clientwrite/client.go 10.0.0.66:1234