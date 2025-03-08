package service

import (
	"chat/model"
	"chat/serializer"
	"fmt"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	code := 200
	count := 0
	fmt.Println("-----", service.UserName)
	fmt.Println("-----", service.Password)
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count != 0 {
		code = 400
		return serializer.Response{
			Status: code,
			Msg:    "用户已存在",
		}
	}
	user = model.User{
		UserName: service.UserName,
	}
	if err := user.SetPassword(service.Password); err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "加密错误",
		}
	}
	model.DB.Create(&user)
	return serializer.Response{
		Status: code,
		Msg:    "创建成功",
	}
}
