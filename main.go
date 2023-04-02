package main

import (
	cache2 "douyin-grab/pkg/cache"
	"douyin-grab/pkg/logger"
	queue2 "douyin-grab/pkg/queue"
	"douyin-grab/web"
	"douyin-grab/wsocket"
)

const (
	VERSION = `0.0.1`
)

func main() {
	logger.Init("")

	// 初始化bigcache
	cache := cache2.NewCache()
	if err := cache.SetDefaultCaChe(); err != nil {
		logger.Error("初始化默认缓存异常: %s", err)
		return
	}

	qu := queue2.NewQueueSrv()

	douYinSrv := wsocket.NewDouYin(qu, cache)

	sendClientSrv := wsocket.NewSendClientSrv(qu)
	sendClientSrv = sendClientSrv
	webSrv := web.NewWeb(cache, douYinSrv)
	// 启动web服务
	webSrv.RunWeb()
	// go web.RunWsClient(qu)

	// app := cli.NewApp()
	// app.Name = `直播间弹幕`
	// app.Version = VERSION
	// app.Before = func(ctx *cli.Context) error {
	// 	return nil
	// }
	//
	// var live_room_url, wss_url string
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:        "live_room_url, lrurl",
	// 		Usage:       "live room url",
	// 		Destination: &live_room_url,
	// 	},
	// 	cli.StringFlag{
	// 		Name:        "wss_url, wssurl",
	// 		Usage:       "live room wws url",
	// 		Destination: &wss_url,
	// 	},
	// }
	//
	// var err error
	// app.Action = func(ctx *cli.Context) error {
	// 	// if len(live_room_url) == 0 {
	// 	// 	live_room_url = constv.LiveUrl // 默认直播间url
	// 	// }
	// 	// // logger.Info("LiveRoomUrl: %s", live_room_url)
	// 	//
	// 	// if len(wss_url) == 0 {
	// 	// 	wss_url = constv.LiveWsUrl // 默认直播间wss_url
	// 	// }
	// 	// // logger.Info("LiveRoomWssUrl: %s", wss_url)
	//
	// 	// // 获取直播间信息
	// 	// _, ttwid := grab.FetchLiveRoomInfo(live_room_url)
	// 	//
	// 	// // 与直播间进行websocket通信，获取评论数据
	// 	// header := http.Header{}
	// 	// header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") // 设置User-Agent头
	// 	// header.Set("Origin", constv.DouYinOrigin)
	// 	// cookie := &http.Cookie{
	// 	// 	Name:  "ttwid",
	// 	// 	Value: ttwid,
	// 	// }
	// 	// header.Add("Cookie", cookie.String())
	// 	// wsclient := wsocket.NewWSClient(qu).SetRequestInfo(wss_url, header)
	// 	// wsclient.ConnWSServer(ttwid)
	// 	// wsclient.RunWSClient()
	//
	// 	// worker服务
	// 	// go nmid.RunWorker()
	//
	// 	// todo 消费队列数据
	// 	// go sendClientSrv.SendStrToDistal()
	// 	// go sendClientSrv.SendStrToLocal()
	//
	// 	return nil
	// }
	//
	// err = app.Run(os.Args)
	// if err != nil {
	// 	panic(err)
	// }

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// os.Exit(0)
}
