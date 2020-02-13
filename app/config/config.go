package config

import "time"

//public
var (
	RoutersMap map[string]string
)

//private
var (
	config Configer

	lastUpdateTime time.Time
)

func LoadConfig(){
	//读取配置 (暂时写死 后期从配置中心读取 热更新)
	config = Configer{
			Nodes:[]string{"127.0.0.1:7171"},
			Zone:"sh1",
			Env:"test",
			AppID:"provider_0",
			Addrs:[]string{"http://127.0.0.1:5678", "grpc://172.0.0.1:9999"},
			Metadata:map[string]string{"weight": "10"},
			RoutersMap:map[string]string{"/auth":"provider_0","/users":"provider_1"},
		}
	RoutersMap = config.RoutersMap
	lastUpdateTime = time.Now()

	//启动热更新
	startHotUpdate()
}

//配置项热更新
func startHotUpdate(){

	lastUpdateTime = time.Now()
}

//获取配置
func GetConfigByKey() (ret Configer) {
	return config
}

//网关服务配置项
type Configer struct {

	//discover配置
	Nodes []string
	Zone string
	Env string
	AppID string
	Addrs []string
	Metadata map[string]string

	//路由配置  route url => app_id
	RoutersMap map[string]string
}