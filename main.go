package main

import (
	"douyin-grab/pkg/cache"
	"douyin-grab/pkg/logger"
	queue2 "douyin-grab/pkg/queue"
	"douyin-grab/web"
	"douyin-grab/wsocket"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

const (
	VERSION = `0.0.1`
)

func main() {
	// 缓存
	caChe := cache.NewCache()
	err := caChe.SetDefaultCaChe()
	if err != nil {
		fmt.Println("初始化默认cache异常:", err)
		return
	}
	// 队列
	qu := queue2.NewQueueSrv()

	// 抖音wsocket客户端
	wsDouYinClient := wsocket.NewWSClient(qu, caChe)

	// 远程wsocket客户端
	sendClientSrv := wsocket.NewSendClientSrv(qu)
	sendClientSrv = sendClientSrv
	// go sendClientSrv.SendStr() todo 开启

	// 本地web服务
	webSrv := web.NewWeb(caChe, wsDouYinClient)
	go webSrv.RunWeb()

	app := cli.NewApp()
	app.Name = `抖音弹幕采集`
	app.Version = VERSION
	app.Before = func(ctx *cli.Context) error {
		return nil
	}

	logger.Init("")
	var live_room_url, wss_url string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "live_room_url, lrurl",
			Usage:       "live room url",
			Destination: &live_room_url,
		},
		cli.StringFlag{
			Name:        "wss_url, wssurl",
			Usage:       "live room wws url",
			Destination: &wss_url,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		// if len(live_room_url) == 0 {
		// 	live_room_url = constv.DEFAULTLIVEROOMURL // 默认直播间url
		// }
		// logger.Info("live room url: %s", live_room_url)
		//
		// if len(wss_url) == 0 {
		// 	wss_url = constv.DEFAULTLIVEWSSURL // 默认直播间wss_url
		// }
		// logger.Info("live room wss_url: %s", wss_url)

		// // 获取直播间信息
		// _, ttwid := grab.FetchLiveRoomInfo(live_room_url)
		//
		// // 与直播间进行websocket通信，获取评论数据
		// header := http.Header{}
		// header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") // 设置User-Agent头
		// header.Set("Origin", constv.DOUYIORIGIN)
		// cookie := &http.Cookie{
		// 	Name:  "ttwid",
		// 	Value: ttwid,
		// }
		// header.Add("Cookie", cookie.String())
		// wsclient := wsocket.NewWSClient(constv.DEFAULTLIVEROOMURL, constv.DEFAULTLIVEWSSURL, qu)
		// wsclient.Run()

		// worker服务
		// go nmid.RunWorker()

		// 消费队列数据
		// go sendClientSrv.SendStr()

		// go func() {
		// 	time.Sleep(time.Second * 30)
		// 	wsclient.Close()
		// 	fmt.Println("停止了")
		// }()

		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	os.Exit(0)
}
