package controller

import (
    "log"
    "speak-sphere/pkg/server/dao"
    "speak-sphere/pkg/server/model"
    "speak-sphere/pkg/server/util"
    "strconv"

    "github.com/gin-gonic/gin"
)

func AddAIPrompt(c *gin.Context) {
    var aiPrompt model.AIPrompt
    if err := c.ShouldBind(&aiPrompt); err != nil {
        util.JsonHttpResponse(c, 1, "参数错误", nil)
        return
    }

    userID, err := util.GetUserID(c)
    if err != nil {
        util.JsonHttpResponse(c, 2, err.Error(), nil)
        return
    }
    aiPrompt.UserID = userID

    dao.AddAIPrompt(&aiPrompt)
    util.JsonHttpResponse(c, 0, "AI提示词添加成功", aiPrompt)
}

func GetAIPrompts(c *gin.Context) {
    userID, err := util.GetUserID(c)
    if err != nil {
        util.JsonHttpResponse(c, 2, err.Error(), nil)
        return
    }

    prompts := dao.FindAIPromptsByUserID(userID)
    util.JsonHttpResponse(c, 0, "success", prompts)
}

func GetAIPromptByID(c *gin.Context) {
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

    prompt, found := dao.FindAIPromptByID(id)
    if !found {
        util.JsonHttpResponse(c, 1, "AI提示词不存在", nil)
        return
    }

    if prompt.UserID != userID && prompt.UserID != 0 {
        util.JsonHttpResponse(c, 1, "无权访问此提示词", nil)
        return
    }

    util.JsonHttpResponse(c, 0, "success", prompt)
}

func UpdateAIPrompt(c *gin.Context) {
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

    existingPrompt, found := dao.FindAIPromptByID(id)
    if !found {
        util.JsonHttpResponse(c, 1, "AI提示词不存在", nil)
        return
    }
    if existingPrompt.UserID != userID {
        util.JsonHttpResponse(c, 1, "无权修改此提示词", nil)
        return
    }

    var updateData map[string]interface{}
    if err := c.ShouldBindJSON(&updateData); err != nil {
        util.JsonHttpResponse(c, 1, "参数不合法: "+err.Error(), nil)
        return
    }

    log.Printf("更新AI提示词: userID=%d, promptID=%d, data=%+v", userID, id, updateData)

    updateData["id"] = id
    dao.UpdateAIPromptWithMap(updateData)

    updatedPrompt, _ := dao.FindAIPromptByID(id)
    util.JsonHttpResponse(c, 0, "AI提示词更新成功", updatedPrompt)
}

func DeleteAIPrompt(c *gin.Context) {
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

    prompt, found := dao.FindAIPromptByID(id)
    if !found {
        util.JsonHttpResponse(c, 1, "AI提示词不存在", nil)
        return
    }
    if prompt.UserID != userID {
        util.JsonHttpResponse(c, 1, "无权删除此提示词", nil)
        return
    }

    dao.DeleteAIPromptByID(id)
    util.JsonHttpResponse(c, 0, "AI提示词删除成功", nil)
}

func GetDefaultAIPrompt(c *gin.Context) {
    userID, err := util.GetUserID(c)
    if err != nil {
        util.JsonHttpResponse(c, 2, err.Error(), nil)
        return
    }

    prompt, found := dao.FindDefaultAIPromptByUserID(userID)
    if !found {
        util.JsonHttpResponse(c, 1, "未设置默认AI提示词", nil)
        return
    }

    util.JsonHttpResponse(c, 0, "success", prompt)
}
