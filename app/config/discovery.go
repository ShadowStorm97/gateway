package config

import (
	"fmt"
	"math/rand"
	"github.com/bilibili/discovery/naming"
	"time"
	"context"
)

var Consu *consumer

// This Example register a server provider into discovery.
func DiscoveryRegister() (cancel context.CancelFunc) {
	conf := &naming.Config{
		Nodes: config.Nodes, // NOTE: 配置种子节点(1个或多个)，client内部可根据/discovery/nodes节点获取全部node(方便后面增减节点)
		Zone:  config.Zone,
		Env:   config.Env,
	}
	dis := naming.New(conf)
	ins := &naming.Instance{
		Zone:  config.Zone,
		Env:   config.Env,
		AppID: config.AppID,
		// Hostname:"", // NOTE: hostname 不需要，会优先使用discovery new时Config配置的值，如没有则从os.Hostname方法获取！！！
		Addrs:    config.Addrs,
		LastTs:   time.Now().Unix(),
		Metadata: config.Metadata,
	}
	cancel, _ = dis.Register(ins)
	fmt.Println("register")
	// Unordered output4
	return
}

type consumer struct {
	conf  *naming.Config
	appID string
	dis   naming.Resolver
	ins   []*naming.Instance
}

// This Example show how get watch a server provider and get provider instances.
func ResolverWatch() {
	conf := &naming.Config{
		Nodes: config.Nodes,
		Zone:  config.Zone,
		Env:   config.Env,
	}
	dis := naming.New(conf)
	c := &consumer{
		conf:  conf,
		appID: config.AppID,
		dis:   dis.Build(config.AppID),
	}
	Consu = c
	rsl := dis.Build(c.appID)
	ch := rsl.Watch()
	go c.getInstances(ch)
}

func (c *consumer) getInstances(ch <-chan struct{}) {
	for { // NOTE: 通过watch返回的event chan =>
		if _, ok := <-ch; !ok {
			return
		}
		// NOTE: <= 实时fetch最新的instance实例
		ins, ok := c.dis.Fetch()
		if !ok {
			continue
		}
		// get local zone instances, otherwise get all zone instances.
		if in, ok := ins.Instances[c.conf.Zone]; ok {
			c.ins = in
		} else {
			for _, in := range ins.Instances {
				c.ins = append(c.ins, in...)
			}
		}
	}
}

func (c *consumer) GetInstance(appID string) (ins *naming.Instance) {
	// get instance by loadbalance
	// you can use any loadbalance algorithm what you want.
	len := len(c.ins)
	fmt.Println("ins len :",len,"c",c.appID)
	index := rand.Intn(len)
	ins = c.ins[index]
	return
}