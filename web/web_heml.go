package web

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
			    max-height: 500px;
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
		  text-align: left;
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
		h2{text-align: center;}
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
            <input type="text" id="input" name="input" size="60" placeholder="输入直播间url" >
            <label for="input2">直播间wsUrl</label>
            <input type="text" id="input2" name="input2" size="60" placeholder="输入直播间wsUrl">
			<br/>
            <button type="submit">开始</button>
        </form>
			<button id="stop">停止</button>
    </div>
	<div class="form-container-2">
        <h2>弹幕</h2>
        <table class="list-box">
		  <thead>
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
		

			// 实时获取弹幕数据===================================================
			const socket = new WebSocket("ws://localhost:40001/ws");
			// 添加事件监听器，当连接建立时触发
			socket.addEventListener('open', (event) => {
			  console.log('连接远程成功');
			});
			// 添加事件监听器，当收到服务器发送的消息时触发
			socket.addEventListener('message', (event) => {
			  console.log('收到远程数据:', event.data);
			  const tbody = document.querySelector('.list-box tbody');
			  const row = tbody.insertRow();
			  const titleCell = row.insertCell();
			  titleCell.textContent = event.data;

			  const bax = document.querySelector('.form-container-2');
			  bax.scrollTop = bax.scrollHeight;	
			});
			// 添加事件监听器，当连接关闭时触发
			socket.addEventListener('close', (event) => {
			  console.log('远程服务端关闭');
			});
			// 添加事件监听器，当连接发生错误时触发
			socket.addEventListener('error', (event) => {
			  console.error('连接远程异常:', event);
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
                    alert("成功开始");
				}
            };
            xhr.send(JSON.stringify({
                'input': formData.get('input'),
                'input2': formData.get('input2'),
            }));
        });

		// 停止 ================================================
		document.getElementById('stop').addEventListener('click', function() {
		  let xhr = new XMLHttpRequest();
		  xhr.open('POST', '/stop', true);
		  xhr.setRequestHeader('Content-Type', 'application/json');
		  xhr.onreadystatechange = function() {
		    if (xhr.readyState === XMLHttpRequest.DONE) {
				if (xhr.status !== 200) {
					alert(xhr.response);
                }else{
                    alert("成功停止");
				}
		    }
		  };
		  xhr.send(JSON.stringify({}));
		});

    </script>
</body>
</html>
`
}
