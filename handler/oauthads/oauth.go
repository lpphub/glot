package oauthads

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"strings"
)

type OAuthCallbackParam struct {
	Code        string `form:"code"`
	State       string `form:"state"`
	Error       string `form:"error"`
	ErrorReason string `form:"error_reason"`
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

var (
	httpClient *resty.Client
)

func getHttpClient() *resty.Client {
	if httpClient == nil {
		httpClient = resty.New()
	}
	return httpClient
}

func GetOAuthUrl(ctx *gin.Context) {
	plat := ctx.Query("plat")
	if plat == "" {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	var url string
	switch plat {
	case "facebook":
		url = getFacebookAuthUrl(ctx)
	case "google":
		url = getGoogleOAuthUrl(ctx)
	}
	if url == "" {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	render.JsonWithSuccess(ctx, url)
}

func FacebookOAuthCallback(ctx *gin.Context) {
	var cbParam OAuthCallbackParam
	_ = ctx.ShouldBindQuery(&cbParam)

	if cbParam.Error != "" {
		zlog.Errorf(ctx, "facebook callback error: %s", cbParam.ErrorReason)
		render.JsonWithSuccess(ctx, "ok")
		return
	}

	token, err := getFacebookAccessToken(ctx, cbParam.Code)
	if err != nil {
		zlog.Errorf(ctx, err.Error())
		render.JsonWithFail(ctx, 3003, "")
		return
	}
	tenantId := strings.Split(cbParam.State, "_")[1]
	zlog.Infof(ctx, "fb - tenantId: %s, accessToken: %s, expiresIn: %d", tenantId, token.AccessToken, token.ExpiresIn)
	// todo 存储 access_token

	render.JsonWithSuccess(ctx, "ok")
}

func GoogleOAuthCallback(ctx *gin.Context) {
	var cbParam OAuthCallbackParam
	_ = ctx.ShouldBindQuery(&cbParam)

	if cbParam.Error != "" {
		zlog.Errorf(ctx, "google callback error: %s", cbParam.Error)
		render.JsonWithSuccess(ctx, "ok")
		return
	}

	token, err := getGoogleAccessToken(ctx, cbParam.Code)
	if err != nil {
		zlog.Errorf(ctx, err.Error())
		render.JsonWithFail(ctx, 3003, "")
		return
	}
	tenantId := strings.Split(cbParam.State, "_")[1]
	zlog.Infof(ctx, "google - tenantId: %s, accessToken: %s, expiresIn: %d", tenantId, token.AccessToken, token.ExpiresIn)
	// todo 存储 access_token

	render.JsonWithSuccess(ctx, "ok")
}
