package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/lpphub/golib/zlog"
	"github.com/tidwall/gjson"
	"glot/component/errcode"
	"glot/helper"
	"time"
)

type WxApi struct {
	domain string
	client *resty.Client
}

var Wx = &WxApi{
	domain: "https://api.weixin.qq.com",
	client: resty.New(),
}

const (
	_wxAppid  = "wxf771e99449615cb2"
	_wxSecret = "19288d1e967e8370d3c0875e320cf900"

	_wxTokenKey = "wx_access_token"
)

func (wx *WxApi) GetAccessToken(ctx *gin.Context) (string, error) {
	if token, found := helper.LocalCache.Get(_wxTokenKey); found {
		return token.(string), nil
	}

	params := map[string]string{
		"appid":      _wxAppid,
		"secret":     _wxSecret,
		"grant_type": "client_credential",
	}
	resp, err := wx.client.R().SetQueryParams(params).Get(wx.domain + "/cgi-bin/token")
	if err != nil {
		zlog.Errorf(ctx, "http_error: %s", err.Error())
		return "", err
	}
	var (
		token     = gjson.GetBytes(resp.Body(), "access_token").String()
		expiresIn = gjson.GetBytes(resp.Body(), "expires_in").Int()
	)
	// 缓存token todo redis替换
	helper.LocalCache.Set(_wxTokenKey, token, time.Duration(expiresIn-10)*time.Second)
	return token, nil
}

// GetWxACodeUnLimit 获取小程序码
func (wx *WxApi) GetWxACodeUnLimit(ctx *gin.Context, scene string) (string, error) {
	// todo 如果用微信云托管则可去掉token
	token, err := wx.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"scene": scene,
		"page":  "pages/index/index",
	}
	url := fmt.Sprintf("%s?access_token=%s", wx.domain+"/wxa/getwxacodeunlimit", token)
	resp, err := wx.client.R().SetBody(params).Post(url)
	if err != nil {
		zlog.Errorf(ctx, "http_error: %s", err.Error())
		return "", err
	}
	if errno := gjson.GetBytes(resp.Body(), "errcode").Int(); errno != 0 {
		errMsg := gjson.GetBytes(resp.Body(), "errmsg").String()
		return "", errcode.ErrApiFail.Sprintf(fmt.Sprintf("%d-%s", errno, errMsg))
	}

	// 小程序码图片buf
	codeBuf := gjson.GetBytes(resp.Body(), "buffer").Value()
	if buf, ok := codeBuf.([]byte); ok {
		return base64.StdEncoding.EncodeToString(buf), nil
	}
	return "", nil
}

func (wx *WxApi) GetSession(ctx *gin.Context, code string) (string, string, error) {
	params := map[string]string{
		"appid":      _wxAppid,
		"secret":     _wxSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	resp, err := wx.client.R().SetQueryParams(params).Get(wx.domain + "/sns/jscode2session")
	if err != nil {
		zlog.Errorf(ctx, "http_error: %s", err.Error())
		return "", "", err
	}
	if errno := gjson.GetBytes(resp.Body(), "errcode").Int(); errno != 0 {
		errMsg := gjson.GetBytes(resp.Body(), "errmsg").String()
		return "", "", errcode.ErrApiFail.Sprintf(fmt.Sprintf("%d-%s", errno, errMsg))
	}
	var (
		sessionKey = gjson.GetBytes(resp.Body(), "session_key").String()
		openId     = gjson.GetBytes(resp.Body(), "openid").String()
	)
	return sessionKey, openId, nil
}
