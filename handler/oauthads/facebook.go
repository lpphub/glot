package oauthads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/zlog"
	"github.com/samber/lo"
	"glot/middleware"
)

const (
	facebookOAuthUrl    = "https://www.facebook.com/v20.0/dialog/oauth?response_type=code&client_id=%s&config_id=%s&redirect_uri=%s&state=%s"
	facebookGetTokenUrl = "https://graph.facebook.com/v20.0/oauth/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s"

	facebookClientId      = ""
	facebookClientSecret  = ""
	facebookScopeConfigId = ""
	facebookCallbackUrl   = "http://localhost:8080/ads/facebook/oauth2callback"
)

// 获取授权url
func getFacebookAuthUrl(ctx *gin.Context) string {
	var (
		randStr  = lo.RandomString(32, lo.LettersCharset)
		tenantId = middleware.GetLoginTenantId(ctx)

		state = fmt.Sprintf("%s_%d", randStr, tenantId)
	)
	return fmt.Sprintf(facebookOAuthUrl, facebookClientId, facebookScopeConfigId, facebookCallbackUrl, state)
}

func getFacebookAccessToken(ctx *gin.Context, code string) (*OAuthToken, error) {
	resp, err := getHttpClient().R().Get(fmt.Sprintf(facebookGetTokenUrl, facebookClientId, facebookClientSecret, facebookCallbackUrl, code))
	if err != nil {
		zlog.Errorf(ctx, "http_error: %s", err.Error())
		return nil, err
	}

	var token OAuthToken
	err = jsoniter.Unmarshal(resp.Body(), &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
