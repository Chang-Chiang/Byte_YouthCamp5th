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
