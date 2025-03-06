package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseVersion(versionStr string) (Version, error) {
	// 去掉前缀 "v" 或 "V"
	versionStr = strings.TrimPrefix(strings.ToLower(versionStr), "v")

	// 按 "." 分割版本号
	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return Version{}, errors.New("版本号格式不正确，必须是 vX.Y.Z 格式")
	}

	// 解析 Major, Minor, Patch
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, errors.New("Major 版本号必须是整数")
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, errors.New("Minor 版本号必须是整数")
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, errors.New("Patch 版本号必须是整数")
	}

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// CompareVersions 比较两个版本号
// 返回值：
// -1: version1 < version2
//
//	0: version1 == version2
//	1: version1 > version2
func CompareVersions(version1Str, version2Str string) (int, error) {
	version1, err := ParseVersion(version1Str)
	if err != nil {
		return 0, fmt.Errorf("解析版本号 %s 失败: %v", version1Str, err)
	}

	version2, err := ParseVersion(version2Str)
	if err != nil {
		return 0, fmt.Errorf("解析版本号 %s 失败: %v", version2Str, err)
	}

	// 比较 Major
	if version1.Major != version2.Major {
		if version1.Major < version2.Major {
			return -1, fmt.Errorf("版本号 %s 序列须高于 %s", version1Str, version2Str)
		} else {
			return 1, nil
		}
	}

	// 比较 Minor
	if version1.Minor != version2.Minor {
		if version1.Minor < version2.Minor {
			return -1, fmt.Errorf("版本号 %s 序列须高于 %s", version1Str, version2Str)
		} else {
			return 1, nil
		}
	}

	// 比较 Patch
	if version1.Patch != version2.Patch {
		if version1.Patch < version2.Patch {
			return -1, fmt.Errorf("版本号 %s 序列须高于 %s", version1Str, version2Str)
		} else {
			return 1, nil
		}
	}

	// 版本号完全相同
	return 0, fmt.Errorf("版本号 %s 序列须高于 %s", version1Str, version2Str)
}
