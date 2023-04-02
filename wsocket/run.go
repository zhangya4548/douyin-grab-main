package wsocket

import (
	"douyin-grab/constv"
	"douyin-grab/grab"
	"douyin-grab/pkg/cache"
	queue2 "douyin-grab/pkg/queue"
	"log"
	"net/http"
)

/*
@Time : 2023/4/2 23:03
@Author : zhangguangqiang
@File : run
@Software: GoLand
*/

type DouYinSrv struct {
	qu    *queue2.QueueSrv
	cache *cache.Cache
}

func NewDouYin(qu *queue2.QueueSrv,
	cache *cache.Cache) *DouYinSrv {
	return &DouYinSrv{
		qu:    qu,
		cache: cache,
	}
}

func (d *DouYinSrv) Run() {
	LiveRoomUrl, err := d.cache.Get("LiveRoomUrl")
	if err != nil {
		log.Println("异常:", err.Error())
		return
	}

	WssUrl, err := d.cache.Get("WssUrl")
	if err != nil {
		log.Println("异常:", err.Error())
		return
	}

	// 获取直播间信息
	_, ttwid := grab.FetchLiveRoomInfo(LiveRoomUrl)

	// 与直播间进行websocket通信，获取评论数据
	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") // 设置User-Agent头
	header.Set("Origin", constv.DouYinOrigin)
	cookie := &http.Cookie{
		Name:  "ttwid",
		Value: ttwid,
	}
	header.Add("Cookie", cookie.String())
	wsClient := NewWSClient(d.qu).SetRequestInfo(WssUrl, header)
	wsClient.ConnWSServer(ttwid)
	wsClient.RunWSClient()

	// // 检测是否要停止
	// go func() {
	// 	for {
	// 		duration := constv.DefaultHeartbeatTime
	// 		timer := time.NewTimer(duration)
	// 		<-timer.C
	// 		Stop, err := d.cache.Get("Stop")
	// 		if err != nil {
	// 			log.Println("异常:", err.Error())
	// 			return
	// 		}
	// 		if Stop == "true" {
	// 			wsClient.Close()
	// 		}
	// 	}
	// }()
}
