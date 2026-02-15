package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HandleDelete(remotePath string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	parts := strings.SplitN(remotePath, "/", 2)
	bucketName := parts[0]
	fileName := parts[1]
	url := fmt.Sprintf("%s/api/buckets/%s/files/%s", cfg.ServerURL, bucketName, url.PathEscape(fileName))
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("sto-Key", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error in deleting file,err: %s\n", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed: %s\n", string(body))
		return
	}
	fmt.Printf("File '%s' deleted successfully from bucket '%s'!\n", fileName, bucketName)
}
