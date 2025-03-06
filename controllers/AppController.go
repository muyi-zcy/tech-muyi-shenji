package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"tech-muyi-shenji/app"
)

// AppController 是一个控制器，用于处理与应用程序相关的请求。
type AppController struct {
	BaseController
}

// GetAppInfo 处理获取应用信息的请求。
func (c *AppController) GetAppInfo() {
	getAppCode := c.Ctx.Input.Param(":appCode")
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
	var err error
	var appInfo *app.App
	if version == "dev" {
		if !app.GetAppManagerInstance().IsAppDevExist(appCode) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppDevInfo(appCode)
	} else {
		if !app.GetAppManagerInstance().IsAppVersionExist(appCode, version) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppVersionInfo(appCode, version)
	}
	if err != nil {
		logs.Error(err.Error())
		c.error("00111052", "应用获取失败", nil)
		return
	}
	c.ok(appInfo)
}

func (c *AppController) GetReleaseAppInfo() {
	getAppCode := c.Ctx.Input.Param(":appCode")
	version := ""
	appCode := ""
	if strings.Contains(getAppCode, ".") {
		parts := strings.Split(getAppCode, ".")
		appCode = parts[0]
		version = strings.Join(parts[1:], ".")
	} else {
		appCode = getAppCode
	}
	var err error
	var appInfo *app.App
	if version == "" {
		if !app.GetAppManagerInstance().IsAppReleaseExist(appCode) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppReleaseInfo(appCode)
	}
	if version == "dev" {
		if !app.GetAppManagerInstance().IsAppDevExist(appCode) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppDevInfo(appCode)
	} else {
		if !app.GetAppManagerInstance().IsAppVersionExist(appCode, version) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppVersionInfo(appCode, version)
	}
	if err != nil {
		logs.Error(err.Error())
		c.error("00111052", "应用获取失败", nil)
		return
	}
	c.ok(appInfo)
}

// CreateApp 处理创建应用的请求。
// 该函数从请求体中解析应用信息，并检查应用是否已存在。如果应用不存在，则创建新应用。
// 如果解析失败或应用已存在，返回相应的错误信息。
func (c *AppController) CreateApp() {
	var appInfo app.App
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111061", "应用创建失败", nil)
		return
	}

	// 检查应用是否已存在
	if app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111062", "应用已存在", nil)
		return
	}

	// 创建新应用
	createApp, err := app.GetAppManagerInstance().CreateApp(appInfo.AppCode, appInfo.AppName, appInfo.Version, appInfo.Description)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111063", "应用创建失败:"+err.Error(), nil)
		return
	}
	c.ok(createApp)
}

// UpdateApp 处理更新应用的请求。
// 该函数从请求体中解析应用信息，并检查应用是否存在。如果应用存在，则更新应用信息。
// 如果解析失败或应用不存在，返回相应的错误信息。
func (c *AppController) UpdateApp() {
	var appInfo app.App
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111071", "应用更新失败", nil)
		return
	}

	// 检查应用是否存在
	if !app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111072", "应用不存在", nil)
		return
	}

	// 更新应用信息
	updateApp, err := app.GetAppManagerInstance().UpdateApp(appInfo.AppCode, appInfo.AppName, appInfo.Description)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111073", "应用更新失败:"+err.Error(), nil)
		return
	}
	c.ok(updateApp)
}

// DeleteApp 处理删除应用的请求。
// 该函数从请求参数中获取应用代码，并检查应用是否存在。如果应用存在，则删除应用。
// 如果应用不存在，返回相应的错误信息。
func (c *AppController) DeleteApp() {
	deleteAppCode := c.Ctx.Input.Param(":appCode")

	// 判断应用代码是否包含版本信息
	version := ""
	appCode := ""
	if strings.Contains(deleteAppCode, ".") {
		parts := strings.Split(deleteAppCode, ".")
		appCode = parts[0]
		version = strings.Join(parts[1:], ".")
		if !app.GetAppManagerInstance().IsAppVersionExist(appCode, version) {
			c.error("00111081", "删除失败，应用或当前版本不存在", nil)
			return
		}
	} else {
		appCode = deleteAppCode
		if !app.GetAppManagerInstance().IsAppDevExist(appCode) {
			c.error("00111082", "删除失败，应用不存在", nil)
			return
		}
	}

	// 删除应用
	err := app.GetAppManagerInstance().DeleteApp(appCode, version)
	if err != nil {
		return
	}
	c.ok(true)
}

// Release 处理发布应用的请求。
// 该函数从请求体中解析应用信息，并检查应用是否存在以及版本是否为空或已存在。
// 如果应用存在且版本合法，则发布应用。如果解析失败或应用不存在，返回相应的错误信息。
func (c *AppController) Publish() {
	var appInfo app.App
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111091", "应用发布失败", nil)
		return
	}

	// 检查应用是否存在
	if !app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111092", "应用不存在", nil)
		return
	}

	// 检查应用版本是否为空
	if appInfo.Version == "" {
		c.error("00111093", "应用版本不能为空", nil)
		return
	}

	// 检查应用版本是否已存在
	if app.GetAppManagerInstance().IsAppVersionExist(appInfo.AppCode, appInfo.Version) {
		c.error("00111094", "应用版本已存在", nil)
		return
	}

	// 发布应用
	err = app.GetAppManagerInstance().Publish(appInfo.AppCode, appInfo.Version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111095", "应用发布失败:"+err.Error(), nil)
		return
	}
	c.ok(true)
}

func (c *AppController) Release() {
	var appInfo app.App
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111091", "应用发布失败", nil)
		return
	}
	if appInfo.Version == "" {
		c.error("00111092", "应用版本不能为空", nil)
		return
	}
	if !app.GetAppManagerInstance().IsAppVersionExist(appInfo.AppCode, appInfo.Version) {
		c.error("00111093", "应用版本不存在", nil)
		return
	}
	err = app.GetAppManagerInstance().Release(appInfo.AppCode, appInfo.Version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111094", "应用版本切换失败:"+err.Error(), nil)
		return
	}
	c.ok(true)
}

func (c *AppController) Rollback() {
	var appInfo app.App
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111095", "应用发布失败", nil)
		return
	}
	if appInfo.Version == "" {
		c.error("00111096", "应用版本不能为空", nil)
		return
	}
	if !app.GetAppManagerInstance().IsAppVersionExist(appInfo.AppCode, appInfo.Version) {
		c.error("00111097", "应用版本不存在", nil)
		return
	}
	err = app.GetAppManagerInstance().Rollback(appInfo.AppCode, appInfo.Version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111098", "应用回滚失败:"+err.Error(), nil)
		return
	}
	c.ok(true)
}
