package common

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/util/captcha"
)

// CaptchaResult 存储验证码的结构
type CaptchaResult struct {
	Id         string `json:"id"`
	Base64Blob string `json:"base64Blob"`
	//VerifyValue string `json:"code"`
}

type CaptchaHandler struct{}

// GetCaptcha 生成图形化验证码
func (a *CaptchaHandler) GetCaptcha(ctx *gin.Context) {
	id, b64s, err := captcha.CaptMake()
	if err != nil {
		result.Fail(ctx)
	}
	captchaResult := CaptchaResult{
		Id:         id,
		Base64Blob: b64s,
	}
	result.OkWithData(captchaResult, ctx)
}

// VerifyCaptcha 验证captcha是否正确
func (a *CaptchaHandler) VerifyCaptcha(ctx *gin.Context) {

	id := ctx.PostForm("id")
	capt := ctx.PostForm("capt")
	if id == "" || capt == "" {
		result.Fail(ctx)
	}
	if captcha.CaptVerify(id, capt) == true {
		result.Ok(ctx)
	} else {
		result.Fail(ctx)
	}
}
