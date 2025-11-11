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

func AddBook(c *gin.Context) {
	var book *model.Book
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
	if err := c.ShouldBind(&book); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
		return
	}
	if _, found := dao.FindBookByTitle(userID, book.Title); found {
		util.JsonHttpResponse(c, 1, "相同标题的词书已存在, 请换一个标题", nil)
		return
	}
	book.UserID = user.ID

	file, _ := c.FormFile("file")
	if file == nil {
		book = dao.AddBook(book)
		util.JsonHttpResponse(c, 0, "成功创建了一个空词书", book)
		return
	} else {
		fileSize := file.Size
		if fileSize > conf.Cfg.Book.MaxFileSize {
			util.JsonHttpResponse(c, 1, "文件大小超过限制(10M)，请重新选择文件", nil)
			return
		}
		var entries []*model.Entry
		if strings.HasSuffix(file.Filename, ".txt") {
			entries, err = util.ParseTxtFile(c)
			if err != nil {
				util.JsonHttpResponse(c, 1, err.Error(), nil)
				return
			}
		} else if strings.HasSuffix(file.Filename, ".xlsx") {
			entries, err = util.ParseXlsxFile(c)
			if err != nil {
				util.JsonHttpResponse(c, 1, err.Error(), nil)
				return
			}
		} else {
			util.JsonHttpResponse(c, 1, "只支持 .txt 和 .xlsx 类型的文件 ", nil)
			return
		}

		path := filepath.Join(conf.Cfg.Book.UploadPath, util.DatetimeToString(time.Now()), file.Filename)
		book.FilePath = path
		book.MD5 = util.GetFileMD5(path)
		book = dao.AddBook(book)
		for i := 0; i < len(entries); i++ {
			entries[i].UserID = user.ID
			entries[i].BookID = book.ID
		}
		dao.BatchInsertEntry(entries)
		_ = c.SaveUploadedFile(file, path)
		util.JsonHttpResponse(c, 0, "success", book)
	}
}

func DeleteBookByID(c *gin.Context) {
	if bookID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if book, found := dao.FindBookByID(bookID); !found {
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
			if user.Level != 0 && user.ID != book.UserID {
				util.JsonHttpResponse(c, 1, "您没有权限删除该词书", nil)
				return
			}
			dao.DeleteBookByID(bookID)
			util.JsonHttpResponse(c, 0, "词书删除成功", nil)
		}
	}
}

func UpdateBook(c *gin.Context) {
	var book *model.Book
	if err := c.ShouldBind(&book); err != nil {
		util.JsonHttpResponse(c, 1, "参数不合法", nil)
		return
	}
	if _, found := dao.FindBookByID(book.ID); !found {
		util.JsonHttpResponse(c, 1, "该词书不存在", nil)
	} else {
		util.JsonHttpResponse(c, 0, "词书更新成功", dao.UpdateBook(book))
	}
}

func FindBookByID(c *gin.Context) {
	if bookID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "词书查询失败, 请输入合法的ID值", nil)
	} else {
		if book, found := dao.FindBookByID(bookID); !found {
			util.JsonHttpResponse(c, 1, "没查找到结果", nil)
		} else {
			util.JsonHttpResponse(c, 0, "success", book)
		}
	}
}

func ListBook(c *gin.Context) {
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
	util.JsonHttpResponse(c, 0, "success", dao.ListBook(pageSize, currentPage))
}

func FindBookByCategory(c *gin.Context) {
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
	util.JsonHttpResponse(c, 0, "success", dao.FindBooksByCategory(category, pageSize, currentPage))
}

func CountBook(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	}
	util.JsonHttpResponse(c, 0, "success", dao.CountBook(userID))
}
