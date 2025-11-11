package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetFileMD5(filename string) string {
	f, _ := os.Open(filename)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	md5Handle := md5.New()
	_, _ = io.Copy(md5Handle, f)
	md := md5Handle.Sum(nil)
	md5str := fmt.Sprintf("%x", md)
	return md5str
}

func ParseTxtFile(c *gin.Context) ([]*model.Entry, error) {
	file, _, _ := c.Request.FormFile("file")
	content, _ := io.ReadAll(file)
	entries := strings.Split(string(content), "\n")

	userID, _ := GetUserID(c)

	var res []*model.Entry
	for row, entry := range entries {
		s := strings.Trim(entry, "\n")
		s = strings.Trim(s, " ")
		s = strings.Trim(s, "\t")
		if len(s) > conf.Cfg.Entry.MaxWordLength {
			return nil, fmt.Errorf("单词长度过长: %s, 第(%d)行\n", s, row+1)
		}
		if len(s) > 0 {
			res = append(res, &model.Entry{
				Word:   s,
				UserID: userID,
			})
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("词书内没有内容")
	}

	return res, nil
}

func ParseXlsxFile(c *gin.Context) ([]*model.Entry, error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("文件为空，请检查是否上传了文件")
	}
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("Excel文件无法打开")
	}

	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("Excel文件无法读取")
	}

	var res []*model.Entry
	for _, row := range rows {
		if len(row[0]) > conf.Cfg.Entry.MaxWordLength {
			return nil, fmt.Errorf("单词长度过长: %s", row[0])
		}
		var entry = &model.Entry{
			Word: row[0],
		}
		if len(row) > 1 {
			if len(row[1]) > conf.Cfg.Entry.MaxMeaningLength {
				return nil, fmt.Errorf("释义长度过长: %s", row[1])
			}
			entry.Meaning = row[1]
		}
		res = append(res, entry)
	}
	return res, nil
}
