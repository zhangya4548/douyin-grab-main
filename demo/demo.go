package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func main() {
	resp, err := http.Get("https://live.douyin.com/168465302284")
	if err != nil {
		fmt.Println("http get error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		return
	}

	reg := regexp.MustCompile(`"ws":"(wss://webcast3-ws-web-lf\.douyin\.com/webcast/im/push/v2/[a-zA-Z0-9_\-]+)"`)
	match := reg.FindStringSubmatch(string(body))

	if len(match) > 1 {
		fmt.Println("websocket url:", match[1])
	} else {
		fmt.Println("websocket url not found")
	}
}
