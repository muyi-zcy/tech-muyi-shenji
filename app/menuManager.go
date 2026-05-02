package app

import (
	"fmt"
	"sync"
	"tech-muyi-shenji/utils"
)

type MenuManager struct{}

var (
	menuManagerInstance *MenuManager
	menuOnce            sync.Once
)

func GetMenuManagerInstance() *MenuManager {
	menuOnce.Do(func() {
		menuManagerInstance = &MenuManager{}
	})
	return menuManagerInstance
}

// 创建菜单
func (mm *MenuManager) CreateMenu(appCode string, parentCode string, menuCode string, menuLabel string, menuPath string, menuIcon string, menuType string, hide bool, linkUrl string) (*[]Menu, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !ValidResourceCode(menuCode) {
		return nil, fmt.Errorf("menuCode only allows numbers, lowercase letters, uppercase letters, underscores, and hyphens")
	}
	if parentCode != "" && !ValidResourceCode(parentCode) {
		return nil, fmt.Errorf("parentCode is invalid")
	}

	// menuType只能是MenuType内包含的枚举
	if parentCode != "" {
		if MenuImperative != MenuType(menuType) && MenuTypeLowCode != MenuType(menuType) && MenuTypeLink != MenuType(menuType) {
			return nil, fmt.Errorf("menuType only allows MenuImperative, MenuStaticPage, MenuTypeLink")
		}
	} else {
		if MenuImperative != MenuType(menuType) && MenuTypeLowCode != MenuType(menuType) && MenuTypeLink != MenuType(menuType) && Group != MenuType(menuType) {
			return nil, fmt.Errorf("menuType only allows MenuImperative, MenuStaticPage, MenuTypeLink")
		}
	}

	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with Code %s does not exist", appCode)
	}

	// 使用getAppMenu获取到对应的菜单json
	var menus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &menus)
	if err != nil {
		return nil, fmt.Errorf("failed to read menu info: %w", err)
	}

	existing := make(map[string]struct{})
	collectMenuCodesDeep(menus, existing)
	if _, dup := existing[menuCode]; dup {
		return nil, fmt.Errorf("menu with Code %s already exists", menuCode)
	}

	if parentCode == "" {
		menus = append(menus, Menu{
			MenuCode:  menuCode,
			MenuLabel: menuLabel,
			MenuPath:  menuPath,
			MenuIcon:  menuIcon,
			MenuType:  menuType,
			Hide:      hide,
			Children:  []Menu{},
			LinkUrl:   linkUrl,
		})
	} else {
		if Group == MenuType(menuType) {
			return nil, fmt.Errorf("child menu`s menuType only allows MenuImperative, MenuStaticPage, MenuTypeLink")
		}
		child := Menu{
			MenuCode:  menuCode,
			MenuLabel: menuLabel,
			MenuPath:  menuPath,
			MenuIcon:  menuIcon,
			MenuType:  menuType,
			Hide:      hide,
			LinkUrl:   linkUrl,
			Children:  []Menu{},
		}
		var ok bool
		menus, ok = appendChildDeep(menus, parentCode, child)
		if !ok {
			return nil, fmt.Errorf("parent menu with Code %s does not exist", parentCode)
		}
	}
	err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
	if err != nil {
		return nil, fmt.Errorf("failed to save menu info: %w", err)
	}
	return &menus, nil
}

// 删除菜单
func (mm *MenuManager) DeleteMenu(appCode string, menuCode string) (*[]Menu, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with Code %s does not exist", appCode)
	}
	var menus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &menus)
	if err != nil {
		return nil, fmt.Errorf("failed to read menu info: %w", err)
	}
	var found bool
	menus, found, err = deleteMenuDeep(menus, menuCode)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("menu with Code %s does not exist", menuCode)
	}
	err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
	if err != nil {
		return nil, fmt.Errorf("failed to save menu info: %w", err)
	}
	pagePath := getAppDevPage(appCode, menuCode)
	if err := utils.RemoveFile(pagePath); err != nil {
		return nil, err
	}
	return &menus, nil
}

// 修改菜单名称、icon、显示隐藏
func (mm *MenuManager) UpdateMenu(appCode string, menuCode string, menuLabel string, menuIcon string, hide bool, linkUrl string) (*[]Menu, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)
	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with Code %s does not exist", appCode)
	}
	var menus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &menus)
	if err != nil {
		return nil, fmt.Errorf("failed to read menu info: %w", err)
	}

	var ok bool
	menus, ok = updateMenuDeep(menus, menuCode, menuLabel, menuIcon, hide, linkUrl)
	if !ok {
		return nil, fmt.Errorf("menu with Code %s does not exist", menuCode)
	}
	err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
	if err != nil {
		return nil, fmt.Errorf("failed to save menu info: %w", err)
	}
	return &menus, nil
}

// 菜单排序，给出所有父子结构菜单结构排序（包括子菜单），需要保证新的结构不缺少、不新增
func (mm *MenuManager) SortMenu(appCode string, newMenus []Menu) (*[]Menu, error) {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)
	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return nil, fmt.Errorf("app with Code %s does not exist", appCode)
	}
	var oldMenus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &oldMenus)
	if err != nil {
		return nil, fmt.Errorf("failed to read menu info: %w", err)
	}

	// 检查是否缺少或新增菜单
	if err := mm.validateMenuConsistency(oldMenus, newMenus); err != nil {
		return nil, err
	}

	// 保存新的菜单顺序,并使用旧菜单里的菜单属性填充所有菜单的信息
	oldMenuMap := make(map[string]Menu)
	mm.buildMenuMap(oldMenus, oldMenuMap)
	for i, menu := range newMenus {
		if oldMenu, ok := oldMenuMap[menu.MenuCode]; ok {
			newMenus[i].MenuLabel = oldMenu.MenuLabel
			newMenus[i].MenuIcon = oldMenu.MenuIcon
			newMenus[i].Hide = oldMenu.Hide
			newMenus[i].MenuPath = oldMenu.MenuPath
			newMenus[i].MenuType = oldMenu.MenuType
			newMenus[i].LinkUrl = oldMenu.LinkUrl
		}
		if menu.Children != nil && len(menu.Children) > 0 {
			for j, child := range menu.Children {
				if oldChild, ok := oldMenuMap[child.MenuCode]; ok {
					newMenus[i].Children[j].MenuLabel = oldChild.MenuLabel
					newMenus[i].Children[j].MenuIcon = oldChild.MenuIcon
					newMenus[i].Children[j].Hide = oldChild.Hide
					newMenus[i].Children[j].MenuPath = oldChild.MenuPath
					newMenus[i].Children[j].MenuType = oldChild.MenuType
					newMenus[i].Children[j].LinkUrl = oldChild.LinkUrl
				}
			}
		}
	}
	err = utils.SaveJSONToFile(newMenus, getAppDevMenu(appCode))
	if err != nil {
		return nil, fmt.Errorf("failed to save menu info: %w", err)
	}

	return &newMenus, nil
}

// 检查新旧菜单是否一致，确保不缺少、不新增
func (mm *MenuManager) validateMenuConsistency(oldMenus, newMenus []Menu) error {
	oldMenuMap := make(map[string]Menu)
	mm.buildMenuMap(oldMenus, oldMenuMap)

	newMenuMap := make(map[string]Menu)
	mm.buildMenuMap(newMenus, newMenuMap)

	// 检查是否有新增的菜单
	for code := range newMenuMap {
		if _, exists := oldMenuMap[code]; !exists {
			return fmt.Errorf("new menu with Code %s is added", code)
		}
	}

	// 检查是否有缺少的菜单
	for code := range oldMenuMap {
		if _, exists := newMenuMap[code]; !exists {
			return fmt.Errorf("menu with Code %s is missing", code)
		}
	}

	return nil
}

// 递归构建菜单的映射关系
func (mm *MenuManager) buildMenuMap(menus []Menu, menuMap map[string]Menu) {
	for _, menu := range menus {
		menuMap[menu.MenuCode] = menu
		if len(menu.Children) > 0 {
			mm.buildMenuMap(menu.Children, menuMap)
		}
	}
}

func (mm *MenuManager) SaveMenuPage(appCode string, menuCode string, config interface{}) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return fmt.Errorf("app with Code %s does not exist", appCode)
	}
	// 获取菜单信息
	var menus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &menus)
	if err != nil {
		return fmt.Errorf("failed to read menu info: %w", err)
	}
	menuMap := make(map[string]Menu)
	mm.buildMenuMap(menus, menuMap)
	menu, ok := menuMap[menuCode]
	if !ok {
		return fmt.Errorf("menu with Code %s does not exist", menuCode)
	}
	if MenuType(menu.MenuType) != MenuTypeLowCode {
		return fmt.Errorf("menu with Code %s is not lowCode type", menuCode)
	}

	utils.CopyFile(getAppDevPage(appCode, menuCode), getAppHistoryPage(appCode, menuCode))
	err = utils.SaveJSONToFile(config, getAppDevPage(appCode, menuCode))
	if err != nil {
		return fmt.Errorf("failed to save menu page: %w", err)
	}
	return nil
}

func (mm *MenuManager) GetMenu(appCode string, version string) (interface{}, error) {
	am := GetAppManagerInstance()
	switch {
	case version == "":
		if !am.IsAppReleaseExist(appCode) {
			return "", fmt.Errorf("app release for Code %s does not exist", appCode)
		}
	case version == "dev":
		if !am.IsAppDevExist(appCode) {
			return "", fmt.Errorf("app with Code %s does not exist", appCode)
		}
	default:
		if !am.IsAppVersionExist(appCode, version) {
			return "", fmt.Errorf("app version for Code %s does not exist", appCode)
		}
	}

	var config interface{}
	var err error
	if version == "" {
		_, err = utils.ReadJSONFromFile(getAppReleaseMenu(appCode), &config)
	} else if version == "dev" {
		_, err = utils.ReadJSONFromFile(getAppDevMenu(appCode), &config)
	} else {
		_, err = utils.ReadJSONFromFile(getAppVersionMenu(appCode, version), &config)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read menu page: %w", err)
	}
	return config, nil
}

func (mm *MenuManager) RemoveMenu(appCode string, menuCode string, parentCode string) error {
	appLock := GetAppLock()
	appLock.Lock(appCode)
	defer appLock.Unlock(appCode)

	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return fmt.Errorf("app with Code %s does not exist", appCode)
	}

	var menus []Menu
	_, err := utils.ReadJSONFromFile(getAppDevMenu(appCode), &menus)
	if err != nil {
		return fmt.Errorf("failed to read menu info: %w", err)
	}

	if parentCode != "" && !ValidResourceCode(parentCode) {
		return fmt.Errorf("parentCode is invalid")
	}

	var targetMenu Menu
	var found bool
	menus, targetMenu, found = extractMenuByCode(menus, menuCode)
	if !found {
		return fmt.Errorf("menu with Code %s does not exist", menuCode)
	}

	if parentCode == "" {
		menus = append(menus, targetMenu)
	} else {
		var ok bool
		menus, ok = appendChildDeep(menus, parentCode, targetMenu)
		if !ok {
			return fmt.Errorf("parent menu with Code %s does not exist", parentCode)
		}
	}

	// 保存更新后的菜单
	err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
	if err != nil {
		return fmt.Errorf("failed to save menu info: %w", err)
	}

	return nil
}

func (mm *MenuManager) GetMenuPage(appCode string, version string, menuCode string) (interface{}, error) {
	am := GetAppManagerInstance()
	switch {
	case version == "":
		if !am.IsAppReleaseExist(appCode) {
			return "", fmt.Errorf("app release for Code %s does not exist", appCode)
		}
	case version == "dev":
		if !am.IsAppDevExist(appCode) {
			return "", fmt.Errorf("app with Code %s does not exist", appCode)
		}
	default:
		if !am.IsAppVersionExist(appCode, version) {
			return "", fmt.Errorf("app version for Code %s does not exist", appCode)
		}
	}

	var config interface{}
	var err error
	if version == "" {
		_, err = utils.ReadJSONFromFile(getAppReleasePage(appCode, menuCode), &config)
	} else if version == "dev" {
		_, err = utils.ReadJSONFromFile(getAppDevPage(appCode, menuCode), &config)
	} else {
		_, err = utils.ReadJSONFromFile(getAppVersionPage(appCode, version, menuCode), &config)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read menu page: %w", err)
	}
	return config, nil
}

func collectMenuCodesDeep(menus []Menu, out map[string]struct{}) {
	for _, m := range menus {
		out[m.MenuCode] = struct{}{}
		collectMenuCodesDeep(m.Children, out)
	}
}

func appendChildDeep(menus []Menu, parentCode string, child Menu) ([]Menu, bool) {
	for i := range menus {
		if menus[i].MenuCode == parentCode {
			menus[i].Children = append(menus[i].Children, child)
			return menus, true
		}
		if len(menus[i].Children) > 0 {
			updated, ok := appendChildDeep(menus[i].Children, parentCode, child)
			if ok {
				menus[i].Children = updated
				return menus, true
			}
		}
	}
	return menus, false
}

func deleteMenuDeep(menus []Menu, menuCode string) ([]Menu, bool, error) {
	for i, menu := range menus {
		if menu.MenuCode == menuCode {
			if len(menu.Children) > 0 {
				return nil, false, fmt.Errorf("menu with Code %s has child menus, cannot delete", menuCode)
			}
			out := append(menus[:i], menus[i+1:]...)
			return out, true, nil
		}
		if len(menu.Children) > 0 {
			newChildren, found, err := deleteMenuDeep(menu.Children, menuCode)
			if err != nil {
				return nil, false, err
			}
			if found {
				menus[i].Children = newChildren
				return menus, true, nil
			}
		}
	}
	return menus, false, nil
}

func updateMenuDeep(menus []Menu, menuCode string, menuLabel, menuIcon string, hide bool, linkUrl string) ([]Menu, bool) {
	for i := range menus {
		if menus[i].MenuCode == menuCode {
			if menuLabel != "" {
				menus[i].MenuLabel = menuLabel
			}
			if menuIcon != "" {
				menus[i].MenuIcon = menuIcon
			}
			if linkUrl != "" {
				menus[i].LinkUrl = linkUrl
			}
			menus[i].Hide = hide
			return menus, true
		}
		if len(menus[i].Children) > 0 {
			updated, ok := updateMenuDeep(menus[i].Children, menuCode, menuLabel, menuIcon, hide, linkUrl)
			if ok {
				menus[i].Children = updated
				return menus, true
			}
		}
	}
	return menus, false
}

func extractMenuByCode(menus []Menu, menuCode string) (newMenus []Menu, target Menu, found bool) {
	for i, menu := range menus {
		if menu.MenuCode == menuCode {
			target = menu
			newMenus = append(menus[:i], menus[i+1:]...)
			return newMenus, target, true
		}
		if len(menu.Children) > 0 {
			newChildren, t, ok := extractMenuByCode(menu.Children, menuCode)
			if ok {
				menu.Children = newChildren
				menus[i] = menu
				return menus, t, true
			}
		}
	}
	return menus, Menu{}, false
}
