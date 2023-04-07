package web

import (
	"douyin-grab/constv"
	"douyin-grab/pkg/cache"
	queue2 "douyin-grab/pkg/queue"
	"douyin-grab/wsocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Web struct {
	cache     *cache.Cache
	qu        *queue2.EsQueue
	douYinSrv *wsocket.WSClient
}

func NewWeb(qu *queue2.EsQueue, cache *cache.Cache, douYinSrv *wsocket.WSClient) *Web {
	return &Web{
		cache:     cache,
		qu:        qu,
		douYinSrv: douYinSrv,
	}
}

func (s *Web) RunWeb() {
	// 定义HTTP路由
	http.HandleFunc("/", s.formHandler)
	http.HandleFunc("/submit", s.submitHandler)
	http.HandleFunc("/getData", s.getDataHandler)
	http.HandleFunc("/stop", s.stopHandler)

	// 启动WebSocket服务器
	http.HandleFunc("/ws", s.wsHandler)
	go func() {
		log.Println("ws启动:", constv.WsPort)
		err := http.ListenAndServe(":"+constv.WsPort, nil)
		if err != nil {
			log.Printf("ws ListenAndServe异常: %s \n", err.Error())
			return
		}
	}()

	// 启动HTTP服务
	log.Println("web启动:", constv.WebPort)
	err := http.ListenAndServe(":"+constv.WebPort, nil)
	if err != nil {
		log.Println("ListenWeb异常:", err.Error())
	}
}

type FormData struct {
	LiveUrl   string `json:"input"`
	LiveWsUrl string `json:"input2"`
}

func (s *Web) formHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.getPageHtml())
}
func (s *Web) submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "请求方式错误", http.StatusMethodNotAllowed)
		return
	}

	var formData FormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	log.Println("submitHandler参数", formData)

	if formData.LiveUrl == "" {
		http.Error(w, "直播间url参数不能为空", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(formData.LiveUrl, constv.DOUYIORIGIN) {
		http.Error(w, "直播间url参数错误", http.StatusBadRequest)
		return
	}

	if formData.LiveWsUrl == "" {
		http.Error(w, "直播间wsUrl参数不能为空", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(formData.LiveWsUrl, "wss://") {
		http.Error(w, "直播间url参数错误", http.StatusBadRequest)
		return
	}

	// 有运行的,先停止
	stopState, err := s.cache.Get("Stop")
	if err != nil {
		http.Error(w, "获取Stop设置错误", http.StatusBadRequest)
		return
	}
	if stopState == "false" {
		s.douYinSrv.Close()
	}

	// 保存配置
	err = s.cache.Set("LiveRoomUrl", formData.LiveUrl)
	if err != nil {
		http.Error(w, "设置LiveRoomUrl错误", http.StatusBadRequest)
		return
	}
	err = s.cache.Set("WssUrl", formData.LiveWsUrl)
	if err != nil {
		http.Error(w, "设置LiveWsUrl错误", http.StatusBadRequest)
		return
	}

	err = s.cache.Set("Stop", "false")
	if err != nil {
		http.Error(w, "设置Stop错误", http.StatusBadRequest)
		return
	}

	// 运行
	s.douYinSrv.SetLiveRoomUrl(formData.LiveUrl)
	s.douYinSrv.SetWSServerUrl(formData.LiveWsUrl)
	go s.douYinSrv.Run()

	jsonData, err := json.Marshal(formData)
	if err != nil {
		http.Error(w, "服务异常", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func (s *Web) stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "请求方式错误", http.StatusMethodNotAllowed)
		return
	}
	log.Println("提交停止了")

	err := s.cache.Set("Stop", "true")
	if err != nil {
		http.Error(w, "停止异常", http.StatusMethodNotAllowed)
		return
	}

	go s.douYinSrv.Close()

	jsonData, err := json.Marshal(FormData{})
	if err != nil {
		http.Error(w, "服务异常", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func (s *Web) getDataHandler(w http.ResponseWriter, r *http.Request) {
	// 获取配置
	LiveRoomUrl, err := s.cache.Get("LiveRoomUrl")
	if err != nil {
		http.Error(w, "获取LiveRoomUrl配置异常", http.StatusInternalServerError)
		return
	}

	WssUrl, err := s.cache.Get("WssUrl")
	if err != nil {
		http.Error(w, "获取WssUrl配置异常", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(FormData{
		LiveUrl:   LiveRoomUrl,
		LiveWsUrl: WssUrl,
	})
	if err != nil {
		http.Error(w, "服务异常", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
