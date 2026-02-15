package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleMakeBucket(bucketName string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	reqBody, _ := json.Marshal(map[string]string{
		"name": bucketName,
	})
	req, _ := http.NewRequest("POST", cfg.ServerURL+"/api/buckets/", bytes.NewBuffer(reqBody))
	req.Header.Set("sto-Key", cfg.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error in making bucket,err: %s\n", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed: %s\n", string(body))
		return
	}

	fmt.Printf("Bucket '%s' created successfully!\n", bucketName)
}

func HandleList(bucketName string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	req, _ := http.NewRequest("GET", cfg.ServerURL+"/api/buckets/"+bucketName+"/files", nil)
	req.Header.Set("sto-Key", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error in listing files,err: %s\n", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed: %s\n", string(body))
		return
	}

	var result struct {
		Status string `json:"status"`
		Files  []struct {
			FileName string `json:"name"`
			Size     int64  `json:"size"`
		} `json:"files"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	files := result.Files
	json.NewDecoder(resp.Body).Decode(&files)

	fmt.Printf("\n Files in '%s':\n", bucketName)
	fmt.Println("------------------------------------------------")
	if len(files) == 0 {
		fmt.Println("(empty bucket)")
	}
	for _, f := range files {
		fmt.Printf("- %-20s  (%d bytes)\n", f.FileName, f.Size)
	}
	fmt.Println("------------------------------------------------")
}
