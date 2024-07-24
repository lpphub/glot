package errcode

import (
	"github.com/lpphub/golib/render"
)

var (
	ErrServerError = render.Error{
		Code: -1,
		Msg:  "server internal error",
	}
	ErrNotLogin = render.Error{
		Code: 1001,
		Msg:  "not login",
	}
	ErrInvalidToken = render.Error{
		Code: 1002,
		Msg:  "invalid token",
	}
	ErrToast = render.Error{
		Code: 1100,
		Msg:  "%s",
	}
	ErrParamInvalid = render.Error{
		Code: 1101,
		Msg:  "invalid param",
	}
	ErrApiFail = render.Error{
		Code: 1102,
		Msg:  "call api fail: %s",
	}

	ErrUserNotFound = render.Error{
		Code: 2001,
		Msg:  "用户不存在",
	}
)
