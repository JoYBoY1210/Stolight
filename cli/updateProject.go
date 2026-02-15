package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleUpdateProject(projectName, allowedBuckets string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	reqBody, _ := json.Marshal(map[string]string{
		"username":        projectName,
		"allowed_buckets": allowedBuckets,
	})
	req, _ := http.NewRequest("POST", cfg.ServerURL+"/api/admin/projects/update", bytes.NewBuffer(reqBody))
	req.Header.Set("sto-Key", cfg.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error in updating project,err: %s\n", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed: %s\n", string(body))
		return
	}
	fmt.Printf("Project '%s' updated successfully with buckets: %s!\n", projectName, allowedBuckets)
}
