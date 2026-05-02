package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"tech-muyi-shenji/common"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) ok(data interface{}) {
	res := common.MyResult{
		"0", true, "", data, nil,
	}
	c.Data["json"] = &res
	c.ServeJSON()
}

func (c *BaseController) queryOk(data interface{}, query common.MyQuery) {
	res := common.MyResult{
		"0", true, "", data, query,
	}
	c.Data["json"] = &res
	c.ServeJSON()

}

func (c *BaseController) error(code string, message string, data interface{}) {

	res := common.MyResult{
		code, false, message, data, nil,
	}
	c.Data["json"] = &res
	c.ServeJSON()
}

func (c *BaseController) invalidAppCode() {
	c.error("00111050", "无效的应用编码", nil)
}

func (c *BaseController) invalidMenuCode() {
	c.error("00111050", "无效的菜单编码", nil)
}
