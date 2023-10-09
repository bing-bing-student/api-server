package service

import (
	"blog/global"
	"blog/model"
	"blog/utils"
	"errors"
)

func Register(username string, password string) (user *model.UserLogin, err error) {
	result := global.DB.Where("username = ?", username).Limit(1).Find(&user)
	if result.RowsAffected != 0 {
		err = errors.New("user already exists")
		return
	}

	user.Username = username                     //接收姓名
	user.Password = utils.BcryptHash(password)   //对明文密码加密
	user.UserID, _ = global.IdGenerator.NextID() //生成增长的 userID
	err = global.DB.Create(user).Error           //存储到数据库
	return
}

func Login(username string, password string) (user *model.UserLogin, err error) {
	//检查用户名是否存在
	result := global.DB.Where("username = ?", username).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		err = errors.New("username does not exist")
		return
	}
	//检查密码是否正确
	if ok := utils.BcryptCheck(password, user.Password); !ok {
		err = errors.New("wrong password")
		return
	}
	return
}
