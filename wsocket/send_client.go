package wsocket

import (
	queue2 "douyin-grab/pkg/queue"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type SendClientSrv struct {
	qu *queue2.QueueSrv
}

func NewSendClientSrv(qu *queue2.QueueSrv) *SendClientSrv {
	return &SendClientSrv{qu: qu}
}

func (s *SendClientSrv) SendStr() {
	// u := url.URL{Scheme: "ws", Host: "lwww.wykji.cn:53331", Path: "/wss/dan/mu/conn"}
	wsRemoteHost := os.Getenv("WsRemoteHost")
	wsRemotePath := os.Getenv("WsRemotePath")
	u := url.URL{Scheme: "ws", Host: wsRemoteHost, Path: wsRemotePath}
	header := make(http.Header)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
		return
	}
	defer c.Close()

	for {
		jsonStr := s.qu.Pop()
		if jsonStr == "" {
			time.Sleep(time.Second * 5)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, []byte(jsonStr)); err != nil {
			fmt.Println("推送数据到服务端异常:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		fmt.Println("推送数据到服务端完:", jsonStr)
	}
}
