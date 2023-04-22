package wsocket

import (
	"bytes"
	"compress/gzip"
	"douyin-grab/pkg/cache"
	queue2 "douyin-grab/pkg/queue"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"douyin-grab/constv"
	"douyin-grab/grab"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
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
	Ttwid       string
	LiveRoomUrl string
	WSServerUrl string
	Header      http.Header
	ClientCon   *websocket.Conn
	qu          *queue2.QueueSrv
	cache       *cache.Cache
}

func (client *WSClient) SetWSServerUrl(WSServerUrl string) {
	client.WSServerUrl = WSServerUrl
}

func (client *WSClient) SetLiveRoomUrl(LiveRoomUrl string) {
	client.LiveRoomUrl = LiveRoomUrl
}

func NewWSClient(qu *queue2.QueueSrv, cache *cache.Cache) *WSClient {
	return &WSClient{
		qu:    qu,
		cache: cache,
	}
}

func (client *WSClient) Run() {
	client.SetRequestInfo()
	client.ConnWSServer()
	client.RunWSClient()
}

func (client *WSClient) SetRequestInfo() *WSClient {
	_, ttwid := grab.FetchLiveRoomInfo(client.LiveRoomUrl)

	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	header.Set("Origin", constv.DOUYIORIGIN)
	cookie := &http.Cookie{
		Name:  "ttwid",
		Value: ttwid,
	}
	header.Add("Cookie", cookie.String())

	client.Header = header
	client.Ttwid = ttwid

	return client
}

func (client *WSClient) ConnWSServer() *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(client.WSServerUrl, client.Header)
	if err != nil {
		log.Println("websocket dial:", err)
	}

	client.ClientCon = c

	return c
}

func (client *WSClient) RunWSClient() {
	if client.ClientCon != nil {
		go func() {
			for {
				_, message, err := client.ClientCon.ReadMessage()
				if err != nil {
					log.Println("read error", err.Error())
					return
				}

				wssPackage := &grab.PushFrame{}
				err = proto.Unmarshal(message, wssPackage)
				if err != nil {
					log.Println("unmarshaling proto wssPackage error: ", err)
				}
				logId := wssPackage.LogId
				log.Println("[douyin] logid", logId)

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

				payloadPackage := &grab.Response{}
				err = proto.Unmarshal(decompressed, payloadPackage)
				if err != nil {
					log.Println("unmarshaling proto payloadPackage error: ", err)
				}

				if payloadPackage.NeedAck {
					client.sendAck(logId, payloadPackage.InternalExt)
				}

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

		go func() {
			for {
				duration := constv.DEFAULTHEARTBEATTIME
				timer := time.NewTimer(duration)
				<-timer.C
				client.heartBeat()
			}
		}()
	}
}

func unPackWebcastChatMessage(payload []byte) string {
	msg := &grab.ChatMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		log.Println("unmarshaling proto unPackWebcastChatMessage error: ", err)
		return ""
	}

	log.Println("[unPackWebcastChatMessage] [📧直播间弹幕消息]", msg.Content)
	return msg.Content
}

func unPackWebcastLikeMessage(payload []byte) string {
	msg := &grab.LikeMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		log.Println("unmarshaling proto unPackWebcastLikeMessage error: ", err)
		return ""
	}
	log.Println("[unPackWebcastLikeMessage] [👍直播间点赞消息]", msg.User.NickName+"点赞")
	return msg.User.NickName + "点赞"
}

func unPackWebcastGiftMessage(payload []byte) string {
	msg := &grab.GiftMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		log.Println("unmarshaling proto unPackWebcastGiftMessage error: ", err)
		return ""
	}
	log.Println("[unPackWebcastGiftMessage] [🎁直播间礼物消息]", msg.Common.Describe)
	return msg.Common.Describe
}

func unPackWebcastMemberMessage(payload []byte) string {
	msg := &grab.MemberMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		log.Println("unmarshaling proto unPackWebcastMemberMessage error: ", err)
		return ""
	}
	log.Println("[unPackWebcastMemberMessage] [🚹🚺直播间成员进入消息]", msg.User.NickName+"进入直播间")
	return msg.User.NickName + "进入直播间"
}

func (client *WSClient) sendAck(logId uint64, InternalExt string) {
	obj := &grab.PushFrame{}
	obj.PayloadType = "ack"
	obj.LogId = logId
	obj.PayloadType = InternalExt
	data, err := proto.Marshal(obj)
	if err != nil {
		log.Println("send ack error", err)
	}

	client.SendBytes(data)
	log.Println("[sendAck] [🌟发送Ack]")
}

func (client *WSClient) heartBeat() {
	obj := &grab.PushFrame{}
	obj.PayloadType = "hb"
	data, err := proto.Marshal(obj)
	if err != nil {
		log.Println("send ack error", err)
	}

	client.SendBytes(data)
	log.Println("[ping] [💗发送ping心跳]")
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
