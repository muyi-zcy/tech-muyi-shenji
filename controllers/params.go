package controllers

import "strings"

// splitAppCodeVersion 将路由中的 :appCode 拆成 app 标识与版本后缀（第一段为 appCode，其后含点号整体为版本）。
func splitAppCodeVersion(raw string) (appCode, version string) {
	if raw == "" {
		return "", ""
	}
	i := strings.IndexByte(raw, '.')
	if i < 0 {
		return raw, ""
	}
	return raw[:i], raw[i+1:]
}

// appCodeVersionForDevAPI 管理端查询：无后缀表示开发态 dev。
func appCodeVersionForDevAPI(raw string) (appCode, version string) {
	appCode, v := splitAppCodeVersion(raw)
	if v == "" {
		return appCode, "dev"
	}
	return appCode, v
}

// appCodeVersionForReleaseAPI 运行态查询：无后缀表示当前发布（release），非 dev。
func appCodeVersionForReleaseAPI(raw string) (appCode, version string) {
	return splitAppCodeVersion(raw)
}
