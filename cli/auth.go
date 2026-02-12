package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func HandleLogin() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Server URL (e.g. http://localhost:6124): ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if url == "" {
		url = "http://localhost:6124"
	}
	url = strings.TrimSuffix(url, "/")

	fmt.Print("Username [root]: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		username = "root"
	}
	fmt.Print("Password: ")
	pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password")
		return
	}
	fmt.Println()
	password := string(pwdBytes)
	if password == "" {
		fmt.Println("Password cannot be empty")
		return
	}
	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	resp, err := http.Post(url+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf(" Connection failed: %v\nIs the server running?\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Login failed: %s\n", resp.Status)
		return
	}
	var result struct {
		Key  string `json:"key"`
		Role string `json:"role"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return
	}
	err = SaveConfig(result.Key, result.Role, url)
	if err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}
	fmt.Println("Login successful! Session saved.")
}
