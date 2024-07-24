package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
)

func Test(ctx *gin.Context) {
	render.JsonWithSuccess(ctx, "test")
}
