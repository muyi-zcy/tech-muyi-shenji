package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"tech-muyi-shenji/controllers"
)

func init() {
	beego.Router("/api/shenji/app/create", &controllers.AppController{}, "post:CreateApp")
	beego.Router("/api/shenji/app/update", &controllers.AppController{}, "put:UpdateApp")
	beego.Router("/api/shenji/app/:appCode", &controllers.AppController{}, "delete:DeleteApp")
	beego.Router("/api/shenji/app/publish", &controllers.AppController{}, "post:Publish")
	beego.Router("/api/shenji/app/release", &controllers.AppController{}, "post:Release")
	beego.Router("/api/shenji/app/rollback", &controllers.AppController{}, "post:Rollback")

	beego.Router("/api/shenji/app/:appCode", &controllers.AppController{}, "get:GetAppInfo")

	beego.Router("/api/shenji/menu/create", &controllers.MenuController{}, "post:CreateMenu")
	beego.Router("/api/shenji/menu/update", &controllers.MenuController{}, "put:UpdateMenu")
	beego.Router("/api/shenji/menu/delete", &controllers.MenuController{}, "post:DeleteMenu")
	beego.Router("/api/shenji/menu/sort", &controllers.MenuController{}, "post:SortMenu")
	beego.Router("/api/shenji/menu/remove", &controllers.MenuController{}, "post:RemoveMenu")
	beego.Router("/api/shenji/menu/:appCode", &controllers.MenuController{}, "get:GetMenu")

	beego.Router("/api/shenji/pages/:appCode/:menuCode", &controllers.PageController{}, "post:SavePage")
	beego.Router("/api/shenji/pages/:appCode/:menuCode", &controllers.PageController{}, "get:GetPage")

	beego.Router("/api/shenji/:appCode/app.json", &controllers.AppController{}, "get:GetReleaseAppInfo")
	beego.Router("/api/shenji/:appCode/menu.json", &controllers.MenuController{}, "get:GetReleaseMenu")
	beego.Router("/api/shenji/:appCode/pages/:menuCode\\.json", &controllers.PageController{}, "get:GetReleasePage")

}
