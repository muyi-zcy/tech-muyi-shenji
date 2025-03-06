package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"tech-muyi-shenji/app"
)

// AppController 是一个控制器，用于处理与应用程序相关的请求。
type MenuController struct {
	BaseController
}

// 创建分组
func (c *MenuController) CreateMenu() {
	type CreateMenuRequest struct {
		AppCode    string `json:"appCode"`
		ParentCode string `json:"parentCode"`
	}
	var request CreateMenuRequest
	appData := c.Ctx.Input.RequestBody
	err := json.Unmarshal(appData, &request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111110", "菜单创建失败", nil)
		return
	}
	appCode := request.AppCode
	parentCode := request.ParentCode

	var menuInfo app.Menu
	data := c.Ctx.Input.RequestBody
	err = json.Unmarshal(data, &menuInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111111", "菜单创建失败", nil)
		return
	}
	createMenu, err := app.GetMenuManagerInstance().CreateMenu(appCode, parentCode, menuInfo.MenuCode, menuInfo.MenuLabel, menuInfo.MenuPath, menuInfo.MenuIcon, menuInfo.MenuType, menuInfo.Hide, menuInfo.LinkUrl)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111112", "菜单创建失败:"+err.Error(), nil)
		return
	}
	c.ok(createMenu)
}

// 删除分组
func (c *MenuController) DeleteMenu() {
	type CreateMenuRequest struct {
		AppCode  string `json:"appCode"`
		MenuCode string `json:"menuCode"`
	}
	var request CreateMenuRequest
	appData := c.Ctx.Input.RequestBody
	err := json.Unmarshal(appData, &request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111113", "菜单创建失败", nil)
		return
	}
	appCode := request.AppCode
	menuCode := request.MenuCode

	menus, err := app.GetMenuManagerInstance().DeleteMenu(appCode, menuCode)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111114", "菜单删除失败:"+err.Error(), nil)
		return
	}
	c.ok(menus)
}

// 修改分组
func (c *MenuController) UpdateMenu() {
	type CreateMenuRequest struct {
		AppCode  string `json:"appCode"`
		MenuCode string `json:"menuCode"`
	}
	var request CreateMenuRequest
	appData := c.Ctx.Input.RequestBody
	err := json.Unmarshal(appData, &request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111113", "菜单创建失败", nil)
		return
	}
	appCode := request.AppCode
	menuCode := request.MenuCode

	var menuInfo app.Menu
	data := c.Ctx.Input.RequestBody
	err = json.Unmarshal(data, &menuInfo)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111114", "菜单更新失败", nil)
		return
	}
	updateMenu, err := app.GetMenuManagerInstance().UpdateMenu(appCode, menuCode, menuInfo.MenuLabel, menuInfo.MenuIcon, menuInfo.Hide, menuInfo.LinkUrl)
	if err != nil {
		logs.Error(err.Error())
	}
	c.ok(updateMenu)
}

// 排序
func (c *MenuController) SortMenu() {
	var requestBody struct {
		AppCode string     `json:"appCode"`
		Menus   []app.Menu `json:"menus"`
	}
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &requestBody)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111115", "菜单排序失败", nil)
		return
	}
	appCode := requestBody.AppCode
	newMenus := requestBody.Menus
	sortMenus, err := app.GetMenuManagerInstance().SortMenu(appCode, newMenus)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111116", "菜单排序失败:"+err.Error(), nil)
		return
	}
	c.ok(sortMenus)
}

func (c *MenuController) RemoveMenu() {
	type CreateMenuRequest struct {
		AppCode    string `json:"appCode"`
		MenuCode   string `json:"menuCode"`
		ParentCode string `json:"parentCode"`
	}
	var request CreateMenuRequest
	appData := c.Ctx.Input.RequestBody
	err := json.Unmarshal(appData, &request)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111115", "菜单排序失败", nil)
		return
	}
	err = app.GetMenuManagerInstance().RemoveMenu(request.AppCode, request.MenuCode, request.ParentCode)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111115", "菜单排序失败:"+err.Error(), nil)
		return
	}
	c.ok(true)
}

func (c *MenuController) GetMenu() {
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

	menus, err := app.GetMenuManagerInstance().GetMenu(appCode, version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111117", "菜单获取失败", nil)
		return
	}
	c.ok(menus)
}

func (c *MenuController) GetReleaseMenu() {
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
	menus, err := app.GetMenuManagerInstance().GetMenu(appCode, version)
	if err != nil {
		logs.Error(err.Error())
		c.error("00111117", "菜单获取失败", nil)
		return
	}
	c.ok(menus)
}
