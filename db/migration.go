package db

import (
	"github.com/joyboy1210/stolight/models"
	"gorm.io/gorm"
)

func Mirgrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.File{}, &models.ChunkMetaData{})
}
