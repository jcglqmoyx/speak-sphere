package conf

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"speak-sphere/pkg/server/model"
)

var DB *gorm.DB

func InitDatabase(cfg *DBConfig) {
	DB, _ = gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	db, _ := DB.DB()
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour * time.Duration(cfg.ConnMaxLifetime))
	err := DB.AutoMigrate(&model.Book{})
	err = DB.AutoMigrate(&model.Entry{})
	err = DB.AutoMigrate(&model.Dictionary{})
	err = DB.AutoMigrate(&model.User{})
	err = DB.AutoMigrate(&model.LLMService{})
	err = DB.AutoMigrate(&model.AIPrompt{})
	if err != nil {
		println("Error occurred when opening database file")
	}
}
