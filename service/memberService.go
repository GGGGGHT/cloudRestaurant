package service

import (
	"cloudRestaurant/dao"
	"cloudRestaurant/model"
	"cloudRestaurant/param"
	"cloudRestaurant/tool"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"math/rand"
	"time"
)

type MemberService struct {
}

func (mc *MemberService) SendCode(phone string) bool {
	// 生产验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	config := tool.GetConfig()

	// 调用阿里云SDK 完成发送
	client, err := dysmsapi.NewClientWithAccessKey(config.Sms.RegionId, config.Sms.AppKey, config.Sms.AppSecret)

	if err != nil {
		fmt.Errorf("%v\n", err.Error())
		return false
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = config.Sms.SignName
	request.TemplateCode = "SMS_212702617"
	par, err := json.Marshal(map[string]interface{}{
		"code": code,
	})
	request.TemplateParam = string(par)

	response, err := client.SendSms(request)

	if err != nil {
		fmt.Print(err.Error())
		return false
	}
	fmt.Printf("response is %#v\n", response)
	if response.Code == "OK" {
		// 获取发送结果
		smsCode := model.SmsCode{
			Phone:      phone,
			BizId:      response.BizId,
			Code:       code,
			CreateTime: time.Now().UnixNano(),
		}

		memberDao := dao.MemberDao{
			Orm: tool.DbEngine,
		}

		return memberDao.InsertCode(smsCode) > 0
	}
	return false
}

func (mc *MemberService) SmsLogin(param *param.SmsLoginParam) *model.Member {
	// 获取到手机号和验证码
	// 验证手机号和验证码是否正确
	md := dao.MemberDao{Orm: tool.DbEngine}
	sms := md.ValidateSmsCode(param.Phone, param.Code)

	if sms.Id == 0 {
		return nil
	}

	// 根据手机号member表中查询记录
	member := md.GetMemberByPhone(param.Phone)
	// 如果不存在 则新创建一条member 并保存
	if nil == member || member.Id == 0 {
		member = &model.Member{
			UserName:     param.Phone,
			Mobile:       param.Phone,
			RegisterTime: time.Now().Unix(),
			IsActive:     0,
			City:         "",
		}

		member.Id = md.AddMember(member)
	}

	return member
}
