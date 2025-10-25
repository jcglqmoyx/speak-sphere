package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/dao"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
)

func AddEntry(c *gin.Context) {
	var entry model.Entry
	if err := c.ShouldBind(&entry); err != nil {
		util.JsonHttpResponse(c, 1, err.Error(), nil)
		return
	}
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, "用户不存在, 无法添加", nil)
		return
	}

	book, found := dao.FindBookByID(entry.BookID)
	if !found {
		util.JsonHttpResponse(c, 1, "词书不存在，请先创建该词书", nil)
		return
	}
	if book.UserID != userID {
		util.JsonHttpResponse(c, 1, "您无权添加该词条", nil)
		return
	}
	if _, found := dao.FindEntryByWord(entry.Word, entry.BookID); found {
		util.JsonHttpResponse(c, 1, "单词已经存在于当前词书中，无需重复添加", nil)
	} else {
		entry.UserID = userID
		dao.AddEntry(&entry)
		util.JsonHttpResponse(c, 0, "success", entry)
	}
}

func DeleteEntryByID(c *gin.Context) {
	if entryID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if _, found := dao.FindEntryByID(entryID); !found {
			util.JsonHttpResponse(c, 1, "该条目不存在, 无法删除", nil)
		} else {
			dao.DeleteEntryByID(entryID)
			util.JsonHttpResponse(c, 0, "条目删除成功", nil)
		}
	}
}

func UpdateEntry(c *gin.Context) {
	var entry *model.Entry
	if err := c.ShouldBind(&entry); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
	} else {
		if userID, err := util.GetUserID(c); err != nil {
			util.JsonHttpResponse(c, 1, "用户不存在, 无法更新", nil)
		} else {
			if obj, found := dao.FindEntryByID(entry.ID); !found {
				util.JsonHttpResponse(c, 1, "该条目不存在, 无法更新", nil)
			} else {
				if obj.UserID != userID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					entry.UserID = obj.UserID
					dao.UpdateEntry(entry)
					util.JsonHttpResponse(c, 0, "success", "词条更新成功")
				}
			}
		}
	}
}

func SetEntryUnwanted(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		entryID, _ := strconv.Atoi(id)
		if entry, found := dao.FindEntryByID(entryID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "用户不存在, 无法设置", nil)
			} else {
				if userID != entry.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					dao.SetEntryUnwanted(entryID)
					util.JsonHttpResponse(c, 0, "success", "词条设置为不想学习")
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}

}

func UpdateEntryStudyCount(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		entryID, _ := strconv.Atoi(id)
		if entry, found := dao.FindEntryByID(entryID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
			} else {
				if userID != entry.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					user, _ := dao.FindUserByID(userID)
					entry.StudyCount++
					arr, _ := util.ParseReviewFrequencyFormula(user.ReviewFrequencyFormula)
					if entry.StudyCount == 1 {
						entry.DateToReview = util.AddDaysToIntDate(util.DateToInt(util.DateToString(time.Now())), arr[0])
					} else {
						if entry.StudyCount > len(arr) {
							entry.DateToReview = conf.Cfg.Entry.DefaultDateToReview
						} else {
							entry.DateToReview = util.AddDaysToIntDate(util.DateToInt(util.DateToString(time.Now())), arr[entry.StudyCount-1])
						}
					}
					dao.UpdateEntry(entry)
					util.JsonHttpResponse(c, 0, "success", entry)
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func ResetStudyCountToZero(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		entryID, _ := strconv.Atoi(id)
		if entry, found := dao.FindEntryByID(entryID); !found {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		} else {
			if userID, err := util.GetUserID(c); err != nil {
				util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
			} else {
				if userID != entry.UserID {
					util.JsonHttpResponse(c, 1, "您无权修改该条目", nil)
				} else {
					entry.StudyCount = 0
					dao.UpdateEntry(entry)
					util.JsonHttpResponse(c, 0, "success", entry)
				}
			}
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func FindEntryByID(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		entryID, _ := strconv.Atoi(id)
		if entry, found := dao.FindEntryByID(entryID); found {
			util.JsonHttpResponse(c, 0, "success", entry)
		} else {
			util.JsonHttpResponse(c, 1, "未找到相关条目", nil)
		}
	} else {
		util.JsonHttpResponse(c, 1, "词条查询失败, 请输入合法的ID值", nil)
	}
}

func GetEntriesToLearn(c *gin.Context) {
	if userID, err := util.GetUserID(c); err != nil {
		util.JsonHttpResponse(c, 1, "用户不存在, 无法查询", nil)
	} else {
		user, _ := dao.FindUserByID(userID)
		util.JsonHttpResponse(c, 0, "success", dao.GetEntriesToLearn(user.CurrentBookID, user.DailyCount))
	}
}

func GetEntriesToReview(c *gin.Context) {
	userID, _ := util.GetUserID(c)
	util.JsonHttpResponse(c, 0, "success", dao.GetEntriesToReview(userID))
}

func CountEntry(c *gin.Context) {
	bookID, _ := strconv.Atoi(c.Param("book_id"))
	util.JsonHttpResponse(c, 0, "success", dao.CountEntry(bookID))
}

func ListEntry(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("book_id"))
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
	util.JsonHttpResponse(c, 0, "success", dao.ListEntry(bookID, pageSize, currentPage))
}
