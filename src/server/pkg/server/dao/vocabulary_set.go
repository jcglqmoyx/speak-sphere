package dao

import (
	"speak-sphere/pkg/server/model"
)

func AddVocabularySet(vocabularySet *model.VocabularySet) *model.VocabularySet {
	db := GetDB()
	db.Create(&vocabularySet)
	return vocabularySet
}

func DeleteVocabularySetByID(id int) {
	db := GetDB()
	db.Delete(&model.VocabularySet{}, id)
	db.Delete(&model.Vocabulary{}, "vocabulary_set_id = ?", id)
}

func UpdateVocabularySet(vocabularySet *model.VocabularySet) *model.VocabularySet {
	db := GetDB()
	db.Model(&model.VocabularySet{}).Where("id = ?", vocabularySet.ID).Updates(vocabularySet)
	return vocabularySet
}

func FindVocabularySetByID(id int) (*model.VocabularySet, bool) {
	db := GetDB()
	var vocabularySets []*model.VocabularySet
	db.Limit(1).Where("id= ?", id).Find(&vocabularySets)
	if len(vocabularySets) > 0 {
		return vocabularySets[0], true
	}
	return nil, false
}

func FindVocabularySetByTitle(userID int, title string) (*model.VocabularySet, bool) {
	db := GetDB()
	var vocabularySets []*model.VocabularySet
	db.Limit(1).Where("title = ? AND user_id = ?", title, userID).Find(&vocabularySets)
	if len(vocabularySets) > 0 {
		return vocabularySets[0], true
	}
	return nil, false
}

func ListVocabularySet(pageSize int, currentPage int) []*model.VocabularySet {
	db := GetDB()
	var vocabularySets []*model.VocabularySet
	db.Offset((currentPage - 1) * pageSize).Limit(pageSize).Find(&vocabularySets)
	return vocabularySets
}

func FindVocabularySetsByCategory(category string, pageSize int, currentPage int) []*model.VocabularySet {
	db := GetDB()
	var vocabularySets []*model.VocabularySet
	db.Offset((currentPage-1)*pageSize).Limit(pageSize).Where("category = ?", category).Find(&vocabularySets)
	return vocabularySets
}

func CountVocabularySet(userID int) int64 {
	db := GetDB()
	var count int64
	db.Model(&model.VocabularySet{}).Where("user_id = ?", userID).Count(&count)
	return count
}
