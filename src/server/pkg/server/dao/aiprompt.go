package dao

import (
	"speak-sphere/pkg/server/model"
)

func AddAIPrompt(aiPrompt *model.AIPrompt) *model.AIPrompt {
	db := GetDB()

	// 检查是否需要设置默认提示词
	if aiPrompt.IsDefault && aiPrompt.UserID != 0 { // 只有用户自定义的提示词可以设置默认
		// 将该用户的所有其他提示词设置为非默认（不包括系统默认）
		db.Model(&model.AIPrompt{}).
			Where("user_id = ?", aiPrompt.UserID).
			Update("is_default", false)
	}

	db.Create(&aiPrompt)
	return aiPrompt
}

func DeleteAIPromptByID(id int) {
	db := GetDB()
	db.Delete(&model.AIPrompt{}, id)
}

func UpdateAIPromptWithMap(updateData map[string]interface{}) {
	db := GetDB()

	// 先保存ID值
	id := updateData["id"]

	// 检查是否需要设置默认提示词
	if isDefault, exists := updateData["is_default"]; exists && isDefault == true {
		// 先找到要更新的提示词，获取其用户ID
		prompt, found := FindAIPromptByID(id.(int))
		if found && prompt.UserID != 0 { // 只有用户自定义的提示词可以设置默认
			// 将该用户的所有其他提示词设置为非默认（不包括系统默认）
			db.Model(&model.AIPrompt{}).
				Where("user_id = ? AND id != ?", prompt.UserID, id).
				Update("is_default", false)
		}
	}

	// 移除不应通过API更新的字段
	delete(updateData, "id")
	delete(updateData, "user_id")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")

	db.Model(&model.AIPrompt{}).Where("id = ?", id).Updates(updateData)
}

func FindAIPromptsByUserID(userID int) []*model.AIPrompt {
	db := GetDB()
	var prompts []*model.AIPrompt
	db.Where("user_id = ? OR user_id = 0", userID).Find(&prompts) // 包含系统默认和用户自定义
	return prompts
}

func FindAIPromptByID(id int) (*model.AIPrompt, bool) {
	db := GetDB()
	var prompts []*model.AIPrompt
	db.Limit(1).Where("id = ?", id).Find(&prompts)
	if len(prompts) > 0 {
		return prompts[0], true
	}
	return nil, false
}

func FindDefaultAIPromptByUserID(userID int) (*model.AIPrompt, bool) {
	db := GetDB()
	var prompts []*model.AIPrompt
	db.Limit(1).Where("(user_id = ? OR user_id = 0) AND is_default = ?", userID, true).Find(&prompts)
	if len(prompts) > 0 {
		return prompts[0], true
	}
	return nil, false
}

func CountAIPromptsByUserID(userID int) int64 {
	var count int64
	db := GetDB()
	db.Model(&model.AIPrompt{}).Where("user_id = ?", userID).Count(&count)
	return count
}
