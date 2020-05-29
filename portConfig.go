package main

import (
	"container/list"
	"math/rand"
	"sync"
	"time"
)

type portConfig struct {
	ports *list.List
	size int
	start int
}
var pc *portConfig
var mu sync.Mutex

func initPortConfig(){
	if pc == nil {
		mu.Lock()
		defer mu.Unlock()
		if pc == nil {
			pc = &portConfig{list.New(), 1000, 8000}
			pc.initPorts()
		}
	}
}

func initPortConfigWithInit(size int, start int) {
	if pc == nil {
		mu.Lock()
		defer mu.Unlock()
		if pc == nil {
			pc = &portConfig{list.New(), size, start}
			pc.initPorts()
		}
	}
}

//初始化端口集
func (pc *portConfig) initPorts() {
	if pc.size < 1000 {
		pc.size = 1000
	}

	for i:=0; i<pc.size; i++ {
		pc.ports.PushBack(pc.start+i)
	}
}

//端口资源回放
func (pc *portConfig) setPort(port int) {
	if port >= pc.start && port <= (pc.start + pc.size) {
		pc.ports.PushBack(port)
	}
}

//获取服务端用于外网请求端口
func (pc *portConfig) getPort() int {
	if pc.ports.Len() > 0 {
		port := pc.ports.Front()
		pc.ports.Remove(port)
		return port.Value.(int)
	} else {
		time.Sleep(1 * time.Nanosecond)
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(65535)
	}
}