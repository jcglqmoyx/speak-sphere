package dao

import (
	"speak-sphere/pkg/server/model"
)

func AddLLMService(llmService *model.LLMService) *model.LLMService {
	db := GetDB()
	db.Create(&llmService)
	return llmService
}

func DeleteLLMServiceByID(id int) {
	db := GetDB()
	db.Delete(&model.LLMService{}, id)
}

func UpdateLLMServiceWithMap(updateData map[string]interface{}) {
	db := GetDB()
	// 移除不应通过API更新的字段
	delete(updateData, "user_id")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")

	db.Model(&model.LLMService{}).Where("id = ?", updateData["id"]).Updates(updateData)
}

func FindLLMServicesByUserID(userID int) []*model.LLMService {
	db := GetDB()
	var services []*model.LLMService
	db.Where("user_id = ?", userID).Find(&services)
	return services
}

func FindLLMServiceByID(id int) (*model.LLMService, bool) {
	db := GetDB()
	var services []*model.LLMService
	db.Limit(1).Where("id = ?", id).Find(&services)
	if len(services) > 0 {
		return services[0], true
	}
	return nil, false
}

func FindDefaultLLMServiceByUserID(userID int) (*model.LLMService, bool) {
	db := GetDB()
	var services []*model.LLMService
	db.Limit(1).Where("user_id = ? AND is_default = ?", userID, true).Find(&services)
	if len(services) > 0 {
		return services[0], true
	}
	return nil, false
}

func CountLLMServicesByUserID(userID int) int64 {
	var count int64
	db := GetDB()
	db.Model(&model.LLMService{}).Where("user_id = ?", userID).Count(&count)
	return count
}
