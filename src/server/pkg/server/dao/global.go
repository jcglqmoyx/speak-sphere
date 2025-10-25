package dao

import (
	"gorm.io/gorm"
	"speak-sphere/pkg/server/conf"
)

func GetDB() *gorm.DB {
	return conf.DB.Session(&gorm.Session{})
}
