package controllers

import (
	"encoding/json"
	"io"

	"github.com/beego/beego/v2/core/logs"
	"tech-muyi-shenji/app"
)

type PageController struct {
	BaseController
}

// SavePage 保存开发态低代码页面 JSON。
func (c *PageController) SavePage() {
	appCode := c.Ctx.Input.Param(":appCode")
	menuCode := c.Ctx.Input.Param(":menuCode")
	if !app.ValidResourceCode(appCode) || !app.ValidResourceCode(menuCode) {
		if !app.ValidResourceCode(appCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}

	var request map[string]interface{}
	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111140", "读取请求体失败", nil)
		return
	}
	if err := json.Unmarshal(body, &request); err != nil {
		logs.Error(err.Error())
		c.error("00111141", "解析JSON失败", nil)
		return
	}
	if err := app.GetMenuManagerInstance().SaveMenuPage(appCode, menuCode, request); err != nil {
		logs.Error(err.Error())
		c.error("00111142", "页面保存失败", nil)
		return
	}
	c.ok(true)
}

func (c *PageController) GetPage() {
	appCode, version := appCodeVersionForDevAPI(c.Ctx.Input.Param(":appCode"))
	menuCode := c.Ctx.Input.Param(":menuCode")
	if !app.ValidResourceCode(appCode) || !app.ValidResourceCode(menuCode) {
		if !app.ValidResourceCode(appCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}
	config, err := app.GetMenuManagerInstance().GetMenuPage(appCode, version, menuCode)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111143", "页面获取失败", nil)
		return
	}
	c.ok(config)
}

func (c *PageController) GetReleasePage() {
	appCode, version := appCodeVersionForReleaseAPI(c.Ctx.Input.Param(":appCode"))
	menuCode := c.Ctx.Input.Param(":menuCode")
	if !app.ValidResourceCode(appCode) || !app.ValidResourceCode(menuCode) {
		if !app.ValidResourceCode(appCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}
	config, err := app.GetMenuManagerInstance().GetMenuPage(appCode, version, menuCode)
	if err != nil {
		logs.Error(err.Error())
		c.ok(map[string]interface{}{})
		return
	}
	c.ok(config)
}
