package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
)

func AddDictionary(c *gin.Context) {
	var dict *model.Dictionary
	if err := c.ShouldBindJSON(&dict); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
		return
	}

	if dict.Title == "" {
		util.JsonHttpResponse(c, 1, "词典标题不能为空", nil)
		return
	}

	if len(dict.Title) > conf.Cfg.Dictionary.MaxTitleLength {
		util.JsonHttpResponse(c, 1, "词典标题长度超过限制", nil)
		return
	}

	if len(dict.Prefix) > conf.Cfg.Dictionary.MaxPrefixLength {
		util.JsonHttpResponse(c, 1, "前根长度超过限制", nil)
		return
	}

	if len(dict.Suffix) > conf.Cfg.Dictionary.MaxSuffixLength {
		util.JsonHttpResponse(c, 1, "后缀长度超过限制", nil)
		return
	}

	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, "用户未登录", nil)
		return
	}

	if _, found := dao.FindDictionaryByTitle(userID, dict.Title); found {
		util.JsonHttpResponse(c, 1, "相同标题的词典已存在", nil)
		return
	}

	dict.UserID = userID
	dao.AddDictionary(dict)
	util.JsonHttpResponse(c, 0, "词典添加成功", dict)
}

func DeleteDictionaryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "id参数错误", nil)
		return
	}
	if _, found := dao.FindDictionaryByID(id); !found {
		util.JsonHttpResponse(c, 1, "该词典不存在, 无法删除", nil)
	} else {
		dao.DeleteDictionaryByID(id)
		util.JsonHttpResponse(c, 0, "删除成功", nil)
	}
}

func UpdateDictionary(c *gin.Context) {
	var dictionary *model.Dictionary
	if err := c.ShouldBind(&dictionary); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", err.Error())
		return
	}
	if _, found := dao.FindDictionaryByID(dictionary.ID); !found {
		util.JsonHttpResponse(c, 1, "该词典不还在，无法更新", nil)
	}
	util.JsonHttpResponse(c, 0, "更新成功", dao.UpdateDictionary(dictionary))
}

func FindDictionaryByID(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		dictionaryID, _ := strconv.Atoi(id)
		if data, found := dao.FindDictionaryByID(dictionaryID); found {
			util.JsonHttpResponse(c, 0, "success", data)
		} else {
			util.JsonHttpResponse(c, 1, "该字典不存在", nil)
		}
	}
}

func ListDictionary(c *gin.Context) {
	if userID, err := util.GetUserID(c); err != nil {
		util.JsonHttpResponse(c, 1, "用户未登录", nil)
		return
	} else {
		util.JsonHttpResponse(c, 0, "success", dao.ListDictionary(userID))
	}
}
