package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"strings"
	"tech-muyi-shenji/app"
)

type PageController struct {
	BaseController
}

// 保存页面配置
func (c *PageController) SavePage() {
	appCode := c.Ctx.Input.Param(":appCode")
	menuCode := c.Ctx.Input.Param(":menuCode")

	// 读取整个body为json
	var request map[string]interface{}

	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111122", "读取请求体失败", nil)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111122", "解析JSON失败", nil)
		return
	}

	err = app.GetMenuManagerInstance().SaveMenuPage(appCode, menuCode, request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111122", "页面保存失败", nil)
	}
	c.ok(true)
}

func (c *PageController) GetPage() {
	getAppCode := c.Ctx.Input.Param(":appCode")
	menuCode := c.Ctx.Input.Param(":menuCode")

	version := ""
	appCode := ""
	if strings.Contains(getAppCode, ".") {
		parts := strings.Split(getAppCode, ".")
		appCode = parts[0]
		version = strings.Join(parts[1:], ".")
	} else {
		appCode = getAppCode
		version = "dev"
	}

	config, err := app.GetMenuManagerInstance().GetMenuPage(appCode, version, menuCode)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111114", "页面获取失败", nil)
		return
	}
	c.ok(config)
}

func (c *PageController) GetReleasePage() {
	getAppCode := c.Ctx.Input.Param(":appCode")
	menuCode := c.Ctx.Input.Param(":menuCode")

	version := ""
	appCode := ""
	if strings.Contains(getAppCode, ".") {
		parts := strings.Split(getAppCode, ".")
		appCode = parts[0]
		version = strings.Join(parts[1:], ".")
	} else {
		appCode = getAppCode
	}

	config, err := app.GetMenuManagerInstance().GetMenuPage(appCode, version, menuCode)
	if err != nil {
		logs.Error(err.Error())
		// 返回空json
		c.ok(map[string]interface{}{})
	}
	c.ok(config)
}
