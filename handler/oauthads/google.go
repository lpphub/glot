package oauthads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/zlog"
	"github.com/samber/lo"
	"glot/middleware"
	"strings"
)

const (
	googleOAuthUrl    = "https://accounts.google.com/o/oauth2/auth?response_type=code&access_type=offline&include_granted_scopes=true&client_id=%s&scope=%s&redirect_uri=%s&state=%s"
	googleGetTokenUrl = "https://oauth2.googleapis.com/token?grant_type=authorization_code&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s"

	googleClientId     = ""
	googleClientSecret = ""
	googleCallbackUrl  = "http://localhost:8080/ads/google/oauth2callback"
	//developerToken = ""
)

var (
	defaultScopes = []string{"https://www.googleapis.com/auth/adwords"}
)

// 获取授权登录url
func getGoogleOAuthUrl(ctx *gin.Context) string {
	var (
		randStr  = lo.RandomString(32, lo.LettersCharset)
		tenantId = middleware.GetLoginTenantId(ctx)

		state  = fmt.Sprintf("%s_%d", randStr, tenantId)
		scopes = strings.Join(defaultScopes, " ")
	)
	return fmt.Sprintf(googleOAuthUrl, googleClientId, scopes, googleCallbackUrl, state)
}

func getGoogleAccessToken(ctx *gin.Context, code string) (*OAuthToken, error) {
	resp, err := getHttpClient().R().Get(fmt.Sprintf(googleGetTokenUrl, googleClientId, googleClientSecret, googleCallbackUrl, code))
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
