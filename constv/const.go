package constv

import "time"

const (
	WsPort  = "40001"
	WebPort = "40000"

	DOUYIORIGIN          = "https://live.douyin.com"
	DOUYINHOST           = "webcast3-ws-web-hl.douyin.com"
	DOUYINPATH           = "/webcast/im/push/v2"
	DEFAULTLIVEROOMURL   = "https://live.douyin.com/543806645662"
	DEFAULTLIVEWSSURL    = "wss://webcast3-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:7218565373859810107|wss_push_did:7200771809301956155|dim_log_id:20230405224209DD81988254A7EB00B86D|fetch_time:1680705729675|seq:1|wss_info:0-1680705729675-0-0|wrds_kvs:InputPanelComponentSyncData-1680703245571292273_WebcastRoomRankMessage-1680705693488297889_WebcastRoomStatsMessage-1680705729445790270_HighlightContainerSyncData-4&cursor=h-1_t-1680705729675_r-1_d-1_u-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7200771809301956155&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/108.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7218565373859810107&heartbeatDuration=0&signature=W0CSHnEx8p46PrKQ"
	DEFAULTHEARTBEATTIME = time.Second * 10
	WSURL                = `www.wykji.cn:53331/wss/dan/mu/conn`
)
