package controller

import (
	"path/filepath"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddVocabularySet(c *gin.Context) {
	var vocabularySet *model.VocabularySet
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}
	user, found := dao.FindUserByID(userID)
	if !found {
		util.JsonHttpResponse(c, 1, "用户不存在", nil)
		return
	}
	if err := c.ShouldBind(&vocabularySet); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
		return
	}
	if _, found := dao.FindVocabularySetByTitle(userID, vocabularySet.Title); found {
		util.JsonHttpResponse(c, 1, "相同标题的词书已存在, 请换一个标题", nil)
		return
	}
	vocabularySet.UserID = user.ID

	file, _ := c.FormFile("file")
	if file == nil {
		vocabularySet = dao.AddVocabularySet(vocabularySet)
		util.JsonHttpResponse(c, 0, "成功创建了一个空词书", vocabularySet)
		return
	} else {
		fileSize := file.Size
		if fileSize > conf.Cfg.VocabularySet.MaxFileSize {
			util.JsonHttpResponse(c, 1, "文件大小超过限制(10M)，请重新选择文件", nil)
			return
		}
		var vocabularies []*model.Vocabulary
		if strings.HasSuffix(file.Filename, ".txt") {
			vocabularies, err = util.ParseTxtFile(c)
			if err != nil {
				util.JsonHttpResponse(c, 1, err.Error(), nil)
				return
			}
		} else if strings.HasSuffix(file.Filename, ".xlsx") {
			vocabularies, err = util.ParseXlsxFile(c)
			if err != nil {
				util.JsonHttpResponse(c, 1, err.Error(), nil)
				return
			}
		} else {
			util.JsonHttpResponse(c, 1, "只支持 .txt 和 .xlsx 类型的文件 ", nil)
			return
		}

		path := filepath.Join(conf.Cfg.VocabularySet.UploadPath, util.DatetimeToString(time.Now()), file.Filename)
		vocabularySet.FilePath = path
		vocabularySet.MD5 = util.GetFileMD5(path)
		vocabularySet = dao.AddVocabularySet(vocabularySet)
		for i := 0; i < len(vocabularies); i++ {
			vocabularies[i].UserID = user.ID
			vocabularies[i].VocabularySetID = vocabularySet.ID
		}
		dao.BatchInsertVocabulary(vocabularies)
		_ = c.SaveUploadedFile(file, path)
		util.JsonHttpResponse(c, 0, "success", vocabularySet)
	}
}

func DeleteVocabularySetByID(c *gin.Context) {
	if vocabularySetID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if vocabularySet, found := dao.FindVocabularySetByID(vocabularySetID); !found {
			util.JsonHttpResponse(c, 1, "该词书不存在, 无法删除", nil)
		} else {
			userID, err := util.GetUserID(c)
			if err != nil {
				util.JsonHttpResponse(c, 1, err.Error(), nil)
				return
			}
			user, found := dao.FindUserByID(userID)
			if !found {
				util.JsonHttpResponse(c, 1, "用户不存在", nil)
				return
			}
			if user.Level != 0 && user.ID != vocabularySet.UserID {
				util.JsonHttpResponse(c, 1, "您没有权限删除该词书", nil)
				return
			}
			dao.DeleteVocabularySetByID(vocabularySetID)
			util.JsonHttpResponse(c, 0, "词书删除成功", nil)
		}
	}
}

func UpdateVocabularySet(c *gin.Context) {
	var vocabularySet *model.VocabularySet
	if err := c.ShouldBind(&vocabularySet); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
		return
	}
	if _, found := dao.FindVocabularySetByID(vocabularySet.ID); !found {
		util.JsonHttpResponse(c, 1, "该词书不存在", nil)
	} else {
		util.JsonHttpResponse(c, 0, "词书更新成功", dao.UpdateVocabularySet(vocabularySet))
	}
}

func FindVocabularySetByID(c *gin.Context) {
	if vocabularySetID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "词书查询失败, 请输入合法的ID值", nil)
	} else {
		if vocabularySet, found := dao.FindVocabularySetByID(vocabularySetID); !found {
			util.JsonHttpResponse(c, 1, "没查找到结果", nil)
		} else {
			util.JsonHttpResponse(c, 0, "success", vocabularySet)
		}
	}
}

func ListVocabularySet(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Param("page_size"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "page_size不合法", nil)
		return
	}
	currentPage, err := strconv.Atoi(c.Param("current_page"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "current_page不合法", nil)
		return
	}
	util.JsonHttpResponse(c, 0, "success", dao.ListVocabularySet(pageSize, currentPage))
}

func FindVocabularySetByCategory(c *gin.Context) {
	category := c.Param("category")
	pageSize, err := strconv.Atoi(c.Param("page_size"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "page_size不合法", nil)
		return
	}
	currentPage, err := strconv.Atoi(c.Param("current_page"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "current_page不合法", nil)
		return
	}
	util.JsonHttpResponse(c, 0, "success", dao.FindVocabularySetsByCategory(category, pageSize, currentPage))
}

func CountVocabularySet(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}
	util.JsonHttpResponse(c, 0, "success", dao.CountVocabularySet(userID))
}
