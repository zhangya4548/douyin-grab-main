package constv

import "time"

const (
	NmidServerHost       = "127.0.0.1"
	NmidServerPort       = "6808"
	LogMode              = "debug"
	RuntimeRootPath      = ""
	LogSavePath          = "logs/"
	LogSaveName          = "log"
	LogFileExt           = "log"
	TimeFormat           = "20060102"
	DefaultHeartbeatTime = time.Second * 10
	// 本地服务配置
	WebPort      = "40000"
	WsPort       = "40001"
	WsClientPort = "localhost:40001"

	// 远程服务配置
	DistalWsPort = "lwww.wykji.cn:53331"
	DistalWsPath = "/wss/dan/mu/conn"

	// 抖音配置
	DouYinOrigin = "https://live.douyin.com"
	LiveUrl      = "https://live.douyin.com/333096532012"
	LiveWsUrl    = "wss://webcast3-ws-web-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:7217444389554244409|wss_push_did:7200771809301956155|dim_log_id:20230402221712487343F346D5DB7E1BBC|fetch_time:1680445032718|seq:1|wss_info:0-1680445032718-0-0|wrds_kvs:WebcastRoomStatsMessage-1680445028019016337_WebcastRoomRankMessage-1680444944157360568&cursor=d-1_u-1_h-1_t-1680445032718_r-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7200771809301956155&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/108.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7217444389554244409&heartbeatDuration=0&signature=W0wY1TxcSnZ55iJ/"
)
