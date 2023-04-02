package web

import (
	"douyin-grab/constv"
	queue2 "douyin-grab/pkg/queue"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func RunWsClient(qu *queue2.QueueSrv) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: constv.WsClientPort, Path: "/ws"}
	log.Printf("连接到本地ws服务端%s", u.String())
	// 连接 websocket 服务器
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("连接本地ws服务端失败:", err)
	}
	defer c.Close()
	// 启动两个协程，一个用于接收服务器发送的消息，一个用于向服务器发送消息
	done := make(chan struct{})
	// go receiveMessages(c, done)
	go sendMessages(c, qu)
	// 等待中断信号，当收到中断信号时，关闭连接并退出程序
	for {
		select {
		case <-interrupt:
			log.Println("收到本地ws服务端中断信号，正在关闭连接...")
			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Printf("发送关闭消息到本地ws时出错\n\n: %v", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// 接收服务器发送的消息
func receiveMessages(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("读取本地ws消息时出错\n\n: %v", err)
			return
		}
		log.Printf("收到本地ws服务端的消息: %s", string(message))
	}
}

// 向服务器发送消息
func sendMessages(c *websocket.Conn, qu *queue2.QueueSrv) {
	for {
		// 取对列数据
		message := qu.Pop()
		if message == "" {
			time.Sleep(time.Second * 2)
			continue
		}

		err := c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("发送数据到本地ws服务端异常: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}

		log.Printf("发送数据到本地ws服务端完: %s", message)
	}
}
