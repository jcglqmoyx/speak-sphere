package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
)

func Register(c *gin.Context) {
	type RegisterForm struct {
		Username        string `form:"username" json:"username"`
		Email           string `form:"email" json:"email"`
		Avatar          string `form:"avatar" json:"avatar"`
		Password        string `form:"password" json:"password"`
		ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
	}
	var registerForm RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		util.JsonHttpResponse(c, 1, "参数错误", nil)
		return
	}
	if registerForm.Username == "" || registerForm.Password == "" || registerForm.ConfirmPassword == "" {
		util.JsonHttpResponse(c, 1, "注册失败, 用户名或密码不能为空", nil)
		return
	}
	if registerForm.ConfirmPassword != registerForm.Password {
		util.JsonHttpResponse(c, 1, "注册失败, 两次输入的密码不一致", nil)
		return
	}
	if _, found := dao.FindUserByUsername(registerForm.Username); found {
		util.JsonHttpResponse(c, 1, "该用户名已被占用", nil)
		return
	}
	if _, found := dao.FindUserByEmail(registerForm.Email); found {
		util.JsonHttpResponse(c, 1, "该邮箱已被占用", nil)
		return
	}
	user := model.User{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Avatar:   registerForm.Avatar,
		Password: registerForm.Password,
	}
	defaultDictionaries := conf.Cfg.Dictionary.Dictionaries
	dao.AddUser(&user)
	for _, dictionary := range defaultDictionaries {
		dao.AddDictionary(&model.Dictionary{
			UserID: user.ID,
			Title:  dictionary.Title,
			Prefix: dictionary.Prefix,
			Suffix: dictionary.Suffix,
		})
	}
	util.JsonHttpResponse(c, 0, "注册成功", nil)
}

func Login(c *gin.Context) {
	var dto *model.User
	if err := c.ShouldBind(&dto); err != nil {
		util.JsonHttpResponse(c, 1, "参数错误", nil)
		return
	}
	user, found := dao.FindUserByUsername(dto.Username)
	if !found {
		util.JsonHttpResponse(c, 1, "用户名不存在", nil)
		return
	}
	if util.HashPassword(dto.Password, user.PasswordSalt) != user.Password {
		util.JsonHttpResponse(c, 1, "密码错误", nil)
		return
	}
	token, _ := util.GenerateJWT(user.ID)
	json := map[string]string{
		"token":   token,
		"user_id": strconv.Itoa(user.ID),
	}
	util.JsonHttpResponse(c, 0, "登录成功", json)
}

func DeleteUserByID(c *gin.Context) {
	if userID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if _, found := dao.FindUserByID(userID); !found {
			util.JsonHttpResponse(c, 1, "该用户不存在, 无法删除", nil)
		} else {
			dao.DeleteUserByID(userID)
			util.JsonHttpResponse(c, 0, "用户删除成功", nil)
		}
	}
}

func UpdateUser(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	} 
	
	currentUser, found := dao.FindUserByID(userID)
	if !found {
		util.JsonHttpResponse(c, 1, "用户不存在, 无法更新", nil)
		return
	}
	
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法: " + err.Error(), nil)
		return
	}
	
	// 记录接收到的数据
	log.Printf("收到用户更新请求: userID=%d, data=%+v", userID, updateData)
	
	// 检查邮箱是否被其他用户占用
	if email, ok := updateData["email"].(string); ok && email != "" && email != currentUser.Email {
		if existingUser, found := dao.FindUserByEmail(email); found && existingUser.ID != userID {
			util.JsonHttpResponse(c, 1, "该邮箱已被其他用户占用", nil)
			return
		}
	}

	// 检查词书ID
	if currentBookID, ok := updateData["current_book_id"].(float64); ok {
		bookID := int(currentBookID)
		if bookID != 0 {
			book, found := dao.FindBookByID(bookID)
			if !found {
				util.JsonHttpResponse(c, 1, "词书不存在", nil)
				return
			}
			if book.UserID != userID {
				util.JsonHttpResponse(c, 1, "词书不属于该用户", nil)
				return
			}
		}
	}

	// 验证参数
	if dailyCount, ok := updateData["daily_count"].(float64); ok {
		dailyCountInt := int(dailyCount)
		if dailyCountInt <= 0 {
			util.JsonHttpResponse(c, 1, "每日学习量不合法", nil)
			return
		}
		if dailyCountInt > 500000 {
			util.JsonHttpResponse(c, 1, "每日学习量过大", nil)
			return
		}
	}

	if reviewFrequencyFormula, ok := updateData["review_frequency_formula"].(string); ok {
		if reviewFrequencyFormula == "" {
			util.JsonHttpResponse(c, 1, "复习频率公式不能为空", nil)
			return
		}
		_, ok := util.ParseReviewFrequencyFormula(reviewFrequencyFormula)
		if !ok {
			util.JsonHttpResponse(c, 1, "复习频率公式不合法", nil)
			return
		}
	}

	if timesCountedAsKnown, ok := updateData["times_counted_as_known"].(float64); ok {
		timesCountedAsKnownInt := int(timesCountedAsKnown)
		if timesCountedAsKnownInt <= 0 {
			util.JsonHttpResponse(c, 1, "\"点击几次“认识”算学会\"参数不合法", nil)
			return
		}
		if timesCountedAsKnownInt > 1000 {
			util.JsonHttpResponse(c, 1, "\"点击几次“认识”算学会\"参数不合法", nil)
			return
		}
	}

	// 只更新需要的字段，避免清空密码等敏感信息
	updateData["id"] = userID
	dao.UpdateUserWithMap(updateData)
	
	// 获取更新后的用户信息
	updatedUser, _ := dao.FindUserByID(userID)
	util.JsonHttpResponse(c, 0, "success", updatedUser)
}

func GetUserProfile(c *gin.Context) {
	id, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, err.Error(), nil)
		return
	}
	user, found := dao.FindUserByID(id)
	if found {
		userInfo := model.User{
			Username:               user.Username,
			Email:                  user.Email,
			Avatar:                 user.Avatar,
			CurrentBookID:          user.CurrentBookID,
			DailyCount:             user.DailyCount,
			TimesCountedAsKnown:    user.TimesCountedAsKnown,
			ReviewFrequencyFormula: user.ReviewFrequencyFormula,
		}
		util.JsonHttpResponse(c, 0, "success", userInfo)
	} else {
		util.JsonHttpResponse(c, 1, "用户不存在", nil)
	}
}
