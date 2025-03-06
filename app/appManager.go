package app

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"tech-muyi-shenji/common"
	"tech-muyi-shenji/utils"
)

type AppManager struct {
}

var (
	appManagerInstance *AppManager
	appOnce            sync.Once
)

func GetAppManagerInstance() *AppManager {
	appOnce.Do(func() {
		appManagerInstance = &AppManager{}
	})
	return appManagerInstance
}

func (am *AppManager) IsAppDevExist(appCode string) bool {
	appPath := getAppDevPath(appCode)
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (am *AppManager) IsAppVersionExist(appCode string, version string) bool {
	appPath := getAppVersionPath(appCode, version)
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (am *AppManager) IsAppReleaseExist(appCode string) bool {
	appPath := getAppReleasePath(appCode)
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (am *AppManager) GetAppDevInfo(appCode string) (*App, error) {
	appInfoFilePath := getAppDevPath(appCode)
	if _, err := os.Stat(appInfoFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("app with ID %s does not exist", appCode)
	}

	var app App
	_, err := utils.ReadJSONFromFile(appInfoFilePath, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to read app info: %w", err)
	}
	return &app, nil
}

func (am *AppManager) GetAppVersionInfo(appCode string, version string) (*App, error) {
	appInfoFilePath := getAppVersionPath(appCode, version)
	if _, err := os.Stat(appInfoFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("app with ID %s %s does not exist", appCode, version)
	}

	var app App
	_, err := utils.ReadJSONFromFile(appInfoFilePath, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to read app info: %w", err)
	}
	return &app, nil
}

func (am *AppManager) GetAppReleaseInfo(appCode string) (*App, error) {
	appInfoFilePath := getAppReleasePath(appCode)
	if _, err := os.Stat(appInfoFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("app with ID %s  does not exist", appCode)
	}

	var app App
	_, err := utils.ReadJSONFromFile(appInfoFilePath, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to read app info: %w", err)
	}
	return &app, nil
}

func (am *AppManager) CreateApp(appCode string, appName string, version string, description string) (*App, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	// appCode只允许包含数字、大小写字母、下划线、横杠,使用正则
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(appCode)
	if !valid {
		return nil, fmt.Errorf("appCode only allows numbers, lowercase letters, uppercase letters, underscores, and hyphens")
	}

	if am.IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with ID %s already exists", appCode)
	}

	appId, _ := common.NextId()
	app := App{
		AppId:       strconv.FormatInt(appId, 10),
		AppCode:     appCode,
		AppName:     appName,
		Version:     "v0.0.1",
		Description: description,
	}
	err := utils.SaveJSONToFile(app, getAppDevPath(appCode))
	if err != nil {
		return nil, err
	}

	// 保存一个空数组，用于存储菜单
	err = utils.SaveJSONToFile([]Menu{}, getAppDevMenu(appCode))
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (am *AppManager) UpdateApp(appCode string, appName string, description string) (*App, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !am.IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with ID %s does not exist", appCode)
	}

	var oldapp App
	_, err := utils.ReadJSONFromFile(getAppDevPath(appCode), &oldapp)
	if err != nil {
		return nil, fmt.Errorf("failed to read app info: %w", err)
	}

	if appName != "" {
		oldapp.AppName = appName
	}
	if description != "" {
		oldapp.Description = description
	}

	err = utils.SaveJSONToFile(oldapp, getAppDevPath(appCode))
	if err != nil {
		return nil, err
	}
	return &oldapp, nil
}

func (am *AppManager) DeleteApp(appCode string, version string) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	// 删除文件夹
	if version == "" {
		if am.IsAppDevExist(appCode) {
			err := os.RemoveAll(getAppDevRootPath(appCode))
			if err != nil {
				return fmt.Errorf("failed to delete app: %w", err)
			}
			return nil
		}
	}
	if am.IsAppVersionExist(appCode, version) {
		err := os.RemoveAll(getAppVersionRootPath(appCode, version))
		if err != nil {
			return fmt.Errorf("failed to delete app version: %w", err)
		}
	}
	return nil
}

func (am *AppManager) Publish(appCode string, version string) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	appInfo, err := am.GetAppDevInfo(appCode)
	if err != nil {
		return fmt.Errorf("failed to get app info: %w", err)
	}
	if am.IsAppVersionExist(appCode, version) {
		return fmt.Errorf("app version %s already exists", version)
	}
	_, err = utils.CompareVersions(version, appInfo.Version)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = utils.CopyDirSrc(getAppDevRootPath(appCode), getAppReleaseRootPath(appCode))
	if err != nil {
		return fmt.Errorf("failed to publish app: %w", err)
	}

	err = utils.CopyDirSrc(getAppDevRootPath(appCode), getAppVersionRootPath(appCode, appInfo.Version))
	if err != nil {
		return fmt.Errorf("failed to publish app: %w", err)
	}

	appInfo.Version = version
	err = utils.SaveJSONToFile(appInfo, getAppDevPath(appCode))
	if err != nil {
		return fmt.Errorf("failed to save app info: %w", err)
	}
	return nil
}

func (am *AppManager) Release(appCode string, version string) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !am.IsAppVersionExist(appCode, version) {
		return fmt.Errorf("app version %s does not exist", version)
	}

	err := utils.CopyDirSrc(getAppVersionRootPath(appCode, version), getAppReleaseRootPath(appCode))
	if err != nil {
		return fmt.Errorf("failed to rollback app: %w", err)
	}
	return nil
}

func (am *AppManager) Rollback(appCode string, version string) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !am.IsAppVersionExist(appCode, version) {
		return fmt.Errorf("app version %s does not exist", version)
	}
	devAppInfo, _ := am.GetAppDevInfo(appCode)
	devVersion := devAppInfo.Version

	err := utils.CopyDirSrc(getAppVersionRootPath(appCode, version), getAppDevRootPath(appCode))
	if err != nil {
		return fmt.Errorf("failed to rollback app: %w", err)
	}

	newDevAppInfo, _ := am.GetAppDevInfo(appCode)
	newDevAppInfo.Version = devVersion
	err = utils.SaveJSONToFile(newDevAppInfo, getAppDevPath(appCode))
	if err != nil {
		return fmt.Errorf("failed to save app info: %w", err)
	}
	return nil
}
