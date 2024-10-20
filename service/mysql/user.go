package mysql

import (
	"blog/global"
	"blog/model"
	"blog/utils"
	"errors"
)

func IfExistPhoneNumber(phone string) bool {
	var user *model.UserLogin
	if affectNum := global.DB.Where("phone=?", phone).Find(&user).RowsAffected; affectNum == 1 {
		return true
	}
	return false
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

func ModifyAboutMe(userInfo *model.UserInfo) (result int64) {
	result = global.DB.Where("user_id=?", userInfo.UserID).Updates(&userInfo).RowsAffected
	return result
}

func GetManagerInfo() (managerInfo *model.UserInfo, result bool) {
	managerInfo = &model.UserInfo{}
	if err := global.DB.First(managerInfo).Error; err != nil {
		result = false
	}
	result = true
	return managerInfo, result
}

func GetUserId(phone *string) uint64 {
	var user *model.UserLogin
	global.DB.Where("phone=?", phone).Find(&user)
	return user.UserID
}
