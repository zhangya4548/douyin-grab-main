package main

import (
	"douyin-grab/pkg/cache"
	queue2 "douyin-grab/pkg/queue"
	"douyin-grab/web"
	"douyin-grab/wsocket"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("读取.env文件异常:", err)
		return
	}
	fmt.Println("远程配置", os.Getenv("WsRemoteHost"), os.Getenv("WsRemotePath"))

	caChe := cache.NewCache()
	err = caChe.SetDefaultCaChe()
	if err != nil {
		fmt.Println("初始化默认cache异常:", err)
		return
	}
	qu := queue2.NewQueueSrv2(1000)

	wsDouYinClient := wsocket.NewWSClient(qu, caChe)

	sendClientSrv := wsocket.NewSendClientSrv(qu)
	sendClientSrv = sendClientSrv
	go sendClientSrv.SendStr() // todo 开启

	webSrv := web.NewWeb(qu, caChe, wsDouYinClient)
	go webSrv.RunWeb()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	os.Exit(0)
}
