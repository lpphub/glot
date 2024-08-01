package router

import (
	"github.com/gin-gonic/gin"
	"glot/handler"
	"glot/handler/oauthads"
	"glot/handler/system"
	"glot/handler/tenant"
	"glot/middleware"
)

func Handle(r *gin.Engine) {
	r.GET("/test", handler.Test)

	r.POST("/auth/login", handler.Login)
	r.GET("/auth/get_user", middleware.CheckAuthLogin, handler.LoginByToken)

	sys := r.Group("/system", middleware.CheckAuthLogin)
	{
		sys.GET("/get_user_routes", system.GetUserRoutes)
		sys.GET("/is_exist_route", system.IsExistRoute)

		sys.GET("/user/list", system.PageListUser)
		sys.POST("/user/post", system.SaveUser)
		sys.POST("/user/del", system.DelUser)

		sys.GET("/role/list", system.PageListRole)
		sys.GET("/role/all", system.ListAllRole)
		sys.POST("/role/post", system.SaveRole)
		sys.POST("/role/del", system.DelRole)
		sys.GET("/role/menu", system.GetRoleMenu)
		sys.POST("/role/bind_menu", system.BindRoleMenu)

		sys.POST("/menu/post", system.SaveMenu)
		sys.GET("/menu/list", system.PageListMenu)
		sys.GET("/menu/tree", system.GetMenuTree)
		sys.GET("/menu/button", system.GetMenuButton)
		sys.POST("/menu/del", system.DelMenu)
	}

	tnt := r.Group("/tenant", middleware.CheckAuthLogin)
	{
		tnt.GET("/list", tenant.PageListTenant)
		tnt.POST("/post", tenant.SaveTenant)
		tnt.GET("/role_scope", tenant.ListRoleScope)
	}

	oauth := r.Group("/oauth")
	{
		oauth.GET("/get_auth_url", middleware.CheckAuthLogin, oauthads.GetOAuthUrl)

		oauth.GET("/facebook/oauth2callback", oauthads.FacebookOAuthCallback)
		oauth.GET("/google/oauth2callback", oauthads.GoogleOAuthCallback)
	}

}
