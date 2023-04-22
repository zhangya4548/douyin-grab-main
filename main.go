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
)

func main() {
	caChe := cache.NewCache()
	err := caChe.SetDefaultCaChe()
	if err != nil {
		fmt.Println("初始化默认cache异常:", err)
		return
	}
	qu := queue2.NewQueueSrv()

	wsDouYinClient := wsocket.NewWSClient(qu, caChe)

	sendClientSrv := wsocket.NewSendClientSrv(qu)
	go sendClientSrv.SendStr()

	webSrv := web.NewWeb(qu, caChe, wsDouYinClient)
	go webSrv.RunWeb()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	os.Exit(0)
}
