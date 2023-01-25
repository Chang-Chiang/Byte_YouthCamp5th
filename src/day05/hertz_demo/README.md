# hertz_demo

> + [快速开始 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/getting-started/)

## 1. 安装 hertz
```bash
$ go install github.com/cloudwego/hertz/cmd/hz@latest
```

```bash
$ hz -v
hz version v0.5.1
```

## 2. 新建项目目录

```bash
$ mkdir hertz_demo
$ cd hertz_demo
```

## 3. 生成代码

```bash
$ hz new -module=hertz_demo
```

## 4. 整理 & 拉取依赖

```bash
$ go mod tidy
```

## 5. 编译

```bash
$ go build -o hertz_demo
```

## 6. 运行

```bash
$ ./hertz_demo
```

## 7. 测试

```bash
$ curl http://127.0.0.1:8888/ping
{"message":"pong"}
```



