package dao

import (
	"log"
	"time"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
)

func AddEntry(entry *model.Entry) *model.Entry {
	db := GetDB()
	db.Create(entry)
	return entry
}

func BatchInsertEntry(entries []*model.Entry) {
	db := GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		log.Printf("Transaction start error: %v", tx.Error)
		return
	}

	batchSize := 1000
	for i := 0; i < len(entries); i += batchSize {
		end := i + batchSize
		if end > len(entries) {
			end = len(entries)
		}

		if err := tx.Create(entries[i:end]).Error; err != nil {
			tx.Rollback()
			log.Printf("Batch insert error: %v", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit error: %v", err)
		return
	}
}

func DeleteEntryByID(id int) {
	db := GetDB()
	db.Delete(&model.Entry{}, id)
}

func UpdateEntry(entry *model.Entry) {
	db := GetDB()
	db.Model(&model.Entry{}).Where("id = ?", entry.ID).Updates(map[string]interface{}{"word": entry.Word, "meaning": entry.Meaning, "note": entry.Note, "unwanted": entry.Unwanted, "study_count": entry.StudyCount, "date_to_review": entry.DateToReview})
}

func SetEntryUnwanted(entryID int) {
	db := GetDB()
	db.Model(&model.Entry{}).Where("id = ?", entryID).Update("unwanted", true)
}

func FindEntryByWord(word string, bookID int) (*model.Entry, bool) {
	db := GetDB()
	var entries []*model.Entry
	db.Limit(1).Where("word = ? AND book_id = ?", word, bookID).Find(&entries)
	if len(entries) > 0 {
		return entries[0], true
	} else {
		return nil, false
	}
}

func FindEntryByID(id int) (*model.Entry, bool) {
	db := GetDB()
	var entries []*model.Entry
	db.Limit(1).Where("id = ?", id).Find(&entries)
	if len(entries) > 0 {
		return entries[0], true
	}
	return nil, false
}

func GetEntriesToLearn(bookID int, count int) []*model.Entry {
	db := GetDB()
	var entries []*model.Entry
	db.Limit(count).Where("book_id = ? AND study_count = ? AND unwanted = ?", bookID, 0, false).Find(&entries)
	return entries
}

func GetEntriesToReview(userID int) []*model.Entry {
	db := GetDB()
	var entries []*model.Entry
	db.Where("user_id = ? AND date_to_review <= ?", userID, util.DateToString(time.Now())).Find(&entries)
	return entries
}

func CountEntry(bookID int) int64 {
	db := GetDB()
	var count int64
	db.Model(&model.Entry{}).Where("book_id = ?", bookID).Count(&count)
	return count
}

func ListEntry(bookID int, pageSize int, currentPage int) []*model.Entry {
	db := GetDB()
	var entries []*model.Entry
	db.Where("book_id = ?", bookID).Offset((currentPage - 1) * pageSize).Limit(pageSize).Find(&entries)
	return entries
}
