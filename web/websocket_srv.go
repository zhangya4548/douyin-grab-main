package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients          = make(map[*websocket.Conn]bool)
	broadcastChannel = make(chan []byte)
)

func (s *Web) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	clients[conn] = true
	log.Printf("新本地ws客户端连接: %s", conn.RemoteAddr().String())
	go sendPingMessage(conn)
	for {
		stopState, err := s.cache.Get("Stop")
		if err != nil {
			log.Println("取到缓存异常:", err)
			time.Sleep(time.Second * 2)
			continue
		}
		if stopState == "true" {
			log.Println("停止中")
			time.Sleep(time.Second * 10)
			continue
		}

		messageS := []string{}
		if len(messageS) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		for _, message := range messageS {
			if message == "" {
				continue
			}
			log.Println("取到队列弹幕:", message)
			broadcastChannel <- []byte(message)
		}
	}
}

func sendPingMessage(conn *websocket.Conn) {
	for {
		time.Sleep(5 * time.Second)
		err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
		if err != nil {
			log.Printf("ping本地ws客户端失败: %v", err)
			delete(clients, conn)
			break
		}
	}
}

func broadcast() {
	for {
		message := <-broadcastChannel
		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("广播消息到本地ws失败: %v", err)
				delete(clients, conn)
			} else {
				err = conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
				if err != nil {
					log.Printf("ping本地ws客户端失败: %v", err)
					delete(clients, conn)
				}
			}
		}
	}
}

func init() {
	go broadcast()
}
