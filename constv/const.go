package constv

import "time"

const DOUYIORIGIN = "https://live.douyin.com"
const DOUYINHOST = "webcast3-ws-web-hl.douyin.com"
const DOUYINPATH = "/webcast/im/push/v2"
const DEFAULTLIVEROOMURL = "https://live.douyin.com/473225570786"
const DEFAULTLIVEWSSURL = "wss://webcast3-ws-web-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:7216284531085577019|wss_push_did:7216295709561161249|dim_log_id:20230330191324D7A59BD5C9789120992E|fetch_time:1680174804956|seq:1|wss_info:0-1680174804956-0-0|wrds_kvs:WebcastRoomRankMessage-1680174591094577592_LotteryInfoSyncData-1680174804645992193_WebcastRoomStatsMessage-1680174801043762059&cursor=t-1680174804956_r-1_d-1_u-1_h-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7216295709561161249&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1080&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/111.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7216284531085577019&heartbeatDuration=0&signature=RMLCNvz2xaXIE2t4"
const DEFAULTHEARTBEATTIME = time.Second * 10
const WSURL = `www.wykji.cn:53331/wss/dan/mu/conn`
