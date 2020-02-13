package app

import (
	"fmt"
	"github.com/2020wfw/app/config"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Proxy struct {
	//前置反向代理
	PreHandler func(w http.ResponseWriter, r *http.Request) bool
	//调用反向代理
	MiddelHandler func(w http.ResponseWriter, r *http.Request)
	//调用后反向代理
	BackHandler func(w http.ResponseWriter, r *http.Request)
	Router Router
}

//创建代理对象
func NewProxy(router Router) Proxy {
	return Proxy{
		PreHandler:func(w http.ResponseWriter, r *http.Request) (ret bool) {
			fmt.Println("前置过滤器",r.URL)
			if router.IsUrlCanceled(r){
				ret = true
			}
			return
		},
		MiddelHandler:func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("这里需要转发",r.URL.Path)
			transPort(w,r)
		},
		BackHandler:func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("后置过滤器",r.URL)
		},
		Router:router,
	}
}

//启动反向代理
func (this *Proxy) StartProxy(){

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {

		if !this.Router.IsUrlExist(r) {
			fmt.Println("路由不存在,拦截器已拦截请求",r.URL)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("路由不存在,拦截器已拦截请求"))
			return
		}

		if this.PreHandler != nil{
			if isBreak := this.PreHandler(w,r);isBreak{
				fmt.Println("路由已被注销,拦截器已拦截请求",r.URL)
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("路由已被注销,拦截器已拦截请求"))
				return
			}
		}

		if this.PreHandler != nil{
			this.MiddelHandler(w,r)
		}

		if this.PreHandler != nil{
			this.BackHandler(w,r)
		}
	})

	http.HandleFunc("/route/remove",func(w http.ResponseWriter, r *http.Request) {
		UrlCancel(w,r,this.Router)
	})

	http.HandleFunc("/route/add",func(w http.ResponseWriter, r *http.Request) {
		UrlAdd(w,r,this.Router)
	})

	http.ListenAndServe("127.0.0.1:1234",nil)


	//http.HandleFunc("/auth",func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("哈哈哈哈哈哈哈哈哈哈哈"))
	//		fmt.Println("第二台收到",r.URL)
	//})
	//http.ListenAndServe("127.0.0.1:5678",nil)

}

//注销路由
func UrlCancel (w http.ResponseWriter, r *http.Request,router Router) {
	var (
		msg string
		code int
		url string
	)

	r.ParseForm()
	if len(r.Form) > 0 {
		url = r.Form.Get("url")
	}

	if router.RoutersCheck != nil {

		if _,ok := router.RoutersCheck[url];ok{
			//标记已移除
			router.RoutersCheck[url] = false
			msg = "OK"
			code = http.StatusOK
		} else{
			//返回错误
			msg = "路由不存在"
			code = http.StatusNotFound
		}
	}

	w.WriteHeader(code)
	w.Write([]byte(msg))
}

//恢复被标记路由
func UrlAdd (w http.ResponseWriter, r *http.Request,router Router) {
	var (
		msg string
		code int
		url string
	)

	r.ParseForm()
	if len(r.Form) > 0 {
		url = r.Form.Get("url")
	}

	if router.RoutersCheck != nil {
		if _,ok := router.RoutersCheck[url];ok{
			//恢复标记为正常
			router.RoutersCheck[url] = true
			msg = "OK"
			code = http.StatusOK
		} else{
			//返回错误
			msg = "路由不存在"
			code = http.StatusNotFound
		}
	}

	w.WriteHeader(code)
	w.Write([]byte(msg))
}

//原生代理请求
func transPort (w http.ResponseWriter, r *http.Request) {

	ins := config.Consu.GetInstance(config.RoutersMap[r.URL.Path])
	ipPort := ""
	if len(ins.Addrs) > 0 {
		ipPort = ins.Addrs[0]
	}

	proxy := func(_ *http.Request) (url *url.URL, err error) {
		return url.Parse(ipPort + r.URL.Path)
	}
	transport := &http.Transport{
		Proxy: proxy,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
	}

	client := &http.Client{Transport: transport}
	url := ipPort + r.URL.Path
	req, err := http.NewRequest(r.Method, url, r.Body)
	//注： 设置Request头部信息
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//注： 设置Response头部信息
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	data, _ := ioutil.ReadAll(resp.Body)
	w.Write(data)
}