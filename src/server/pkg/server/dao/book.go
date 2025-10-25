package dao

import (
	"speak-sphere/pkg/server/model"
)

func AddBook(book *model.Book) *model.Book {
	db := GetDB()
	db.Create(&book)
	return book
}

func DeleteBookByID(id int) {
	db := GetDB()
	db.Delete(&model.Book{}, id)
	db.Delete(&model.Entry{}, "book_id = ?", id)
}

func UpdateBook(book *model.Book) *model.Book {
	db := GetDB()
	db.Model(&model.Book{}).Where("id = ?", book.ID).Updates(book)
	return book
}

func FindBookByID(id int) (*model.Book, bool) {
	db := GetDB()
	var books []*model.Book
	db.Limit(1).Where("id= ?", id).Find(&books)
	if len(books) > 0 {
		return books[0], true
	}
	return nil, false
}

func FindBookByTitle(userID int, title string) (*model.Book, bool) {
	db := GetDB()
	var books []*model.Book
	db.Limit(1).Where("title = ? AND user_id = ?", title, userID).Find(&books)
	if len(books) > 0 {
		return books[0], true
	}
	return nil, false
}

func ListBook(pageSize int, currentPage int) []*model.Book {
	db := GetDB()
	var books []*model.Book
	db.Offset((currentPage - 1) * pageSize).Limit(pageSize).Find(&books)
	return books
}

func FindBooksByCategory(category string, pageSize int, currentPage int) []*model.Book {
	db := GetDB()
	var books []*model.Book
	db.Offset((currentPage-1)*pageSize).Limit(pageSize).Where("category = ?", category).Find(&books)
	return books
}

func CountBook(userID int) int64 {
	db := GetDB()
	var count int64
	db.Model(&model.Book{}).Where("user_id = ?", userID).Count(&count)
	return count
}
