package dao

import (
	"speak-sphere/pkg/server/conf"

	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return conf.DB.Session(&gorm.Session{})
}
