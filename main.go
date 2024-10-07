package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/ware"
	"glot/component/errcode"
	"glot/helper"
	"glot/router"
	"net/http"
)

func main() {
	app := gin.New()
	helper.InitResource()
	defer helper.Clear()

	ware.Bootstrap(app, ware.BootstrapConf{
		LogTrace: true,
		Cors:     true,
		CustomRecovery: func(ctx *gin.Context, err any) {
			render.JsonWithError(ctx, errcode.ErrServerError)
		},
	})

	router.Handle(app)

	ware.ListenAndServe(&http.Server{
		Addr:    ":8080",
		Handler: app,
	})
}
