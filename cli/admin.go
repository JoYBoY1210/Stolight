package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func HandleCreateProject() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if cfg.Role != "admin" {
		fmt.Println("Only Root Admin can create projects. Please log in as Root Admin.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new project name (eg. myproject): ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		fmt.Println("Project name cannot be empty")
		return
	}

	fmt.Print("Allowed Buckets (comma separated, or *): ")
	buckets, _ := reader.ReadString('\n')
	buckets = strings.TrimSpace(buckets)

	reqBody, _ := json.Marshal(map[string]string{
		"name":    projectName,
		"buckets": buckets,
	})

	req, _ := http.NewRequest("POST", cfg.ServerURL+"/api/admin/projects/create", bytes.NewBuffer(reqBody))
	req.Header.Set("sto-Key", cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed: %s\n", string(body))
		return
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("------------------------------------------------")
	fmt.Println("Project Created Successfully")
	fmt.Printf("API Key:  %s\n", result["key"])
	fmt.Println("Copy this key. It will not be shown again.")
	fmt.Println("------------------------------------------------")
}
