package main

import (
	_ "tunaiku/routers"
	"tunaiku/utilities/mongo"
	"github.com/astaxie/beego"
	"os"
)

func main() {
	err := mongo.Startup("main")
	if err != nil {
		os.Exit(1)
	}

	beego.Run()
}