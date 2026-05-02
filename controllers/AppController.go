package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"tech-muyi-shenji/app"
)

// AppController 处理与应用相关的 HTTP 请求。
type AppController struct {
	BaseController
}

// GetAppInfo 获取开发态或指定历史版本的应用信息。
func (c *AppController) GetAppInfo() {
	appCode, version := appCodeVersionForDevAPI(c.Ctx.Input.Param(":appCode"))
	if !app.ValidResourceCode(appCode) {
		c.invalidAppCode()
		return
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

// GetReleaseAppInfo 获取已发布或指定版本的应用信息（运行态）。
func (c *AppController) GetReleaseAppInfo() {
	appCode, version := appCodeVersionForReleaseAPI(c.Ctx.Input.Param(":appCode"))
	if !app.ValidResourceCode(appCode) {
		c.invalidAppCode()
		return
	}
	var err error
	var appInfo *app.App
	if version == "" {
		if !app.GetAppManagerInstance().IsAppReleaseExist(appCode) {
			c.error("00111051", "应用不存在", nil)
			return
		}
		appInfo, err = app.GetAppManagerInstance().GetAppReleaseInfo(appCode)
	} else if version == "dev" {
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

// CreateApp 创建应用（开发态）。
func (c *AppController) CreateApp() {
	var appInfo app.App
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &appInfo); err != nil {
		logs.Error(err.Error())
		c.error("00111061", "应用创建失败", nil)
		return
	}
	if !app.ValidResourceCode(appInfo.AppCode) {
		c.invalidAppCode()
		return
	}
	if app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111062", "应用已存在", nil)
		return
	}
	createApp, err := app.GetAppManagerInstance().CreateApp(appInfo.AppCode, appInfo.AppName, appInfo.Version, appInfo.Description)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111063", "应用创建失败", nil)
		return
	}
	c.ok(createApp)
}

// UpdateApp 更新开发态应用信息。
func (c *AppController) UpdateApp() {
	var appInfo app.App
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &appInfo); err != nil {
		logs.Error(err.Error())
		c.error("00111071", "应用更新失败", nil)
		return
	}
	if !app.ValidResourceCode(appInfo.AppCode) {
		c.invalidAppCode()
		return
	}
	if !app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111072", "应用不存在", nil)
		return
	}
	updateApp, err := app.GetAppManagerInstance().UpdateApp(appInfo.AppCode, appInfo.AppName, appInfo.Description)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111073", "应用更新失败", nil)
		return
	}
	c.ok(updateApp)
}

// DeleteApp 删除开发态应用或指定版本快照。
func (c *AppController) DeleteApp() {
	deleteRaw := c.Ctx.Input.Param(":appCode")
	appCode, version := splitAppCodeVersion(deleteRaw)
	if !app.ValidResourceCode(appCode) {
		c.invalidAppCode()
		return
	}
	if version != "" {
		if !app.GetAppManagerInstance().IsAppVersionExist(appCode, version) {
			c.error("00111081", "删除失败，应用或当前版本不存在", nil)
			return
		}
	} else {
		if !app.GetAppManagerInstance().IsAppDevExist(appCode) {
			c.error("00111082", "删除失败，应用不存在", nil)
			return
		}
	}
	if err := app.GetAppManagerInstance().DeleteApp(appCode, version); err != nil {
		logs.Error(err.Error())
		c.error("00111083", "删除失败", nil)
		return
	}
	c.ok(true)
}

// Publish 将开发态发布为新版本并更新线上目录。
func (c *AppController) Publish() {
	var appInfo app.App
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &appInfo); err != nil {
		logs.Error(err.Error())
		c.error("00111091", "应用发布失败", nil)
		return
	}
	if !app.ValidResourceCode(appInfo.AppCode) {
		c.invalidAppCode()
		return
	}
	if !app.GetAppManagerInstance().IsAppDevExist(appInfo.AppCode) {
		c.error("00111092", "应用不存在", nil)
		return
	}
	if appInfo.Version == "" {
		c.error("00111093", "应用版本不能为空", nil)
		return
	}
	if app.GetAppManagerInstance().IsAppVersionExist(appInfo.AppCode, appInfo.Version) {
		c.error("00111094", "应用版本已存在", nil)
		return
	}
	if err := app.GetAppManagerInstance().Publish(appInfo.AppCode, appInfo.Version); err != nil {
		logs.Error(err.Error())
		c.error("00111095", "应用发布失败", nil)
		return
	}
	c.ok(true)
}

// Release 将指定历史版本切换为当前发布内容。
func (c *AppController) Release() {
	var appInfo app.App
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &appInfo); err != nil {
		logs.Error(err.Error())
		c.error("00111091", "应用发布失败", nil)
		return
	}
	if !app.ValidResourceCode(appInfo.AppCode) {
		c.invalidAppCode()
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
	if err := app.GetAppManagerInstance().Release(appInfo.AppCode, appInfo.Version); err != nil {
		logs.Error(err.Error())
		c.error("00111094", "应用版本切换失败", nil)
		return
	}
	c.ok(true)
}

// Rollback 用指定历史版本覆盖开发态。
func (c *AppController) Rollback() {
	var appInfo app.App
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &appInfo); err != nil {
		logs.Error(err.Error())
		c.error("00111095", "应用发布失败", nil)
		return
	}
	if !app.ValidResourceCode(appInfo.AppCode) {
		c.invalidAppCode()
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
	if err := app.GetAppManagerInstance().Rollback(appInfo.AppCode, appInfo.Version); err != nil {
		logs.Error(err.Error())
		c.error("00111098", "应用回滚失败", nil)
		return
	}
	c.ok(true)
}
