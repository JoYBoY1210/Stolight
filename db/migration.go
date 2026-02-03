package db

import (
	"fmt"

	"github.com/joyboy1210/stolight/models"
	"gorm.io/gorm"
)

func Mirgrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.File{}, &models.ChunkMetaData{})
	if err!=nil{
		return err
	}
	fmt.Println("Migration completed successfully")
	return nil
}
