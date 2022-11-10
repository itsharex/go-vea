package main

import (
	"go-vea/configs"
	"go-vea/router"
)

func main() {
	configs.InitConfig()
	router.InitRouter()
	//r := router.NewRouter()
	//_ = r.Run()
}
