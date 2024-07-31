package google

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
	authLoginUrl = "https://accounts.google.com/o/oauth2/auth?response_type=code&access_type=offline&include_granted_scopes=true&client_id=%s&scope=%s&redirect_uri=%s&state=%s"
	getTokenUrl  = "https://oauth2.googleapis.com/token?grant_type=authorization_code&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s"

	clientId     = ""
	clientSecret = ""
	callbackGg   = "http://localhost:8080/ads/google/oauth2callback"
	//developerToken = ""
)

var (
	defaultScopes = []string{"https://www.googleapis.com/auth/adwords"}
)

// GetAuthUrl 获取授权登录url
func GetAuthUrl(ctx *gin.Context) {
	var (
		state  = fmt.Sprintf("%s_%s", lo.RandomString(32, lo.LettersCharset), cast.ToString(middleware.GetLoginTenantId(ctx)))
		scopes = strings.Join(defaultScopes, " ")
	)
	render.JsonWithSuccess(ctx, fmt.Sprintf(authLoginUrl, clientId, scopes, callbackGg, state))
}

// AuthCallback 授权后的回调
func AuthCallback(ctx *gin.Context) {
	var cbParam ads.AuthCallbackParam
	_ = ctx.ShouldBindQuery(&cbParam)

	if cbParam.Error != "" {
		zlog.Errorf(ctx, "google callback error: %s", cbParam.Error)
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
	zlog.Infof(ctx, "google - tenantId: %s, accessToken: %s, expiresIn: %d", tenantId, token.AccessToken, token.ExpiresIn)
	// todo 存储 access_token

	render.JsonWithSuccess(ctx, "ok")
}

func getAccessToken(ctx *gin.Context, code string) (*ads.AuthToken, error) {
	resp, err := ads.GetHttpClient().R().Get(fmt.Sprintf(getTokenUrl, clientId, clientSecret, callbackGg, code))
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
