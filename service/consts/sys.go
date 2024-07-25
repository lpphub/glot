package consts

const (
	// StatusOn 通用状态 1-启用 2-禁用
	StatusOn  = 1
	StatusOff = 2

	MenuDir    = 1 //目录
	MenuOpt    = 2 //菜单
	MenuButton = 3 //按钮
)

var (
	RouteMenu = []int{MenuDir, MenuOpt}
)
