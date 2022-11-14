package common

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/service/syssrv"
	"go-vea/util/captcha"
)

// CaptchaResult 存储验证码的结构
type CaptchaResult struct {
	Id             string `json:"id"`
	Base64Blob     string `json:"base64Blob"`
	CaptchaEnabled bool   `json:"captchaEnabled"`
	//VerifyValue string `json:"code"`
}

type CaptchaHandler struct{}

// GetCaptcha 生成图形化验证码
func (*CaptchaHandler) GetCaptcha(ctx *gin.Context) {
	id, b64s, err := captcha.CaptMake()
	if err != nil {
		result.Fail(ctx)
		return
	}
	captchaResult := CaptchaResult{}
	captchaEnabled, _ := syssrv.SysConfigSrv.SelectCaptchaEnabled(ctx)
	if !captchaEnabled {
		captchaResult.CaptchaEnabled = false
		result.OkWithData(captchaResult, ctx)
	} else {
		captchaResult.CaptchaEnabled = true
		captchaResult.Id = id
		captchaResult.Base64Blob = b64s
		result.OkWithData(captchaResult, ctx)
	}
}

// VerifyCaptcha 验证captcha是否正确
func (*CaptchaHandler) VerifyCaptcha(ctx *gin.Context) {

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
