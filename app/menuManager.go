package app

import (
	"fmt"
	"regexp"
	"sync"
	"tech-muyi-shenji/utils"
)

type MenuManager struct {
	mu sync.RWMutex // 读写锁
}

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

	// 判断menuCode是否满足规则：只能是大小写字母
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(menuCode)
	if !valid {
		return nil, fmt.Errorf("menuCode only allows numbers, lowercase letters, uppercase letters, underscores, and hyphens")
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

	// menuCode不可以重复，包括所有的子菜单
	for _, menu := range menus {
		if menu.MenuCode == menuCode {
			return nil, fmt.Errorf("menu with Code %s already exists", menuCode)
		}
		if menu.Children != nil && len(menu.Children) > 0 {
			for _, child := range menu.Children {
				if child.MenuCode == menuCode {
					return nil, fmt.Errorf("menu with Code %s already exists", menuCode)
				}
			}
		}
	}

	if parentCode == "" {
		// 如果parentCode 为空，则新增到根菜单
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
		// 如果parentCode 不为空，则新增到对应菜单的子菜单
		flag := false
		for i, menu := range menus {
			if menu.MenuCode == parentCode {
				menus[i].Children = append(menu.Children, Menu{
					MenuCode:  menuCode,
					MenuLabel: menuLabel,
					MenuPath:  menuPath,
					MenuIcon:  menuIcon,
					MenuType:  menuType,
					Hide:      hide,
					LinkUrl:   linkUrl,
					Children:  []Menu{},
				})
				flag = true
				break
			}
		}
		if !flag {
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
	// 遍历删除对应的menuCode，如果menuCode对应的子菜单未清空，则无法删除
	flag := false
	for i, menu := range menus {
		if menu.MenuCode == menuCode {
			if len(menu.Children) > 0 {
				return nil, fmt.Errorf("menu with Code %s has child menus, cannot delete", menuCode)
			}
			menus = append(menus[:i], menus[i+1:]...)
			flag = true
		}
		if menu.Children != nil && len(menu.Children) > 0 {
			for j, child := range menu.Children {
				if child.MenuCode == menuCode {
					menus[i].Children = append(menu.Children[:j], menu.Children[j+1:]...)
					flag = true
					break
				}
			}
		}
		if flag {
			break
		}
	}
	if flag {
		err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
		if err != nil {
			return nil, fmt.Errorf("failed to save menu info: %w", err)
		}
		pagePath := getAppDevPage(appCode, menuCode)
		utils.RemoveFile(pagePath)
		return &menus, nil
	}
	return nil, fmt.Errorf("menu with Code %s does not exist", menuCode)
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

	flag := false
	for i, menu := range menus {
		if menu.MenuCode == menuCode {
			// 判断是否为空
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
			flag = true
		}
		if menu.Children != nil && len(menu.Children) > 0 {
			for j, child := range menu.Children {
				if child.MenuCode == menuCode {
					if menuLabel != "" {
						menus[i].Children[j].MenuLabel = menuLabel
					}
					if menuIcon != "" {
						menus[i].Children[j].MenuIcon = menuIcon
					}
					if linkUrl != "" {
						menus[i].Children[j].LinkUrl = linkUrl
					}
					menus[i].Children[j].Hide = hide
					flag = true
					break
				}
			}
		}
		if flag {
			break
		}
	}
	if flag {
		err = utils.SaveJSONToFile(menus, getAppDevMenu(appCode))
		if err != nil {
			return nil, fmt.Errorf("failed to save menu info: %w", err)
		}
		return &menus, nil
	}
	return nil, fmt.Errorf("menu with Code %s does not exist", menuCode)
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
	err = utils.SaveJSONToFile(config, getAppDevPage(appCode, menuCode))
	if err != nil {
		return fmt.Errorf("failed to save menu page: %w", err)
	}
	return nil
}

func (mm *MenuManager) GetMenu(appCode string, version string) (interface{}, error) {
	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return "", fmt.Errorf("app with Code %s does not exist", appCode)
	}

	//  读取配置，返回字符串
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

	// 查找目标菜单
	var targetMenu Menu
	found := false
	for i, menu := range menus {
		if menu.MenuCode == menuCode {
			targetMenu = menu
			menus = append(menus[:i], menus[i+1:]...)
			found = true
			break
		}
		if menu.Children != nil && len(menu.Children) > 0 {
			for j, child := range menu.Children {
				if child.MenuCode == menuCode {
					targetMenu = child
					menus[i].Children = append(menu.Children[:j], menu.Children[j+1:]...)
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("menu with Code %s does not exist", menuCode)
	}

	// 移动菜单
	if parentCode == "" {
		// 移动到根节点
		menus = append(menus, targetMenu)
	} else {
		// 移动到指定父菜单下
		parentFound := false
		for i, menu := range menus {
			if menu.MenuCode == parentCode {
				menus[i].Children = append(menu.Children, targetMenu)
				parentFound = true
				break
			}
			if menu.Children != nil && len(menu.Children) > 0 {
				for j, child := range menu.Children {
					if child.MenuCode == parentCode {
						menus[i].Children[j].Children = append(child.Children, targetMenu)
						parentFound = true
						break
					}
				}
			}
			if parentFound {
				break
			}
		}
		if !parentFound {
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
	if !GetAppManagerInstance().IsAppDevExist(appCode) {
		return "", fmt.Errorf("app with Code %s does not exist", appCode)
	}

	//  读取配置，返回字符串
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
