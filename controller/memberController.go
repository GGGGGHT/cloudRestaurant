package controller

import (
	"cloudRestaurant/param"
	_ "cloudRestaurant/param"
	"cloudRestaurant/service"
	"cloudRestaurant/tool"
	"github.com/gin-gonic/gin"
)

type MemberController struct {
}

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendCode", mc.sendCode)
	engine.POST("/api/login_sms", mc.smsLogin)
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
		tool.Success(context,member)
		return
	}

	tool.Failed(context,"login failed")
}
