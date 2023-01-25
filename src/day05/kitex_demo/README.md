# kitex_demo

> + [快速开始 | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex/getting-started/#基础教程)
> + [代码生成工具 | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex/tutorials/code-gen/code_generation/)
> + [Apache Thrift - Interface Description Language (IDL)](https://thrift.apache.org/docs/idl)
> + [Language Guide (proto3)  | Protocol Buffers  | Google Developers](https://developers.google.com/protocol-buffers/docs/proto3)

## 1. 安装 kitex

```bash
$ go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

## 2. 新建项目目录

```bash
$ mkdir kitex_demo
$ cd kitex_demo
$ go mod init kitex_demo
```

## 3. 编写 IDL

```idl
; echo.thrift

namespace go api

struct Request {
  1: string message
}

struct Response {
  1: string message
}

service Echo {
    Response echo(1: Request req)
}
```

## 4. 生成 echo 服务代码

```bash
# kitex -module module_name -service service_name path_to_your_idl.thrift
$ kitex -module kitex_demo -service kitex_demo echo.thrift
```

## 5. 获取最新的 kitex 框架 & 整理、拉取依赖

```bash
$ go get github.com/cloudwego/kitex@latest
$ go mod tidy
```

## 6. 编写 echo 服务逻辑

```go
// handler.go

package main

import (
	"context"
	api "kitex_demo/kitex_gen/api"
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the EchoImpl interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	return &api.Response{Message: req.Message}, nil
}
```

## 7. 编译运行

```bash
$ sh build.sh
$ sh output/bootstrap.sh
# 执行上述命令后，Echo 服务就开始运行
```

## 8. 编写客户端

```bash
# 创建客户端目录
$ mkdir client
$ cd client

# 创建客户端代码 main.go
```

```go
// main.go  client
package main

import (
	"context"
	"kitex_demo/kitex_gen/api"
	"kitex_demo/kitex_gen/api/echo"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
)

func main() {
    
    // 创建客户端
	client, err := echo.NewClient("kitex_demo", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}

	// 每秒发起一次调用
	for {
		req := &api.Request{Message: "my request"}
		resp, err := client.Echo(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}
```

## 9. 客户端运行，发起调用

```bash
$ go run main.go
```

