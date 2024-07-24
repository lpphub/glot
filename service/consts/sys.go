package consts

const (
	// StatusOn 通用状态 1-启用 2-禁用
	StatusOn  = 1
	StatusOff = 2

	ResourceDir    = 1 //菜单目录
	ResourceMenu   = 2 //菜单
	ResourceButton = 3 //按钮
)

var (
	MenuRoute = []int{ResourceDir, ResourceMenu}
)
