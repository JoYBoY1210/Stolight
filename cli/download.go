package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func HandleDownload(remotePath, localPath string) {
	parts := strings.SplitN(remotePath, "/", 2)
	if len(parts) < 2 {
		println("Invalid remote path. Use bucket/filename format.")
		return
	}
	bucket := parts[0]
	fileName := parts[1]

	cfg, err := LoadConfig()
	if err != nil {
		println("Error loading config:", err.Error())
		return
	}
	if cfg.Token == "" {
		println("Not logged in. Please run 'sto login' first.")
		return
	}
	if cfg.Role != "admin" {
		println("Only admin users can download files.")
		return
	}
	req, err := http.NewRequest("GET", cfg.ServerURL+"/api/download/"+bucket+"/"+fileName, nil)
	if err != nil {
		println("Error creating request:", err.Error())
		return
	}
	req.Header.Set("sto-Key", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Error making request:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println("Download failed with status:", resp.Status)
		return
	}

	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err.Error())
		return
	}
	fmt.Println("File downloaded successfully as", localPath)
}
