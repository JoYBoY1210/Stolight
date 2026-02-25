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
	nodes := config.Cfg.StorageNodes

	stagePath := filepath.Join("./staging", fileId+".raw")
	Stagefile, err := os.Open(stagePath)
	if err != nil {
		models.UpdateFileStatus(fileId, models.FileStatusFailed)
		return fmt.Errorf("Failed to open staged file: %v", err)
	}

	err = models.UpdateFileStatus(fileId, models.FileStatusEncoding)
	if err != nil {
		Stagefile.Close()
		return fmt.Errorf("Failed to update file status: %v", err)
	}

	err = storage.EncodeFile(Stagefile, fileId, nodes)

	Stagefile.Close()

	if err != nil {
		models.UpdateFileStatus(fileId, models.FileStatusFailed)
		return fmt.Errorf("Failed to encode file: %v", err)
	}

	err = models.UpdateFileStatus(fileId, models.FileStatusCompleted)
	if err != nil {
		return fmt.Errorf("Failed to update file status: %v", err)
	}

	err = os.Remove(stagePath)
	if err != nil {
		fmt.Println("Could not delete staged file, will try again later")
	}

	return nil
}
