这是我参与「第五届青训营 」笔记创作活动的第5天

# 一、内容概览

+ Web 框架
+ RPC 框架
+ ORM 框架

# 二、知识点详解

## 1. 简介

+ Gorm: ORM 框架
+ Kitex: 字节内部 Golang 微服务 RPC 框架，高性能、强可扩展、支持多协议且拥有丰富的开源扩展
+ Hertz: 字节内部 HTTP 框架，高易用性、高性能、高扩展性

## 2. Gorm

### 2.1 基本使用

> [模型定义 | GORM - The fantastic ORM library for Golang, aims to be developer friendly.](http://gorm.cn/zh_CN/docs/models.html)

```go
// 定义 gorm model
type Product struct {
    Code  string
    Price uint
}

// 为 model 定义表名
func (p Product) TableName() string {
    return "product"
}

func main() {
    
    // 连接数据库
    db, err := gorm.Open(
        mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"),
        &gorm.Config{})
    if err != nil {
        panic("failed to connnect database")
    }
    
    // 创建数据
    db.Create(&Product{Code: "D42", Price: 100})
    
    // 查询数据
    vat product Product
    db.First(&product, 1)  // 根据整型主键查找
    db.First(&product, "code = ?", "D42")  // 查找 code 字段值为 D42 的记录
    
    // 更新数据，product 的 price 更新为 200
    db.Model(&product).Update(column:"Price", value：200)
    
    // 更新数据，更新多个字段
    db.Model(&product).Updates(Product{Price：200，Code:"F42"})  // 仅更新非零值字段
    db.Model(&product).Updates(map[string]interface{}{"Price": 200，"Code": "F42"})
    
    // 删除数据
    db.Delete(&product, 1)
}
```

### 2.2 支持的数据库

> [连接到数据库](https://gorm.cn/zh_CN/docs/connecting_to_the_database.html)

GORM 通过驱动来连接数据库，如果需要连接其他类型的数据库，可以复用/自行开发驱动

+ MySQL

  > DSN: https://github.com/go-sql-driver/mysql#dsn-data-source-name

  ```go
  import(
      "gorm.io/driver/sqlserver"
      "gorm.io/gorm"
  )
      
  // 连接SQLServer数据库为例
  // github.com/denisenkom/go-mssqldb 
  dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
  db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config})
  ```

+ SQLServer

+ PostgreSQL

+ SQLite

### 2.3 创建数据

> [创建](https://gorm.cn/zh_CN/docs/create.html)

```go
package main 

import(
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Product struct {
    ID uint `gorm:"primarykey"`
    Code string `gorm:"column: code"`
    Price uint  `gorm:"column: user_id"`
}

func main() {
    db, err := gorm.Open(
        mysql.Open(dsn:"username:password@tcp(localhost:9918)/database?charset=utf8"),
        &gorm.Config{})
    if err != nil : "failed to connect database" *
    // 创建一条
    p := &Product{Code: "D42"}
    res := db.Create(p)
    fmt.Println(res.Error)  // 获取err 
    fmt.Println(p.ID)  //返回插入数据的主键
    // 创建多条
    products := []*Product{{Code: "D41"}, {Code: "D42"}, {Code: "D43"}}
    res = db.Create(products)
    fmt.Println(res.Error)  // 获取err 
    for _, p := range products {
    	fmt.Println(p.ID)
    }
}
```

+ Upsert 使用

  ```go
  // 以不处理冲突为例，创建一条数据
  p := &Product{Code: "D42", ID: 1}
  db.Clauses(clause.OnConflict{DoNothing: true}.Create(&p))  // clause.OnConflict 处理数据冲突
  ```

+ 默认值

  ```go
  // 使用 default 标签为字段定义默认值
  type User struct {
      ID int64
      Name string `gorm:"default:galeone"`
      Age int64 `gorm:"default:18"`
  }
  ```

### 2.4 查询数据

> [查询](https://gorm.cn/zh_CN/docs/query.html)

```go
db,err:=gorm.Open(
	mysql.Open(dsn:"username:passwordetcp(localhost:9910)/database?charset=utf8"),
	&gorm.Config{})
if err != nil {
	panic("failed to connect database")
}

// 获取第一条记录(主键升序), 查询不到数据则返回 ErrRecordNotFound 
u := &User{}
db.First(u)  //SELECT*FROM users ORDER BY id LIMIT 1；

// 查询多条数据
users := make([]*User,0)
resuLt := db.Where("age>10").Find(&users)  // SELECT*FROM Users where age>10;
fmt.Println(result.RowsAffected)  // 返回找到的记录数, 相当于 `len(users)`
fmt.Println(result.Error)  // returns error

// IN SELECT * FROM users WHERE name IN('jinzhu', 'jinzhu 2');
db.Where("name IN ？", []string{"jinzhu", "jinzhu 2"}).Find(&users)

// LIKE SELECT * FROM users WHERE name LIKE‘%jin%';
db.Where("name LIKE ?", "%jin%").Find(&users)

// AND SELECT * FROM users WHERE name='iinzhu'AND age >=22;
db.Where("name= ? AND age >= ?", "jinzhu", "22").Find(&users)

// SELECT * FROM users WHERE name="jinzhu";
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)

// SELECT * FROM users WHERE name = "jinzhu" AND age =0;
db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
```

+ 使用 First 查询

  + 查询不到数据会返回 ErrRecordNotFound
  + 查询多条数据，查询不到数据不会返回错误

+ 使用结构体作为查询条件

  GORM 只会查询非零值字段，字段值为 "0"、"false" 或其他 零值，该字段不会被用于构建查询条件 → 使用 Map 来构建查询条件

### 2.5 更新数据

> [更新](https://gorm.cn/zh_CN/docs/update.html)

```go
// 条件更新单个列
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE age > 18;
db.Model(&User{ID:111}).Where("age>?", 18).Update("name", "hello")

// 更新多个列
// 根据 `struct` 更新属性, 只会更新非零值的字段
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111
db.Model(&User{ID:111}).Updates(User{Name:"hello", Age:18})

// 根据 `map` 更新属性
// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
db.Model(&User{ID:111}).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})

// 更新选定字段
// UPDATE users SET name='hello' WHERE id=111;
db.Model(&User{ID:111}).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})

// SQL 表达式更新
// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;
db.Model(&User{ID:111}).Update("age", gorm.Expr(expr:"age * ? + ?", 2, 100))
```

使用 Struct 更新时，只会更新非零值，如果需要更新零值可以使用 Map 更新或使用 Select 选择字段

### 2.6 删除数据

> [删除](https://gorm.cn/zh_CN/docs/delete.html)

+ 物理删除

  ```go
  db.Delete(& User{}, conds.:10) // DELETE FROM users WHERE id=10;
  db.Delete(& User{}, conds.:"10") // DELETE FROM users WHERE id=10;
  db.Delete(& User{}, [] int{1, 2, 3}) //DELETE FROM users WHERE id IN(1,2,3);
  db.Where(query:"name LIKE ?", "% jinzhu%").Delete(User{}) // DELETE from users where name LIKE "% iinzhu%";
  db.Delete(User{}, "email LIKE ?", "% jinzhu%") // DELETE from users where name LIKE "% jinzhu%";
  ```

+ 软删除

  + GORM 提供了 gorm.DeletedAt 用于帮助用户实现软删
  + 拥有软删除能力的 Model 调用 Delete 时，记录不会被从数据库中真正删除。但 GORM 会将 DeletedAt 置为当前时间，并且不能再通过正常的查询方法找到该记录
  + 使用 Unscoped 可以查询到被软删的数据

  ```go
  package main
  import (
  	"gorm.io/driver/mysqL"
  	"gorm.io/gorm"
  )
  type User struct{
  	ID int64
  	Name string `gorm:"defauLt:galeone"`
  	Age int64 `gorm:"default:18"`
  	Deleted gorm.DeletedAt
  }
  
  func main(){
  	db, err := gorm.Open(
  	mysql.Open("username:password@tcp(localhost:9910)/database?charset=utf8"),
  	&gorm.Config{})
  	if err != nil {
  	panic("failed to connect database")
  	}
  	
  	// 翻除一条
  	u := User{ID: 111}  // user 的 ID 是 `111`
  	db.Delete(&u)       // UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;
  	
  	// 批量删除
  	db.Where("age = ?", 20).Delete(&User{}) // UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;
      users := make([]*User,0)
  	// 在查询时会忽略被软删除的记录
  	db.Where("age = 20").Find(&users) // SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
  	// 在查询时不会忽略被软删除的记录
  	db.Unscoped().Where("age =20").Find(&users) // SELECT * FRON users WHERE age = 20;
  }
  ```

### 2.7 事务

> [事务](https://gorm.cn/zh_CN/docs/transactions.html)

+ Gorm 提供了 Begin、Commit、Rollback 方法用于使用事务
+ Gorm 提供了 Transaction 方法用于自动提交事务，避免用户漏写 Commit、Rollback

```go
···
```

### 2.8 HOOK

> [Hook](https://gorm.cn/zh_CN/docs/hooks.html)

+ GORM 提供了 CURD 的 Hook 能力
+ Hook 是在创建、查询、更新、删除等操作之前、之后 自动调用的函数
+ 如果任何 Hook 返回错误，GORM 将停止后续的操作并回滚事务

```go
···
```

### 2.9 性能提高

> [性能](https://gorm.cn/zh_CN/docs/performance.html)

+ 对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，可以使用 SkipDefaultTransaction 关闭默认事务

+ 使用 PrepareStmt 缓存预编译语句可以提高后续调用的速度

```go
···
```

### 2.10 生态

+ [代码生成工具](https://github.com/go-gorm/gen)
+ [分片库方案](https://github.com/go-gorm/sharding)
+ [手动索引](https://github.com/go-gorm/hints)
+ [乐观锁](https://github.com/go-gorm/optimisticlock)
+ [读写分离](https://github.com/go-gorm/dbresolver)
+ [OpenTelemetry 扩展](https://github.com/go-gorm/opentelemetry)

## 3. Kitex

### 3.1 安装

```bash
$ vi ~/.bashrc
# GOPATH=/home/chang/Projects/go
# $GOPATH/bin 添加到环境变量, 因为 kitex 是安装到该路径
export PATH=$PATH:/usr/local/go/bin:/home/chang/Projects/go/bin
$ source ~/.bashrc

$ go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
$ go install github.com/cloudwego/thriftgo@latest
$ kitex --version
v0.4.4
```

### 3.2 定义 IDL

> + Thrift: https://thrift.apache.org/docs/idl
> + Proto3:https://developers.google.com/protocol-buffers/docs/proto3

如果我们要进行 RPC，就需要知道对方的接口是什么，需要传什么参数，同时也需要知道返回值是什么样的。这时候，就需要通过 IDL 来约定双方的协议，就像在写代码的时候需要调用某个函数，我们需要知道函数签名一样。

```go
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

### 3.3 生成代码

`kitex -module example -service example echo.thrift`

```angular
.        
├── build.sh              [构建脚本]     
|
├── echo.thrift         
|
├── handler.go            [IDL service 定义的方法] 
|
├── kitex_gen             [IDL 内容相关的生成代码，主要是基础的 Server/Client 代码] 
|   └── api  
|       ├── echo     
|       |   ├── client.go  
|       |   ├── echo.go   
|       |   ├── invoker.go   
|       |   └── server.go 
|       ├── echo.go
|       └── k-echo.go
| 
├── main.go               [程序入口]
|
└── script     
    ├── bootstrap.sh     
    └── setting.py
```

### 3.4 基本使用

```go
package main

import (
	"context"
	"example/kitex_gen/api"
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the EchoImpl interface.
func (s * EchoImpl) Echo(ctx context.Context, req *api.Request)(resp *api.Response, err error){
	// TODO: Your code here...
	return
}
```

### 3.5 Client 发起请求

+ 创建 Client

  ```go
  import "example/kitex_gen/api/echo"
  import "github. com/cloudwego/kitex/client"
  ···
  c, err := echo.NewClient("example", client. WithHostPorts("0.0.0.0:8888"))
  if err != nil {
  	log.Fatal(err)
  }
  ```

+ 发起请求

  ```go
  import "example/kitex_gen/api"
  ···
  req := &api.Request{Message: "my request"}
  resp, err := c.Echo(context.Background(), req, callopt.WithRPCTimeout(3 * time.Second))
  if err != nil {
  	log.Fatal(err)
  }
  log.Println(resp)
  ```

### 3.6 服务注册与发现

```go
type HelloImpl struct{}

func (h *HelloImp1) Echo(ctx context.Context, req *api.Request)(resp *api.Response, err error){
	resp = & api.Response{
		Message: req.Message,
	}
	return
}
func main(){
	r, err := etcd.NewEtcdRegistry([] string{"127.0.0.1:2379"})
	if errl != nil {
		log.Fatal(err)
		}
	server := hello.NewServer(new(HeltoImp1), server.withRegistry(r), server.withServerBasicInfo(&rpcinfo.EndpointBasicInfo{
	ServiceName:"Hel1o",
	}))
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

### 3.7 生态

> [Kitex | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex/)

+ [XDS 扩展](https://github.com/kitex-contrib/xds)
+ [openelementary 扩展](https://github.com/kitex-contrib/obs-opentelemetry)
+ [ETCD 服务注册与发现扩展](https://github.com/kitex-contrib/registry-etcd)
+ [Nacos 服务注册与发现扩展](https://github.com/kitex-contrib/registry-nacos)
+ [Zookeeper 服务注册与发现扩展](https://github.com/kitex-contrib/registry-zookeeper)
+ [polaris 扩展](https://github.com/kitex-contrib/polaris)
+ [示例代码与业务 Demo](https://github.com/cloudwego/kitex-examples/)

## 4. Hertz

> [字节跳动开源 Go HTTP 框架 Hertz 设计实践 | CloudWeGo](https://www.cloudwego.io/zh/blog/2022/06/21/字节跳动开源-go-http-框架-hertz-设计实践/)

### 4.1 基本使用

> [快速开始 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/getting-started/)

### 4.2 路由

> [路由 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/route/)

### 4.3 参数绑定

> [绑定与校验 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/binding-and-validate/)

### 4.4 中间件

> [中间件概览 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/middleware/)

### 4.5 Client

> [代码示例 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/example/#client)

### 4.6 代码生成工具

> [hz 命令行工具使用 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/toolkit/)

### 4.7 性能

+ 网络库 Netpoll
+ Json 编解码 Sonic
+ sync.Pool 复用对象协议层数据解析优化

### 4.8 生态

> [Hertz | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/#extensions)

# 三、实践——笔记项目

> + [普通版 easy_note](https://github.com/cloudwego/kitex-examples/tree/main/bizdemo/easy_note)
> + [优化版 easy_note](https://github.com/cloudwego/biz-demo/tree/main/easy_note)
> + [‌﻿⁣⁢⁡⁣﻿‌⁢‬﻿⁣‍‬‌‌‍‬⁤⁣⁣‬‬‍﻿⁣⁡⁢⁣⁤⁤﻿﻿‬﻿⁢‌⁢⁡‬‌⁡‬实战项目-Go 语言笔记服务 - Feishu Docs](https://bytedance.feishu.cn/docx/Wwa4dfwScogfjLxclXKcStGEncd)

## 1. 项目介绍

使用 Hertz、Kitex、Gorm 搭建具备一定业务逻辑的后端 API 项目

| 服务名称 |   服务介绍   | 传输协议 |    主要技术栈    |
| :------: | :----------: | :------: | :--------------: |
| demoapi  |   API 服务   |   HTTP   | Gorm/Kitex/Hertz |
| demouser | 用户数据管理 | Protobuf |    Gorm/Kitex    |
| demonote | 笔记数据管理 |  Thrift  |    Gorm/Kitex    |

## 2. 项目功能

![image-20230120230717410](./assets/image-20230120230717410.png)

## 3. 项目调用关系

![image-20230120230802847](./assets/image-20230120230802847.png)

## 4. IDL 介绍

```idl
message User {
    int64 user_id = 1;
    string user_name = 2;
    string avatar = 3;
}

message CreateUserRequest {
    string user_name = 1;
    string password = 2;
}

message CreateUserResponse {
    BaseResp base_resp = 1;
}

message MGetUserRequest {
    repeated int64 user_ids = 1;
}

message MGetUserResponse {
    repeated User users = 1;
    BaseResp base_resp = 2;
}

message CheckUserRequest{
    string user_name = 1;
    string password = 2;
}

message CheckUserResponse{
    int64 user_id = 1;
    BaseResp base_resp = 2;
}

service UserService {
    rpc CreateUser (CreateUserRequest) returns (CreateUserResponse) {}
    rpc MGetUser (MGetUserRequest) returns (MGetUserResponse) {}
    rpc CheckUser (CheckUserRequest) returns (CheckUserResponse) {}
}
```

```idl
struct Note {
    1:i64 note_id
    2:i64 user_id
    3:string user_name
    4:string user_avatar
    5:string title
    6:string content
    7:i64 create_time
}

struct CreateNoteRequest {
    1:string title
    2:string content
    3:i64 user_id
}

struct CreateNoteResponse {
    1:BaseResp base_resp
}

struct DeleteNoteRequest {
    1:i64 note_id
    2:i64 user_id
}

struct DeleteNoteResponse {
    1:BaseResp base_resp
}

struct UpdateNoteRequest {
    1:i64 note_id
    2:i64 user_id
    3:optional string title
    4:optional string content
}

struct UpdateNoteResponse {
    1:BaseResp base_resp
}

struct MGetNoteRequest {
    1:list<i64> note_ids
}

struct MGetNoteResponse {
    1:list<Note> notes
    2:BaseResp base_resp
}

struct QueryNoteRequest {
    1:i64 user_id
    2:optional string search_key
    3:i64 offset
    4:i64 limit
}

struct QueryNoteResponse {
    1:list<Note> notes
    2:i64 total
    3:BaseResp base_resp
}

service NoteService {
    CreateNoteResponse CreateNote(1:CreateNoteRequest req)
    MGetNoteResponse MGetNote(1:MGetNoteRequest req)
    DeleteNoteResponse DeleteNote(1:DeleteNoteRequest req)
    QueryNoteResponse QueryNote(1:QueryNoteRequest req)
    UpdateNoteResponse UpdateNote(1:UpdateNoteRequest req)
}
```

## 5. 技术栈

![image-20230121100427433](./assets/image-20230121100427433.png)

## 6. 关键代码

### 6.1 Hertz

```go
// ./cmd/api/handlers/create_note.go
// CreateNote create note info
func CreateNote(ctx context.Context, c *app.RequestContext) {
	var noteVar NoteParam
	if err := c.Bind(&noteVar); err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	if len(noteVar.Title) == 0 || len(noteVar.Content) == 0 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}

	claims := jwt.ExtractClaims(ctx, c)
	userID := int64(claims[constants.IdentityKey].(float64))
	err := rpc.CreateNote(context.Background(), &notedemo.CreateNoteRequest{
		UserId:  userID,
		Content: noteVar.Content, Title: noteVar.Title,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, nil)
}
```

### 6.2 Kitex Client

```go
// ./cmd/api/rpc/note.go

func initNoteRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := noteservice.NewClient(
		constants.NoteServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	noteClient = c
}

// CreateNote create note info
func CreateNote(ctx context.Context, req *notedemo.CreateNoteRequest) error {
	resp, err := noteClient.CreateNote(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}
```

### 6.3 Kitex Server

```go
// ./cmd/note/service/create_note.go

type CreateNoteService struct {
	ctx context.Context
}

// NewCreateNoteService new CreateNoteService
func NewCreateNoteService(ctx context.Context) *CreateNoteService {
	return &CreateNoteService{ctx: ctx}
}

// CreateNote create note info
func (s *CreateNoteService) CreateNote(req *notedemo.CreateNoteRequest) error {
	noteModel := &db.Note{
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
	}
	return db.CreateNote(s.ctx, []*db.Note{noteModel})
}
```

### 6.4 Gorm

```go
// ./cmd/note/dal/db/note.go

type Note struct {
	gorm.Model
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (n *Note) TableName() string {
	return constants.NoteTableName
}

// CreateNote create note info
func CreateNote(ctx context.Context, notes []*Note) error {
	if err := DB.WithContext(ctx).Create(notes).Error; err != nil {
		return err
	}
	return nil
}

// MGetNotes multiple get list of note info
func MGetNotes(ctx context.Context, noteIDs []int64) ([]*Note, error) {
	var res []*Note
	if len(noteIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", noteIDs).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// UpdateNote update note info
func UpdateNote(ctx context.Context, noteID, userID int64, title, content *string) error {
	params := map[string]interface{}{}
	if title != nil {
		params["title"] = *title
	}
	if content != nil {
		params["content"] = *content
	}
	return DB.WithContext(ctx).Model(&Note{}).Where("id = ? and user_id = ?", noteID, userID).
		Updates(params).Error
}

// DeleteNote delete note info
func DeleteNote(ctx context.Context, noteID, userID int64) error {
	return DB.WithContext(ctx).Where("id = ? and user_id = ? ", noteID, userID).Delete(&Note{}).Error
}
```

# 四、总结

+ Gorm 框架
+ Kitex 框架
+ Hertz 框架
+ 笔记项目

# 五、参考

+ [GORM 指南](https://gorm.cn/zh_CN/docs)
+ [Kitex | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex/)
+ [Hertz | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/)
