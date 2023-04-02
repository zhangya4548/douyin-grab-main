package web

import (
	"douyin-grab/constv"
	"douyin-grab/pkg/cache"
	"douyin-grab/wsocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Web struct {
	cache     *cache.Cache
	douYinSrv *wsocket.DouYinSrv
}

func NewWeb(cache *cache.Cache, douYinSrv *wsocket.DouYinSrv) *Web {
	return &Web{cache: cache, douYinSrv: douYinSrv}
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
		log.Printf("ListenWeb异常: %s \n", err.Error())
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
	log.Println("submitHandler参数", r.Body)

	var formData FormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	if formData.LiveUrl == "" {
		http.Error(w, "直播间url参数错误", http.StatusBadRequest)
		return
	}

	if formData.LiveWsUrl == "" {
		http.Error(w, "直播间wsUrl参数错误", http.StatusBadRequest)
		return
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
	log.Println("submitHandler参数", r.Body)

	// 停止
	err := s.cache.Set("Stop", "true")
	if err != nil {
		http.Error(w, "设置Stop错误", http.StatusBadRequest)
		return
	}

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

func (s *Web) getPageHtml() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Form Page</title>
    <style>
        body {
            background-color: #f2f2f2;
            font-family: Arial, sans-serif;
        }
		.container {
            display: flex;
            flex-direction: row;
            justify-content: center; /* 水平居中对齐 */
            height: auto; /* 设置容器高度为100% */
        }
        .form-container {
			float: left;
            max-width: 400px;
           /* margin: 30px auto;*/
            padding: 30px;
            background-color: white;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
		
        .form-container-2 {
			    float: left;
			    margin-left: 5px;
			    max-height: 600px;
			    max-width: 900px;
			    padding: 5px;
			    background-color: white;
			    border-radius: 5px;
			    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			    overflow: auto;
        }
        .form-container h2 {
            margin-top: 0;
        }
        .form-container label {
            display: block;
            margin-top: 20px;
            font-weight: bold;
        }
        .form-container input[type=text], .form-container textarea {
            width: 100%;
            padding: 10px;
            margin-top: 5px;
            margin-bottom: 20px;
            border: none;
            border-radius: 5px;
            box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
            transition: box-shadow 0.3s;
        }
        .form-container input[type=text]:focus, .form-container textarea:focus {
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
        }
        .form-container button[type=submit] {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
        .form-container button[type=submit]:hover {
            background-color: #3e8e41;
        }
        .form-container p {
            margin-top: 20px;
            font-size: 14px;
            color: #666;
        }
        .form-container code {
            font-family: monospace;
        }
		/* CSS代码 */
		.list-box {
		  width: 800px;
		  border-collapse: collapse;
		}
		.list-box th,
		.list-box td {
		  border: 1px solid #ddd;
		  padding: 8px;
		  text-align: center;
		}
		.list-box th {
		  background-color: #f2f2f2;
		}
		.list-box tbody tr:nth-child(even) {
		  background-color: #f2f2f2;
		}
		td{
			text-overflow:ellipsis; overflow:hidden; white-space:nowrap;
		}
		h2,th{text-align: center;}
		#stop {
		    background-color: #af4c4c;
		    color: white;
		    padding: 10px 20px;
		    border: none;
		    border-radius: 5px;
		    cursor: pointer;
		}
    </style>
</head>
<body>
<div class="container">
    <div class="form-container">
        <h2>设置</h2>
        <form method="post" action="/submit">
            <label for="input">直播间url</label>
            <input type="text" id="input" name="input" size="60" placeholder="输入直播间url">
            <label for="input2">直播间wsUrl</label>
            <input type="text" id="input2" name="input2" size="60" placeholder="输入直播间wsUrl">
			<br/>
            <button type="submit">开始</button>
        </form>
			<button id="stop">停止</button>
    </div>
	<div class="form-container-2">
        <h2>完成记录</h2>
        <table class="list-box">
		  <thead>
		    <tr>
		      <th class="table-class" style="text-align:left;">内容</th>
		    </tr>
		  </thead>
		  <tbody>
		    
		  </tbody>
		</table>
    </div>
</div>
    <script>
		window.onload = function() {
			// 获取配置===============================================================
            // 创建XMLHttpRequest对象
            const xhr = new XMLHttpRequest();
            // 设置请求方法和请求地址
            xhr.open('GET', '/getData');
            // 设置响应类型为json
            xhr.responseType = 'json';
            // 监听请求完成事件
            xhr.onload = function() {
                // 如果请求成功，则修改表单对应的value
                if (xhr.status !== 200) {
					alert(xhr.response);
                }else{
					document.getElementById('input').value = xhr.response.input;
                    document.getElementById('input2').value = xhr.response.input2;
				}
            };
            // 发送请求
            xhr.send();
		

			// 实时获取上传完成数据===================================================
			const socket = new WebSocket("ws://47.242.27.85:40001/ws");
			// 添加事件监听器，当连接建立时触发
			socket.addEventListener('open', (event) => {
			  console.log('WebSocket connected');
			});
			// 添加事件监听器，当收到服务器发送的消息时触发
			socket.addEventListener('message', (event) => {
			 // console.log('Received message from server:', event.data);
			  const data = JSON.parse(event.data);
			  const tbody = document.querySelector('.list-box tbody');
			  const row = tbody.insertRow();
			  const titleCell = row.insertCell();
			  const urlCell = row.insertCell();
			  const statusCell = row.insertCell();
			  const timeCell = row.insertCell();
			  titleCell.textContent = data.title;
			  const link = document.createElement('a');
			  link.href = data.url;
			  link.target = '_blank';
			  link.textContent = '点击查看';
			  urlCell.appendChild(link);
              statusCell.textContent = data.is_upload;
			  //statusCell.textContent = data.is_upload ? '√' : '×';
			  timeCell.textContent = data.create_time;
			});
			// 添加事件监听器，当连接关闭时触发
			socket.addEventListener('close', (event) => {
			  console.log('WebSocket connection closed');
			});
			// 添加事件监听器，当连接发生错误时触发
			socket.addEventListener('error', (event) => {
			  console.error('WebSocket error:', event);
			});

        };


		// 表单设置提交=========================================================
        document.querySelector('form').addEventListener('submit', function(e) {
            e.preventDefault();
            var formData = new FormData(this);
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/submit');
            xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');
            xhr.onload = function() {
			    if (xhr.status !== 200) {
					alert(xhr.response);
                }else{
                    alert("设置成功");
				}

            };
            xhr.send(JSON.stringify({
                'input': formData.get('input'),
                'input2': formData.get('input2'),
                'textarea': formData.get('textarea')
            }));
        });


		function timestampToDatetime(timestamp, format) {
		  const date = new Date(timestamp);
		  const year = date.getFullYear();
		  const month = ('0' + (date.getMonth() + 1)).slice(-2);
		  const day = ('0' + date.getDate()).slice(-2);
		  const hour = ('0' + date.getHours()).slice(-2);
		  const minute = ('0' + date.getMinutes()).slice(-2);
		  const second = ('0' + date.getSeconds()).slice(-2);
		  const milliseconds = ('00' + date.getMilliseconds()).slice(-3);
		  let result = format || 'YYYY-MM-DD HH:mm:ss.SSS';
		  result = result.replace(/YYYY/g, year);
		  result = result.replace(/MM/g, month);
		  result = result.replace(/DD/g, day);
		  result = result.replace(/HH/g, hour);
		  result = result.replace(/mm/g, minute);
		  result = result.replace(/ss/g, second);
		  result = result.replace(/SSS/g, milliseconds);
		  return result;
		}
    </script>
</body>
</html>
`
}
