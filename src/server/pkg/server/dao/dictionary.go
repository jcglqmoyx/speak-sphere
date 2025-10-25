package dao

import "speak-sphere/pkg/server/model"

func AddDictionary(dictionary *model.Dictionary) *model.Dictionary {
	db := GetDB()
	db.Create(&dictionary)
	return dictionary
}

func DeleteDictionaryByID(id int) {
	db := GetDB()
	db.Delete(&model.Dictionary{}, id)
}

func UpdateDictionary(dictionary *model.Dictionary) *model.Dictionary {
	db := GetDB()
	db.Model(&model.Dictionary{}).Where("id = ?", dictionary.ID).Updates(dictionary)
	return dictionary
}

func FindDictionaryByID(id int) (*model.Dictionary, bool) {
	db := GetDB()
	var dictionaries []*model.Dictionary
	db.Limit(1).Where("id = ?", id).Find(&dictionaries)
	if len(dictionaries) > 0 {
		return dictionaries[0], true
	} else {
		return nil, false
	}
}

func FindDictionaryByTitle(userID int, title string) (*model.Dictionary, bool) {
	db := GetDB()
	var dictionaries []*model.Dictionary
	db.Limit(1).Where("user_id = ? and title = ?", userID, title).Find(&dictionaries)
	if len(dictionaries) > 0 {
		return dictionaries[0], true
	} else {
		return nil, false
	}
}

func ListDictionary(userID int) []*model.Dictionary {
	db := GetDB()
	var dictionaries []*model.Dictionary
	db.Find(&dictionaries)
	db.Where("user_id = ?", userID).Find(&dictionaries)
	return dictionaries
}
