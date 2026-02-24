package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
)

func StageFile(src io.Reader, fileName string, size int64, bucketName string) (string, int64, error) {

	fileID := uuid.New().String()
	file := &models.File{
		ID:       fileID,
		Name:     fileName,
		Size:     size,
		BucketID: bucketName,
		Status:   models.FileStatusStaging,
	}
	err := models.CreateFile(file)
	if err != nil {
		return "", 0, fmt.Errorf("could not create file record: %v", err)
	}

	err = os.MkdirAll("./staging", os.ModePerm)
	if err != nil {
		return "", 0, fmt.Errorf("could not create staging directory: %v", err)
	}
	stagePath := filepath.Join("./staging", fileID+".raw")
	out, err := os.Create(stagePath)
	if err != nil {
		return "", 0, fmt.Errorf("could not make the file: %v", err)
	}
	var success bool
	defer func() {
		out.Close()
		if !success {
			os.Remove(stagePath)
		}
	}()

	written, err := io.Copy(out, src)
	if err != nil {
		return "", 0, fmt.Errorf("could not write the file: %v", err)
	}
	if size > 0 && written != size {
		return "", 0, fmt.Errorf("incomplete upload: expected %d bytes, got %d", size, written)
	}
	success = true
	err = models.UpdateFileStatusAndSize(fileID, models.FileStatusPending, written)
	if err != nil {
		return "", 0, fmt.Errorf("could not update file status: %v", err)
	}
	//idhar worker we have to call later. we will add the file to the queue here.
	return fileID, written, nil
}
