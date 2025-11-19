package controller

import (
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddVocabulary(c *gin.Context) {
	var vocabulary model.Vocabulary
	if err := c.ShouldBind(&vocabulary); err != nil {
		util.JsonHttpResponse(c, 1, err.Error(), nil)
		return
	}
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, "用户不存在, 无法添加", nil)
		return
	}

	vocabularySet, found := dao.FindVocabularySetByID(vocabulary.VocabularySetID)
	if !found {
		util.JsonHttpResponse(c, 1, "词书不存在，请先创建该词书", nil)
		return
	}
	if vocabularySet.UserID != userID {
		util.JsonHttpResponse(c, 1, "您无权添加该词条", nil)
		return
	}
	if _, found := dao.FindVocabularyByVocabulary(vocabulary.Vocabulary, vocabulary.VocabularySetID); found {
		util.JsonHttpResponse(c, 1, "单词已经存在于当前词书中，无需重复添加", nil)
	} else {
		vocabulary.UserID = userID
		dao.AddVocabulary(&vocabulary)
		util.JsonHttpResponse(c, 0, "success", vocabulary)
	}
}

func DeleteVocabularyByID(c *gin.Context) {
	if vocabularyID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if _, found := dao.FindVocabularyByID(vocabularyID); !found {
			util.JsonHttpResponse(c, 1, "该条目不存在, 无法删除", nil)
		} else {
			dao.DeleteVocabularyByID(vocabularyID)
			util.JsonHttpResponse(c, 0, "条目删除成功", nil)
		}
	}
}

func UpdateVocabulary(c *gin.Context) {
	var vocabulary *model.Vocabulary
	if err := c.ShouldBind(&vocabulary); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
	} else {
		if userID, err := util.GetUserID(c); err != nil {
			util.JsonHttpResponse(c, 1, "用户不存在, 无法更新", nil)
		} else {
			if obj, found := dao.FindVocabularyByID(vocabulary.ID); !found {
				util.JsonHttpResponse(c, 1, "该条目不存在, 无法更新", nil)
			} else {
				if obj.UserID != userID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					vocabulary.UserID = obj.UserID
					dao.UpdateVocabulary(vocabulary)
					util.JsonHttpResponse(c, 0, "success", "词条更新成功")
				}
			}
		}
	}
}

func SetVocabularyUnwanted(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		vocabularyID, _ := strconv.Atoi(id)
		if vocabulary, found := dao.FindVocabularyByID(vocabularyID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "用户不存在, 无法设置", nil)
			} else {
				if userID != vocabulary.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					dao.SetVocabularyUnwanted(vocabularyID)
					util.JsonHttpResponse(c, 0, "success", "词条设置为不想学习")
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}

}

func UpdateVocabularyStudyCount(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		vocabularyID, _ := strconv.Atoi(id)
		if vocabulary, found := dao.FindVocabularyByID(vocabularyID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
			} else {
				if userID != vocabulary.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					user, _ := dao.FindUserByID(userID)
					vocabulary.StudyCount++
					arr, _ := util.ParseReviewFrequencyFormula(user.ReviewFrequencyFormula)
					if vocabulary.StudyCount == 1 {
						vocabulary.DateToReview = util.AddDaysToIntDate(util.DateToInt(util.DateToString(time.Now())), arr[0])
					} else {
						if vocabulary.StudyCount > len(arr) {
							vocabulary.DateToReview = conf.Cfg.Vocabulary.DefaultDateToReview
						} else {
							vocabulary.DateToReview = util.AddDaysToIntDate(util.DateToInt(util.DateToString(time.Now())), arr[vocabulary.StudyCount-1])
						}
					}
					dao.UpdateVocabulary(vocabulary)
					util.JsonHttpResponse(c, 0, "success", vocabulary)
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func ResetStudyCountToZero(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		vocabularyID, _ := strconv.Atoi(id)
		if vocabulary, found := dao.FindVocabularyByID(vocabularyID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
			} else {
				if userID != vocabulary.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					vocabulary.StudyCount = 0
					dao.UpdateVocabulary(vocabulary)
					util.JsonHttpResponse(c, 0, "success", vocabulary)
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func FindVocabularyByID(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		vocabularyID, _ := strconv.Atoi(id)
		if vocabulary, found := dao.FindVocabularyByID(vocabularyID); found {
			util.JsonHttpResponse(c, 0, "success", vocabulary)
		} else {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func GetVocabulariesToLearn(c *gin.Context) {
	if userID, err := util.GetUserID(c); err != nil {
		util.JsonHttpResponse(c, 1, "用户不存在, 无法查询", nil)
	} else {
		user, _ := dao.FindUserByID(userID)
		util.JsonHttpResponse(c, 0, "success", dao.GetVocabulariesToLearn(user.CurrentVocabularySetID, user.DailyCount))
	}
}

func GetVocabulariesToReview(c *gin.Context) {
	userID, _ := util.GetUserID(c)
	util.JsonHttpResponse(c, 0, "success", dao.GetVocabulariesToReview(userID))
}

func CountVocabulary(c *gin.Context) {
	vocabularySetID, _ := strconv.Atoi(c.Param("vocabulary_set_id"))
	util.JsonHttpResponse(c, 0, "success", dao.CountVocabulary(vocabularySetID))
}

func ListVocabulary(c *gin.Context) {
	vocabularySetID, err := strconv.Atoi(c.Query("vocabulary_set_id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "error", "词条查询失败, 请输入合法的id值")
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "error", "词条查询失败, 请输入合法的page_size值")
	}
	currentPage, err := strconv.Atoi(c.Query("currentPage"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "error", "词条查询失败, 请输入合法的current_page值")
	}
	util.JsonHttpResponse(c, 0, "success", dao.ListVocabulary(vocabularySetID, pageSize, currentPage))
}

func CheckVocabularyInVocabularySet(c *gin.Context) {
	vocabulary := c.Query("vocabulary")
	vocabularySetID, err := strconv.Atoi(c.Query("vocabulary_set_id"))
	if err != nil {
		util.JsonHttpResponse(c, 1, "词书ID不合法", nil)
		return
	}

	if vocabulary == "" {
		util.JsonHttpResponse(c, 1, "单词不能为空", nil)
		return
	}

	// 检查用户是否有权限访问该词书
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, "用户不存在", nil)
		return
	}

	vocabularySet, found := dao.FindVocabularySetByID(vocabularySetID)
	if !found {
		util.JsonHttpResponse(c, 1, "词书不存在", nil)
		return
	}

	if vocabularySet.UserID != userID {
		util.JsonHttpResponse(c, 1, "您无权访问该词书", nil)
		return
	}

	// 检查单词是否在词书中
	_, exists := dao.FindVocabularyByVocabulary(vocabulary, vocabularySetID)
	util.JsonHttpResponse(c, 0, "success", map[string]bool{
		"exists": exists,
	})
}
