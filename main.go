package main

import (
	"go-web-template/configs"
	"go-web-template/router"
)

func main() {
	configs.InitConfig()
	router.InitRouter()
	//r := router.NewRouter()
	//_ = r.Run()
}
