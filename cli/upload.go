package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleUpload(localPath, remotePath string) {
	parts := strings.SplitN(remotePath, "/", 2)
	bucketName := parts[0]
	name := parts[1]

	file, err := os.Open(localPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	if name == "" {
		name = filepath.Base(localPath)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	fmt.Printf("Uploading %s to bucket %s...\n", localPath, bucketName)
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	writer.Close()
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	url := fmt.Sprintf("%s/api/upload/%s", cfg.ServerURL, bucketName)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("sto-Key", cfg.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n Network error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("\nUpload Failed: %s\n", string(respBody))
		return
	}
	fmt.Printf("\nUpload successful!\n")

}
