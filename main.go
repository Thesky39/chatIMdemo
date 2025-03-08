package main

import (
	"chat/cache"
	"chat/conf"
	"chat/router"
	"chat/service"
)

func main() {
	cache.Init()
	conf.Init()
	go service.Manager.Start()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
