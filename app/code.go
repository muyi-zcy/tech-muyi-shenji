package app

import "regexp"

var resourceCodePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidResourceCode 校验 appCode、menuCode、parentCode 等，与创建应用时的规则一致，避免路径与文件名异常字符。
func ValidResourceCode(code string) bool {
	return code != "" && resourceCodePattern.MatchString(code)
}
