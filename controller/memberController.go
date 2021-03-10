package controller

import (
	"cloudRestaurant/param"
	_ "cloudRestaurant/param"
	"cloudRestaurant/service"
	"cloudRestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mojocn/base64Captcha"
)

type MemberController struct {
}

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendCode", mc.sendCode)
	engine.POST("/api/login_sms", mc.smsLogin)
	engine.GET("/api/captcha", mc.captcha)
	engine.GET("/api/verifycha", mc.vertifyCaptcha)
}

func (mc *MemberController) sendCode(context *gin.Context) {
	// 发送验证码
	phone, exist := context.GetQuery("phone")
	if !exist {
		tool.Failed(context, "非法参数")
		return
	}

	ms := service.MemberService{}
	if ms.SendCode(phone) {
		tool.Success(context, "send Success")
		return
	}

	tool.Failed(context, "send error")
}

func (mc *MemberController) smsLogin(context *gin.Context) {
	var par param.SmsLoginParam
	if err := tool.Decoder(context.Request.Body, &par); err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}

	// 完成手机+验证码登录
	ms := service.MemberService{}
	member := ms.SmsLogin(&par)
	if member != nil {
		tool.Success(context, member)
		return
	}

	tool.Failed(context, "login failed")
}

// 生成验证码
func (mc *MemberController) captcha(context *gin.Context) {
	tool.GenerateCaptcha(context)
}

func (mc *MemberController) vertifyCaptcha(context *gin.Context) {
	var captcha tool.CaptchaResult
	err := tool.Decoder(context.Request.Body, &captcha)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}

	if tool.Rstore.Verify(captcha.Id, captcha.VerifyValue, false) {
		fmt.Printf("验证通过!\n")
	} else {
		fmt.Printf("验证失败!\n")
	}

}
