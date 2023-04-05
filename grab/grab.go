package grab

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type RoomInfo struct {
	App struct {
		InitialState struct {
			RoomStore struct {
				RoomInfo struct {
					RoomId string `json:"roomId"`
					Room   struct {
						Title        string `json:"title"`
						UserCountStr string `json:"user_count_str"`
					} `json:"room"`
				} `json:"roomInfo"`
			} `json:"roomStore"`
		} `json:"initialState"`
	} `json:"app"`
}

func FetchLiveRoomInfo(roomUrl string) (*RoomInfo, string) {
	req, err := http.NewRequest("GET", roomUrl, nil)
	if err != nil {
		log.Println("fetch live room info err", err)
		return nil, ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") // 设置User-Agent头
	cookie := &http.Cookie{Name: "__ac_nonce", Value: "063abcffa00ed8507d599"}
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("fetch live room info err", err)
		return nil, ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("read res body err", err)
		return nil, ""
	}

	pattern := regexp.MustCompile(`<script id="RENDER_DATA" type="application/json">(.*?)</script>`)
	data := pattern.FindSubmatch(body)
	decodedUrl, err := url.QueryUnescape(string(data[1]))
	if err != nil {
		log.Println("url decode err", err)
		return nil, ""
	}

	var roomInfo RoomInfo
	err = json.Unmarshal([]byte(decodedUrl), &roomInfo)
	if err != nil {
		log.Println("json unmarshal err", err)
		return nil, ""
	}
	log.Println("roomid:", roomInfo.App.InitialState.RoomStore.RoomInfo.RoomId)
	log.Println("title:", roomInfo.App.InitialState.RoomStore.RoomInfo.Room.Title)
	log.Println("user_count:", roomInfo.App.InitialState.RoomStore.RoomInfo.Room.UserCountStr)

	var ttwid string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "ttwid" {
			ttwid = cookie.Value
		}
	}
	log.Println("ttwid:", ttwid)

	return &roomInfo, ttwid
}
