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
	qu *queue2.EsQueue
}

func NewSendClientSrv(qu *queue2.EsQueue) *SendClientSrv {
	return &SendClientSrv{qu: qu}
}

func (s *SendClientSrv) SendStr() {
	// 定义websocket客户端
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

	log.Println("连接上远程wsocket服务端")
	for {
		jsonStrS := make([]interface{}, 1000)
		gets, quantity := s.qu.GetAll(jsonStrS)
		if gets == 0 {
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Printf("获取到队列数据: %d, 队列剩余:%d \n", gets, quantity)
		for _, v := range jsonStrS {
			if v == nil {
				continue
			}
			// val, ok, quantity := s.qu.Get()
			// if !ok {
			// 	time.Sleep(time.Second * 10)
			// 	continue
			// }
			// fmt.Printf("获取到队列数据: %s, 队列剩余:%v \n", val.(string), quantity)

			jsonStr := v.(string)
			if jsonStr == "" {
				continue
			}
			if err := c.WriteMessage(websocket.TextMessage, []byte(jsonStr)); err != nil {
				fmt.Println("推送数据到服务端异常:", err)
				time.Sleep(time.Second * 2)
				continue
			}
			fmt.Println("推送数据到服务端完:", jsonStr)
		}

	}
}
