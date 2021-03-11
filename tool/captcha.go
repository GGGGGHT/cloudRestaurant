package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type CaptchaResult struct {
	Id          string `json:"id"`
	Base64Blob  string `json:"base_64_blob"`
	VerifyValue string `json:"code"`
}

var Captcha *base64Captcha.Captcha

// 生成图形化验证码
func GenerateCaptcha(ctx *gin.Context) {
	driverString := base64Captcha.DriverString{
		Height:          60,
		Width:           240,
		NoiseCount:      0,
		ShowLineOptions: 0,
		Length:          6,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 254,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driver := driverString.ConvertFonts()

	captcha := base64Captcha.NewCaptcha(driver, Rstore)
	id, b64s, err := captcha.Generate()
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
	}
	fmt.Printf("id: %v\t b64s: %v\n", id, b64s)

	Success(ctx, map[string]interface{}{
		"code":      1,
		"data":      b64s,
		"captchaId": id,
		"msg":       "success",
	})
}
