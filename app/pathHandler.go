package app

import (
	"fmt"
	"tech-muyi-shenji/utils"
	"time"
)

const rootPath = "data/"
const devPath = "data/dev/"
const releasePath = "data/app/"
const versionPath = "data/version/"
const tempPath = "data/temp/"

func getAppDevRootPath(appCode string) string {
	return utils.JoinPaths(devPath, appCode)
}

func getAppDevPath(appCode string) string {
	return utils.JoinPaths(devPath, appCode, "app.json")
}

func getAppReleaseRootPath(appCode string) string {
	return utils.JoinPaths(releasePath, appCode)
}

func getAppVersionRootPath(appCode string, version string) string {
	return utils.JoinPaths(versionPath, appCode+"."+version)
}

func getAppVersionPath(appCode string, version string) string {
	return utils.JoinPaths(versionPath, appCode+"."+version, "app.json")
}

func getAppReleasePath(appCode string) string {
	return utils.JoinPaths(releasePath, appCode, "app.json")
}

func getAppDevMenu(appCode string) string {
	return utils.JoinPaths(devPath, appCode, "menu.json")
}

func getAppVersionMenu(appCode string, version string) string {
	return utils.JoinPaths(versionPath, appCode+"."+version, "menu.json")
}

func getAppReleaseMenu(appCode string) string {
	return utils.JoinPaths(releasePath, appCode, "menu.json")
}

func getAppDevPage(appCode string, menuCode string) string {
	return utils.JoinPaths(devPath, appCode, "pages", menuCode+".json")
}

func getAppVersionPage(appCode string, version string, menuCode string) string {
	return utils.JoinPaths(versionPath, appCode+"."+version, "pages", menuCode+".json")
}

func getAppReleasePage(appCode string, menuCode string) string {
	return utils.JoinPaths(releasePath, appCode, "pages", menuCode+".json")
}

func getAppHistoryPage(appCode string, menuCode string) string {
	// 实现获取时间戳
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	return utils.JoinPaths(tempPath, appCode, "pages", menuCode+"."+timestamp+".json")
}
