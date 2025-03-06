package app

type MenuType string

const (
	Group           MenuType = "group"
	MenuImperative  MenuType = "imperative"
	MenuTypeLowCode MenuType = "lowCode"
	MenuTypeLink    MenuType = "link"
)

type App struct {
	AppId       string `json:"appId"`
	AppCode     string `json:"appCode"`
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type Menu struct {
	MenuCode  string `json:"menuCode"`
	MenuLabel string `json:"menuLabel"`
	MenuPath  string `json:"menuPath"`
	MenuIcon  string `json:"menuIcon"`
	Children  []Menu `json:"children"`
	Hide      bool   `json:"hide"`
	MenuType  string `json:"menuType"`
	LinkUrl   string `json:"linkUrl"`
}
