package wsocket

import (
	"bytes"
	"compress/gzip"
	queue2 "douyin-grab/pkg/queue"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"douyin-grab/constv"
	"douyin-grab/grab"
	"douyin-grab/pkg/logger"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type DYCookieJar struct {
	cookies []*http.Cookie
}

func (jar *DYCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *DYCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

type WSClient struct {
	WSServerUrl string
	Header      http.Header
	ClientCon   *websocket.Conn
	qu          *queue2.QueueSrv
}

func NewWSClient(qu *queue2.QueueSrv) *WSClient {
	return &WSClient{
		qu: qu,
	}
}

func (client *WSClient) SetRequestInfo(WSServerUrl string, header http.Header) *WSClient {
	client.WSServerUrl = WSServerUrl
	client.Header = header

	return client
}

func (client *WSClient) ConnWSServer(ttwid string) *websocket.Conn {
	// 创建一个 CookieJar，设置 Cookie
	// cookieJar := &DYCookieJar{cookies: []*http.Cookie{
	// 	&http.Cookie{Name: "ttwid", Value: TTWID},
	// }}
	// dialer := websocket.Dialer{
	// 	HandshakeTimeout: 5 * time.Second,
	// 	Jar:              cookieJar,
	// }
	// c, _, err := dialer.Dial(client.WSServerUrl, client.Header)
	c, _, err := websocket.DefaultDialer.Dial(client.WSServerUrl, client.Header)
	if err != nil {
		logger.Error("连接抖音Ws异常: %s", err)
	}

	client.ClientCon = c

	return c
}

func (client *WSClient) RunWSClient() {
	if client.ClientCon != nil {
		// read
		go func() {
			for {
				_, message, err := client.ClientCon.ReadMessage()
				if err != nil {
					logger.Error("读取抖音Ws异常: %s", err.Error())
					return
				}

				// --push frame--//
				wssPackage := &grab.PushFrame{}
				err = proto.Unmarshal(message, wssPackage)
				if err != nil {
					logger.Fatal("解析抖音wssPackage异常: %s", err)
				}
				logId := wssPackage.LogId
				logger.Info("解析到抖音logid: %s", logId)

				// --gizp decompress--//
				compressedDataReader := bytes.NewReader(wssPackage.Payload)
				gzipReader, err := gzip.NewReader(compressedDataReader)
				if err != nil {
					panic(err)
				}
				defer gzipReader.Close()

				decompressed, err := ioutil.ReadAll(gzipReader)
				if err != nil {
					panic(err)
				}
				// println(string(decompressed))

				// --response--//
				payloadPackage := &grab.Response{}
				err = proto.Unmarshal(decompressed, payloadPackage)
				if err != nil {
					logger.Fatal("解析抖音payloadPackage异常: %s", err)
				}

				// 返回ack
				if payloadPackage.NeedAck {
					client.sendAck(logId, payloadPackage.InternalExt)
				}

				// 打印各种消息
				for _, msg := range payloadPackage.MessagesList {
					str := ""
					if msg.Method == "WebcastChatMessage" {
						str = unPackWebcastChatMessage(msg.Payload)
					}
					if msg.Method == "WebcastLikeMessage" {
						str = unPackWebcastLikeMessage(msg.Payload)
					}
					if msg.Method == "WebcastGiftMessage" {
						str = unPackWebcastGiftMessage(msg.Payload)
					}
					if msg.Method == "WebcastMemberMessage" {
						str = unPackWebcastMemberMessage(msg.Payload)
					}
					client.qu.Push(str)
				}
			}
		}()

		// heartbeat
		go func() {
			for {
				duration := constv.DefaultHeartbeatTime
				timer := time.NewTimer(duration)
				<-timer.C
				client.heartBeat()
			}
		}()
	}
}

// 直播间弹幕消息
func unPackWebcastChatMessage(payload []byte) string {
	msg := &grab.ChatMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		logger.Fatal("解析抖音直播间弹幕消息异常: %s", err)
		return ""
	}

	logger.Info("[📧直播间弹幕消息] %s", msg.Content)
	return msg.Content
}

// 直播间点赞消息
func unPackWebcastLikeMessage(payload []byte) string {
	msg := &grab.LikeMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		logger.Fatal("解析抖音直播间点赞消息: %s", err)
		return ""
	}
	// likemsg, err := json.Marshal(msg)
	// if err != nil {
	// 	logger.Fatal("json marshal error: ", err)
	// }

	logger.Info("[👍直播间点赞消息] %s", msg.User.NickName+"点赞")
	return msg.User.NickName + "点赞"
}

// 直播间礼物消息
func unPackWebcastGiftMessage(payload []byte) string {
	msg := &grab.GiftMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		logger.Fatal("解析抖音直播间礼物消息: %s", err)
		return ""
	}
	// giftmsg, err := json.Marshal(msg)
	// if err != nil {
	// 	logger.Fatal("json marshal error: ", err)
	// }

	// logger.Info("[unPackWebcastGiftMessage] [🎁直播间礼物消息] json %s", giftmsg)
	logger.Info("[🎁直播间礼物消息] %s", msg.Common.Describe)
	return msg.Common.Describe
}

// 欢迎进入直播间
func unPackWebcastMemberMessage(payload []byte) string {
	msg := &grab.MemberMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		logger.Fatal("解析抖音直播间进入直播间消息异常: %s", err)
		return ""
	}
	// membermsg, err := json.Marshal(msg)
	// if err != nil {
	// 	logger.Fatal("json marshal error: ", err)
	// }

	// logger.Info("[unPackWebcastMemberMessage] [🚹🚺直播间成员进入消息] json %s", membermsg)
	logger.Info("[🚹🚺直播间成员进入消息] %s", msg.User.NickName+"进入直播间")
	return msg.User.NickName + "进入直播间"
}

// 发送ack
func (client *WSClient) sendAck(logId uint64, InternalExt string) {
	obj := &grab.PushFrame{}
	obj.PayloadType = "ack"
	obj.LogId = logId
	obj.PayloadType = InternalExt
	data, err := proto.Marshal(obj)
	if err != nil {
		logger.Error("发送Ack到抖音Ws异常: %s", err)
	}

	client.SendBytes(data)
	logger.Info("[🌟发送Ack到抖音Ws]")
}

// 发送心跳
func (client *WSClient) heartBeat() {
	obj := &grab.PushFrame{}
	obj.PayloadType = "hb"
	data, err := proto.Marshal(obj)
	if err != nil {
		logger.Error("发送ping心跳到抖音Ws异常: %s", err)
	}

	client.SendBytes(data)
	logger.Info("[💗发送ping心跳到抖音Ws]")
}

func (client *WSClient) SendBytes(buf []byte) error {
	return client.ClientCon.WriteMessage(websocket.BinaryMessage, buf)
}

func (client *WSClient) SendTexts(buf []byte) error {
	return client.ClientCon.WriteMessage(websocket.TextMessage, buf)
}

func (client *WSClient) Close() {
	if client.ClientCon != nil {
		client.ClientCon.Close()
	}
}
