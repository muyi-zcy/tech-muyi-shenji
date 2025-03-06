package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "tech-muyi-shenji/routers"
)

func main() {
	logs.SetLogger("console")
	beego.Run()
}
