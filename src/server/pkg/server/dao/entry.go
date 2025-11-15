package dao

import (
	"log"
	"speak-sphere/pkg/server/model"
	"speak-sphere/pkg/server/util"
	"time"
)

func AddVocabulary(vocabulary *model.Vocabulary) *model.Vocabulary {
	db := GetDB()
	db.Create(vocabulary)
	return vocabulary
}

func BatchInsertVocabulary(vocabularies []*model.Vocabulary) {
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
	for i := 0; i < len(vocabularies); i += batchSize {
		end := i + batchSize
		if end > len(vocabularies) {
			end = len(vocabularies)
		}

		if err := tx.Create(vocabularies[i:end]).Error; err != nil {
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

func DeleteVocabularyByID(id int) {
	db := GetDB()
	db.Delete(&model.Vocabulary{}, id)
}

func UpdateVocabulary(vocabulary *model.Vocabulary) {
	db := GetDB()
	db.Model(&model.Vocabulary{}).Where("id = ?", vocabulary.ID).Updates(map[string]interface{}{"vocabulary": vocabulary.Vocabulary, "meaning": vocabulary.Meaning, "note": vocabulary.Note, "unwanted": vocabulary.Unwanted, "study_count": vocabulary.StudyCount, "date_to_review": vocabulary.DateToReview})
}

func SetVocabularyUnwanted(vocabularyID int) {
	db := GetDB()
	db.Model(&model.Vocabulary{}).Where("id = ?", vocabularyID).Update("unwanted", true)
}

func FindVocabularyByVocabulary(vocabulary string, vocabularySetID int) (*model.Vocabulary, bool) {
	db := GetDB()
	var vocabularies []*model.Vocabulary
	db.Limit(1).Where("vocabulary = ? AND vocabulary_set_id = ?", vocabulary, vocabularySetID).Find(&vocabularies)
	if len(vocabularies) > 0 {
		return vocabularies[0], true
	} else {
		return nil, false
	}
}

func FindVocabularyByID(id int) (*model.Vocabulary, bool) {
	db := GetDB()
	var vocabularies []*model.Vocabulary
	db.Limit(1).Where("id = ?", id).Find(&vocabularies)
	if len(vocabularies) > 0 {
		return vocabularies[0], true
	}
	return nil, false
}

func GetVocabulariesToLearn(vocabularySetID int, count int) []*model.Vocabulary {
	db := GetDB()
	var vocabularies []*model.Vocabulary
	db.Limit(count).Where("vocabulary_set_id = ? AND study_count = ? AND unwanted = ?", vocabularySetID, 0, false).Find(&vocabularies)
	return vocabularies
}

func GetVocabulariesToReview(userID int) []*model.Vocabulary {
	db := GetDB()
	var vocabularies []*model.Vocabulary
	db.Where("user_id = ? AND date_to_review <= ?", userID, util.DateToString(time.Now())).Find(&vocabularies)
	return vocabularies
}

func CountVocabulary(vocabularySetID int) int64 {
	db := GetDB()
	var count int64
	db.Model(&model.Vocabulary{}).Where("vocabulary_set_id = ?", vocabularySetID).Count(&count)
	return count
}

func ListVocabulary(vocabularySetID int, pageSize int, currentPage int) []*model.Vocabulary {
	db := GetDB()
	var vocabularies []*model.Vocabulary
	db.Where("vocabulary_set_id = ?", vocabularySetID).Offset((currentPage - 1) * pageSize).Limit(pageSize).Find(&vocabularies)
	return vocabularies
}
