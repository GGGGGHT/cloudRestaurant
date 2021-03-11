package dao

import (
	"cloudRestaurant/model"
	"cloudRestaurant/tool"
	"fmt"
)

type MemberDao struct {
	*tool.Orm
}

func (md *MemberDao) InsertCode(sms model.SmsCode) int64 {
	res, err := md.InsertOne(&sms)
	if err != nil {
		//logger.Errorf(err.)
		fmt.Errorf("%v", err.Error())
	}
	return res
}

func (md *MemberDao) ValidateSmsCode(phone, code string) *model.SmsCode {
	var sms model.SmsCode
	if _, err := md.Where("phone = ? and code = ?", phone, code).Get(&sms); err != nil {
		fmt.Println(err.Error())
	}
	return &sms
}

func (md *MemberDao) GetMemberByPhone(phone string) *model.Member {
	var member model.Member

	if _, err := md.Where("Mobile = ?", phone).Get(&member); err != nil {
		fmt.Println(err.Error())
	}

	return &member
}

func (md *MemberDao) AddMember(member *model.Member) int64 {
	res, err := md.InsertOne(member)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return res
}
