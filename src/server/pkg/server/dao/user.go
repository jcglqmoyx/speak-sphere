package dao

import (
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
)

func AddUser(user *model.User) *model.User {
	userCount := CountUser()
	if userCount == 0 {
		user.Level = 0
	} else {
		user.Level = 1
	}

	user.PasswordSalt = util.GenerateSalt()
	user.Password = util.HashPassword(user.Password, user.PasswordSalt)

	db := GetDB()
	db.Create(&user)
	return user
}

func DeleteUserByID(id int) {
	db := GetDB()
	db.Delete(&model.User{}, id)
}

func UpdateUserWithMap(updateData map[string]interface{}) {
	db := GetDB()
	// 移除不应通过API更新的字段
	delete(updateData, "password")
	delete(updateData, "password_salt")
	delete(updateData, "level")

	db.Model(&model.User{}).Where("id = ?", updateData["id"]).Updates(updateData)
}

func CountUser() int64 {
	var userCount int64
	db := GetDB()
	if err := db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return 0
	}
	return userCount
}

func FindUserByUsername(username string) (*model.User, bool) {
	db := GetDB()
	var users []*model.User
	db.Limit(1).Where("username = ?", username).Find(&users)
	if len(users) > 0 {
		return users[0], true
	}
	return nil, false
}

func FindUserByEmail(email string) (*model.User, bool) {
	db := GetDB()
	var users []*model.User
	db.Limit(1).Where("email= ?", email).Find(&users)
	if len(users) > 0 {
		return users[0], true
	}
	return nil, false
}

func FindUserByID(id int) (*model.User, bool) {
	db := GetDB()
	var users []*model.User
	db.Limit(1).Where("id= ?", id).Find(&users)
	if len(users) > 0 {
		return users[0], true
	}
	return nil, false
}
