package ads

import "github.com/go-resty/resty/v2"

type AuthCallbackParam struct {
	Code        string `form:"code"`
	State       string `form:"state"`
	Error       string `form:"error"`
	ErrorReason string `form:"error_reason"`
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

var (
	httpClient *resty.Client
)

func GetHttpClient() *resty.Client {
	if httpClient == nil {
		httpClient = resty.New()
	}
	return httpClient
}
