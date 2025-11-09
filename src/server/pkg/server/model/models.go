package model

import (
	"time"

	"gorm.io/gorm"
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

type LLMService struct {
	ID          int            `gorm:"primaryKey" json:"id" form:"id"`
	Name        string         `gorm:"column:name" json:"name" form:"name"`
	Provider    string         `gorm:"column:provider" form:"provider" json:"provider"`
	Endpoint    string         `gorm:"column:endpoint" form:"endpoint" json:"endpoint,omitempty"`
	Model       string         `gorm:"column:model" form:"model" json:"model,omitempty"`
	APIKey      string         `gorm:"column:api_key" form:"api_key" json:"api_key,omitempty"`
	IsDefault   bool           `gorm:"column:is_default;default:false" json:"is_default"`
	Description string         `gorm:"column:description" json:"description"`
	UserID      int            `gorm:"column:user_id" form:"user_id" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

// AIPrompt AI提示词管理
type AIPrompt struct {
	ID          int            `gorm:"primaryKey" json:"id" form:"id"`
	Name        string         `gorm:"column:name" json:"name" form:"name"`                          // 提示词名称
	Content     string         `gorm:"column:content;type:text" json:"content" form:"content"`       // 提示词内容
	Description string         `gorm:"column:description" json:"description" form:"description"`      // 提示词描述
	IsDefault   bool           `gorm:"column:is_default;default:false" json:"is_default"`            // 是否默认提示词
	UserID      int            `gorm:"column:user_id" form:"user_id" json:"user_id"`                 // 用户ID，0表示系统默认
	CreatedAt   time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
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
	ReviewFrequencyFormula string         `gorm:"column:review_frequency_formula;default:'2_4_8_16_32_64_128_256_512_'" form:"review_frequency_formula" json:"review_frequency_formula"`
	CreatedAt              time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
