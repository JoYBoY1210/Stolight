package queue

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/storage"
)

func Worker(fileId string) error {
	file, err := models.GetFileByID(fileId)
	if err != nil {
		return fmt.Errorf("Failed to retrieve file: %v", err)
	}
	nodes := config.Cfg.StorageNodes
	bucket, err := models.GetBucketByID(file.BucketID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve bucket: %v", err)
	}
	stagePath := filepath.Join("./staging", fileId+".raw")
	Stagefile, err := os.Open(stagePath)
	if err != nil {
		return fmt.Errorf("Failed to open staged file: %v", err)
	}
	defer Stagefile.Close()
	err = storage.EncodeFile(Stagefile, file.Name, nodes, file.Size, bucket.Name)
	if err != nil {
		return fmt.Errorf("Failed to encode file: %v", err)
	}
	err = os.Remove(stagePath)
	if err != nil {
		fmt.Println("Could not delete staged file, will try again later")
	}
	return nil

}
