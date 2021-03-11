package controller

import (
	"cloudRestaurant/model"
	"cloudRestaurant/param"
	_ "cloudRestaurant/param"
	"cloudRestaurant/service"
	"cloudRestaurant/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mojocn/base64Captcha"
	"strconv"
	"time"
)

const UPLOADFILEPATH = "/uploadfile/"

type MemberController struct {
}

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendCode", mc.sendCode)
	engine.POST("/api/login_sms", mc.smsLogin)
	engine.GET("/api/captcha", mc.captcha)
	engine.GET("/api/verifycha", mc.verify)
	engine.POST("/api/login_pwd", mc.pwdLogin)
	// 上传头像
	engine.POST("/api/upload/avatar", mc.uploadAvatar)
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
		json, _ := json.Marshal(member)
		if err := tool.SetSession(context, "user_"+strconv.FormatInt(member.Id, 10), json); err != nil {
			tool.Failed(context, "保存session失败")
			return
		}
		tool.Success(context, member)
		return
	}

	tool.Failed(context, "login failed")
}

// 生成验证码
func (mc *MemberController) captcha(context *gin.Context) {
	tool.GenerateCaptcha(context)
}

func (mc *MemberController) verify(context *gin.Context) {
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

func (mc *MemberController) pwdLogin(context *gin.Context) {
	// 解析参数
	var loginParam param.LoginParam
	if err := tool.Decoder(context.Request.Body, &loginParam); err != nil {
		tool.Failed(context, "parsed failed")
		return
	}

	// 验证验证码
	if !tool.Rstore.Verify(loginParam.Id, loginParam.Value, false) {
		tool.Failed(context, "验证码不正确,请重新输入!")
		return
	}
	// 登录
	ms := service.MemberService{}
	member := ms.PwdLogin(loginParam.Name, loginParam.Password)
	if member == nil || member.Id == 0 {
		tool.Failed(context, "用户名或密码不正确,请重试试")
		return
	}

	json, _ := json.Marshal(member)
	if err := tool.SetSession(context, "user_"+strconv.FormatInt(member.Id, 10), json); err != nil {
		tool.Failed(context, "保存session失败")
		return
	}

	tool.Success(context, member)
}

func (mc *MemberController) uploadAvatar(context *gin.Context) {
	// 解析参数 file文件 用户ID
	userId := context.PostForm("user_id")
	fmt.Println(userId)
	file, err := context.FormFile("avatar")
	if err != nil {
		tool.Failed(context, "文件解析失败")
		return
	}
	// 判断用户ID是否已经登录
	session := tool.GetSession(context, "user_"+userId)
	if session == nil {
		tool.Failed(context, "参数不合法")
		return
	}

	var member *model.Member
	if err = json.Unmarshal(session.([]byte), &member); err != nil {
		tool.Failed(context, "参数不合法")
		return
	}
	// 保存到本地(文件服务器上)
	fileName := UPLOADFILEPATH + strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	if err = context.SaveUploadedFile(file, fileName); err != nil {
		tool.Failed(context, "更新文件失败")
		return
	}
	// 将文件路径 保存到用户表中头像字段
	ms := service.MemberService{}
	path := ms.UploadAvatar(member.Id, fileName)
	if path != "" {
		tool.Success(context, "http://localhost:8080"+path)
		return
	}
}
