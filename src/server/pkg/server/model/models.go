package model

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID        int            `gorm:"primaryKey" json:"id" form:"id"`
	Title     string         `gorm:"column:title" form:"title" json:"title,omitempty"`
	Category  string         `gorm:"column:category;default:'Uncategorized'" form:"category" json:"category,omitempty"`
	UserID    int            `gorm:"column:user_id" form:"user_id" json:"user_id"`
	MD5       string         `gorm:"column:md5" form:"md5" json:"md5,omitempty"`
	FilePath  string         `gorm:"column:file_path" form:"file_path" json:"file_path,omitempty"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

type Dictionary struct {
	ID        int            `gorm:"primaryKey" json:"id" form:"id"`
	Title     string         `gorm:"column:title" form:"title" json:"title,omitempty"`
	Prefix    string         `gorm:"column:prefix" form:"prefix" json:"prefix,omitempty"`
	Suffix    string         `gorm:"column:suffix" form:"suffix" json:"suffix,omitempty"`
	UserID    int            `gorm:"column:user_id" form:"user_id" json:"user_id"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

type Entry struct {
	ID           int            `gorm:"primaryKey" json:"id" form:"id"`
	Word         string         `gorm:"column:word" form:"word" json:"word"`
	Meaning      string         `gorm:"column:meaning" form:"meaning" json:"meaning"`
	BookID       int            `gorm:"column:book_id" form:"book_id" json:"book_id"`
	UserID       int            `gorm:"column:user_id" form:"user_id" json:"user_id"`
	Note         string         `gorm:"column:note" form:"note" json:"note"`
	Unwanted     bool           `gorm:"column:unwanted" form:"unwanted" json:"unwanted"`
	StudyCount   int            `gorm:"column:study_count" form:"study_count" json:"study_count"`
	DateToReview int            `gorm:"column:date_to_review;default:99991231" form:"date_to_review" json:"date_to_review"`
	CreatedAt    time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

type User struct {
	ID                     int            `gorm:"primaryKey" json:"id" form:"id"`
	Username               string         `gorm:"column:username;unique" form:"username" json:"username,omitempty"`
	Email                  string         `gorm:"column:email;unique" form:"email" json:"email,omitempty"`
	Avatar                 string         `gorm:"column:avatar" form:"avatar" json:"avatar"`
	Password               string         `gorm:"column:password" form:"password" json:"password,omitempty"`
	PasswordSalt           string         `gorm:"column:password_salt" form:"password_salt" json:"password_salt,omitempty"`
	Level                  int            `gorm:"column:level;default:1" form:"level" json:"level"`
	TotalEntryCount        int            `gorm:"total_entry_count" form:"total_entry_count" json:"total_entry_count"`
	CurrentBookID          int            `gorm:"current_book_id" form:"current_book_id" json:"current_book_id"`
	DailyCount             int            `gorm:"daily_count;default:10" form:"daily_count" json:"daily_count"`
	TimesCountedAsKnown    int            `gorm:"times_counted_as_known;default:2" form:"times_counted_as_known" json:"times_counted_as_known"`
	ReviewFrequencyFormula string         `gorm:"column:review_frequency_formula;default:'1_3_7_15_30_90_180_'" form:"review_frequency_formula" json:"review_frequency_formula"`
	CreatedAt              time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
