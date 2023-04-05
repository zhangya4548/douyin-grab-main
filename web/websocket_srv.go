package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// 创建 websocket 升级器
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 保存连接的客户端
	clients = make(map[*websocket.Conn]bool)
	// 发送消息的 channel
	broadcastChannel = make(chan []byte)
)

func (s *Web) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 http 请求为 websocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	// 添加客户端到 clients 映射
	clients[conn] = true
	log.Printf("新本地ws客户端连接: %s", conn.RemoteAddr().String())
	// 启动一个协程，每 5 秒向客户端发送一个 ping 消息，以保持连接
	go sendPingMessage(conn)
	// 循环从连接中读取消息
	for {
		stopState, err := s.cache.Get("Stop")
		if err != nil {
			log.Println("取到缓存异常:", err)
			time.Sleep(time.Second * 2)
			continue
		}
		if stopState == "true" {
			log.Println("停止中")
			s.qu.Empty()
			time.Sleep(time.Second * 10)
			continue
		}

		// 取队列数据
		messageS := s.qu.GetAll()
		if len(messageS) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		for _, message := range messageS {
			if message == "" {
				continue
			}
			log.Println("取到队列弹幕:", message)
			// 将消息放入 broadcastChannel 中，等待广播
			broadcastChannel <- []byte(message)
		}
	}
}

// 每 5 秒向客户端发送一个 ping 消息，以保持连接
func sendPingMessage(conn *websocket.Conn) {
	for {
		time.Sleep(5 * time.Second)
		err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
		if err != nil {
			log.Printf("ping本地ws客户端失败: %v", err)
			// 如果发送 ping 消息失败，则将客户端从 clients 映射中删除
			delete(clients, conn)
			break
		}
	}
}

// 启动一个 goroutine，从 broadcastChannel 中读取消息并发送给所有客户端
func broadcast() {
	for {
		message := <-broadcastChannel
		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("广播消息到本地ws失败: %v", err)
				// 如果发送消息失败，则将客户端从 clients 映射中删除
				delete(clients, conn)
			} else {
				// 发送 PING 帧以保持连接状态
				err = conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
				if err != nil {
					log.Printf("ping本地ws客户端失败: %v", err)
					// 如果发送 ping 消息失败，则将客户端从 clients 映射中删除
					delete(clients, conn)
				}
			}
		}
	}
}

func init() {
	// 启动一个 goroutine，从 broadcastChannel 中读取消息并发送给所有客户端
	go broadcast()
}
