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
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("从本地ws客户端读取时出错: %v", err)
			// 如果读取消息失败，则将客户端从 clients 映射中删除
			delete(clients, conn)
			break
		}

		log.Printf("收到本地ws客户端消息 %s: %s", conn.RemoteAddr().String(), string(message))
		// 广播接收到的消息给所有客户端
		go broadcastMessage(message)
	}
}

// 广播消息给所有客户端
func broadcastMessage(message []byte) {
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("广播消息到本地ws失败: %v", err)
			// 如果广播消息失败，则将客户端从 clients 映射中删除
			delete(clients, conn)
		}
	}
}

// 每 5 秒向客户端发送一个 ping 消息，以保持连接
func sendPingMessage(conn *websocket.Conn) {
	for {
		time.Sleep(5 * time.Second)
		err := conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			log.Printf("ping本地ws客户端失败: %v", err)
			// 如果发送 ping 消息失败，则将客户端从 clients 映射中删除
			delete(clients, conn)
			break
		}
	}
}
