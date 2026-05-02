package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"tech-muyi-shenji/app"
)

// MenuController 处理菜单相关请求。
type MenuController struct {
	BaseController
}

type createMenuBody struct {
	AppCode    string `json:"appCode"`
	ParentCode string `json:"parentCode"`
	MenuCode   string `json:"menuCode"`
	MenuLabel  string `json:"menuLabel"`
	MenuPath   string `json:"menuPath"`
	MenuIcon   string `json:"menuIcon"`
	MenuType   string `json:"menuType"`
	Hide       bool   `json:"hide"`
	LinkUrl    string `json:"linkUrl"`
}

func (c *MenuController) CreateMenu() {
	var body createMenuBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		logs.Error(err.Error())
		c.error("00111110", "菜单创建失败", nil)
		return
	}
	if !app.ValidResourceCode(body.AppCode) {
		c.invalidAppCode()
		return
	}
	if body.ParentCode != "" && !app.ValidResourceCode(body.ParentCode) {
		c.error("00111050", "无效的父菜单编码", nil)
		return
	}
	if !app.ValidResourceCode(body.MenuCode) {
		c.invalidMenuCode()
		return
	}
	createMenu, err := app.GetMenuManagerInstance().CreateMenu(body.AppCode, body.ParentCode, body.MenuCode, body.MenuLabel, body.MenuPath, body.MenuIcon, body.MenuType, body.Hide, body.LinkUrl)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111112", "菜单创建失败", nil)
		return
	}
	c.ok(createMenu)
}

type appMenuCodeBody struct {
	AppCode  string `json:"appCode"`
	MenuCode string `json:"menuCode"`
}

func (c *MenuController) DeleteMenu() {
	var request appMenuCodeBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		logs.Error(err.Error())
		c.error("00111150", "请求体解析失败", nil)
		return
	}
	if !app.ValidResourceCode(request.AppCode) || !app.ValidResourceCode(request.MenuCode) {
		if !app.ValidResourceCode(request.AppCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}
	menus, err := app.GetMenuManagerInstance().DeleteMenu(request.AppCode, request.MenuCode)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111151", "菜单删除失败", nil)
		return
	}
	c.ok(menus)
}

type updateMenuBody struct {
	AppCode   string `json:"appCode"`
	MenuCode  string `json:"menuCode"`
	MenuLabel string `json:"menuLabel"`
	MenuIcon  string `json:"menuIcon"`
	Hide      bool   `json:"hide"`
	LinkUrl   string `json:"linkUrl"`
}

func (c *MenuController) UpdateMenu() {
	var body updateMenuBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		logs.Error(err.Error())
		c.error("00111152", "请求体解析失败", nil)
		return
	}
	if !app.ValidResourceCode(body.AppCode) || !app.ValidResourceCode(body.MenuCode) {
		if !app.ValidResourceCode(body.AppCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}
	updateMenu, err := app.GetMenuManagerInstance().UpdateMenu(body.AppCode, body.MenuCode, body.MenuLabel, body.MenuIcon, body.Hide, body.LinkUrl)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111153", "菜单更新失败", nil)
		return
	}
	c.ok(updateMenu)
}

func (c *MenuController) SortMenu() {
	var requestBody struct {
		AppCode string     `json:"appCode"`
		Menus   []app.Menu `json:"menus"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		logs.Error(err.Error())
		c.error("00111115", "菜单排序失败", nil)
		return
	}
	if !app.ValidResourceCode(requestBody.AppCode) {
		c.invalidAppCode()
		return
	}
	sortMenus, err := app.GetMenuManagerInstance().SortMenu(requestBody.AppCode, requestBody.Menus)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111116", "菜单排序失败", nil)
		return
	}
	c.ok(sortMenus)
}

type removeMenuBody struct {
	AppCode    string `json:"appCode"`
	MenuCode   string `json:"menuCode"`
	ParentCode string `json:"parentCode"`
}

func (c *MenuController) RemoveMenu() {
	var request removeMenuBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		logs.Error(err.Error())
		c.error("00111160", "请求体解析失败", nil)
		return
	}
	if !app.ValidResourceCode(request.AppCode) || !app.ValidResourceCode(request.MenuCode) {
		if !app.ValidResourceCode(request.AppCode) {
			c.invalidAppCode()
		} else {
			c.invalidMenuCode()
		}
		return
	}
	if request.ParentCode != "" && !app.ValidResourceCode(request.ParentCode) {
		c.error("00111050", "无效的父菜单编码", nil)
		return
	}
	if err := app.GetMenuManagerInstance().RemoveMenu(request.AppCode, request.MenuCode, request.ParentCode); err != nil {
		logs.Error(err.Error())
		c.error("00111161", "菜单移动失败", nil)
		return
	}
	c.ok(true)
}

func (c *MenuController) GetMenu() {
	appCode, version := appCodeVersionForDevAPI(c.Ctx.Input.Param(":appCode"))
	if !app.ValidResourceCode(appCode) {
		c.invalidAppCode()
		return
	}
	menus, err := app.GetMenuManagerInstance().GetMenu(appCode, version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111117", "菜单获取失败", nil)
		return
	}
	c.ok(menus)
}

func (c *MenuController) GetReleaseMenu() {
	appCode, version := appCodeVersionForReleaseAPI(c.Ctx.Input.Param(":appCode"))
	if !app.ValidResourceCode(appCode) {
		c.invalidAppCode()
		return
	}
	menus, err := app.GetMenuManagerInstance().GetMenu(appCode, version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111117", "菜单获取失败", nil)
		return
	}
	c.ok(menus)
}
