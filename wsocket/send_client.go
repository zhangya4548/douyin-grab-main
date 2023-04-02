// @Time : 2023/3/30 7:28 PM
// @Author : zhangguangqiang
// @File : send_client
// @Software: GoLand

package wsocket

import (
	"douyin-grab/constv"
	queue2 "douyin-grab/pkg/queue"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type SendClientSrv struct {
	qu *queue2.QueueSrv
}

func NewSendClientSrv(qu *queue2.QueueSrv) *SendClientSrv {
	return &SendClientSrv{qu: qu}
}

func (s *SendClientSrv) SendStrToDistal() {
	// 定义websocket客户端
	u := url.URL{Scheme: "ws", Host: constv.WsClientPort, Path: constv.DistalWsPath}
	header := make(http.Header)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Println("连接远端Ws异常:", err)
		return
	}
	defer c.Close()

	for {
		jsonStr := s.qu.Pop()
		if jsonStr == "" {
			time.Sleep(time.Second * 2)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, []byte(jsonStr)); err != nil {
			log.Println("推送数据到远端Ws服务端异常:", err)
			time.Sleep(time.Second * 2)
			continue
		}
		log.Println("推送数据到远端Ws服务端完:", jsonStr)
	}
}

func (s *SendClientSrv) SendStrToLocal() {
	// 定义websocket客户端
	u := url.URL{Scheme: "ws", Host: constv.WsClientPort, Path: "/ws"}
	header := make(http.Header)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Println("连接本地Ws异常:", err)
		return
	}
	defer c.Close()

	for {
		jsonStr := s.qu.Pop()
		if jsonStr == "" {
			time.Sleep(time.Second * 2)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, []byte(jsonStr)); err != nil {
			log.Println("推送数据到本地Ws服务端异常:", err)
			time.Sleep(time.Second * 2)
			continue
		}
		log.Println("推送数据到本地Ws服务端完:", jsonStr)
	}
}
