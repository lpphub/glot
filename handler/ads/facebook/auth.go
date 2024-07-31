package facebook

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"glot/handler/ads"
	"glot/middleware"
	"strings"
)

const (
	authLoginUrl = "https://www.facebook.com/v20.0/dialog/oauth?response_type=code&client_id=%s&config_id=%s&redirect_uri=%s&state=%s"
	getTokenUrl  = "https://graph.facebook.com/v20.0/oauth/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s"

	clientId      = ""
	clientSecret  = ""
	scopeConfigId = ""
	callbackFb    = "http://localhost:8080/ads/facebook/oauth2callback"
)

// GetAuthUrl 获取授权url
func GetAuthUrl(ctx *gin.Context) {
	state := fmt.Sprintf("%s_%s", lo.RandomString(32, lo.LettersCharset), cast.ToString(middleware.GetLoginTenantId(ctx)))
	render.JsonWithSuccess(ctx, fmt.Sprintf(authLoginUrl, clientId, scopeConfigId, callbackFb, state))
}

// AuthCallback 授权后的回调
func AuthCallback(ctx *gin.Context) {
	var cbParam ads.AuthCallbackParam
	_ = ctx.ShouldBindQuery(&cbParam)

	if cbParam.Error != "" {
		zlog.Errorf(ctx, "facebook callback error: %s", cbParam.ErrorReason)
		render.JsonWithSuccess(ctx, "ok")
		return
	}

	token, err := getAccessToken(ctx, cbParam.Code)
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

func getAccessToken(ctx *gin.Context, code string) (*ads.AuthToken, error) {
	resp, err := ads.GetHttpClient().R().Get(fmt.Sprintf(getTokenUrl, clientId, clientSecret, callbackFb, code))
	if err != nil {
		zlog.Errorf(ctx, "http_error: %s", err.Error())
		return nil, err
	}

	var token ads.AuthToken
	err = jsoniter.Unmarshal(resp.Body(), &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
