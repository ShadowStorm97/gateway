package app

import (
	"net/http"
	"github.com/2020wfw/app/config"
)


type Router struct {
	//初始化路由表开关
	RoutersCheck map[string]bool
}

func NewRouter() (ret Router) {
	ret = Router{}
	if len(config.RoutersMap) > 0 {
		ret.RoutersCheck = make(map[string]bool)
		for k,_ := range config.RoutersMap{
			ret.RoutersCheck[k] = true
		}
	}else{
		//网关路由初始化失败 路由表缺失
		panic("网关路由初始化失败 路由表缺失")
	}
	return
}

func (this *Router) IsUrlCanceled (r *http.Request) (ret bool) {
	if val,ok := this.RoutersCheck[r.URL.String()];ok{
		if !val{
			ret = true
		}
	}else{
		ret = true
	}
	return
}

func (this *Router) IsUrlExist (r *http.Request) (ret bool){
	if _,ok := this.RoutersCheck[r.URL.String()];ok{
		ret = true
	}
	return
}





