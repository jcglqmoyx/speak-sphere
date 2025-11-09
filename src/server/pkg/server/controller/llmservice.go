package controller

import (
	"log"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddLLMService(c *gin.Context) {
	var llmService model.LLMService
	if err := c.ShouldBind(&llmService); err != nil {
		util.JsonHttpResponse(c, 1, "参数错误", nil)
		return
	}

	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}
	llmService.UserID = userID

	// 如果设置为默认服务，先取消其他默认服务
	if llmService.IsDefault {
		existingDefault, found := dao.FindDefaultLLMServiceByUserID(userID)
		if found {
			dao.UpdateLLMServiceWithMap(map[string]interface{}{
				"id":         existingDefault.ID,
				"is_default": false,
			})
		}
	}

	dao.AddLLMService(&llmService)
	util.JsonHttpResponse(c, 0, "AI服务添加成功", llmService)
}

func GetLLMServices(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}

	services := dao.FindLLMServicesByUserID(userID)
	util.JsonHttpResponse(c, 0, "success", services)
}

func GetLLMServiceByID(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "ID格式错误", nil)
		return
	}

	service, found := dao.FindLLMServiceByID(id)
	if !found {
		util.JsonHttpResponse(c, 1, "AI服务不存在", nil)
		return
	}

	// 检查服务是否属于当前用户
	if service.UserID != userID {
		util.JsonHttpResponse(c, 1, "无权访问此服务", nil)
		return
	}

	util.JsonHttpResponse(c, 0, "success", service)
}

func UpdateLLMService(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "ID格式错误", nil)
		return
	}

	// 检查服务是否存在且属于当前用户
	existingService, found := dao.FindLLMServiceByID(id)
	if !found {
		util.JsonHttpResponse(c, 1, "AI服务不存在", nil)
		return
	}
	if existingService.UserID != userID {
		util.JsonHttpResponse(c, 1, "无权修改此服务", nil)
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法: "+err.Error(), nil)
		return
	}

	log.Printf("更新AI服务: userID=%d, serviceID=%d, data=%+v", userID, id, updateData)

	// 如果设置为默认服务，先取消其他默认服务
	if isDefault, ok := updateData["is_default"].(bool); ok && isDefault {
		existingDefault, found := dao.FindDefaultLLMServiceByUserID(userID)
		if found && existingDefault.ID != id {
			dao.UpdateLLMServiceWithMap(map[string]interface{}{
				"id":         existingDefault.ID,
				"is_default": false,
			})
		}
	}

	updateData["id"] = id
	dao.UpdateLLMServiceWithMap(updateData)

	updatedService, _ := dao.FindLLMServiceByID(id)
	util.JsonHttpResponse(c, 0, "AI服务更新成功", updatedService)
}

func DeleteLLMService(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "ID格式错误", nil)
		return
	}

	// 检查服务是否存在且属于当前用户
	service, found := dao.FindLLMServiceByID(id)
	if !found {
		util.JsonHttpResponse(c, 1, "AI服务不存在", nil)
		return
	}
	if service.UserID != userID {
		util.JsonHttpResponse(c, 1, "无权删除此服务", nil)
		return
	}

	dao.DeleteLLMServiceByID(id)
	util.JsonHttpResponse(c, 0, "AI服务删除成功", nil)
}

func GetDefaultLLMService(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}

	service, found := dao.FindDefaultLLMServiceByUserID(userID)
	if !found {
		util.JsonHttpResponse(c, 1, "未设置默认AI服务", nil)
		return
	}

	util.JsonHttpResponse(c, 0, "success", service)
}
