package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
	"golang.org/x/term"
	"gorm.io/gorm"
)

func CheckOnStart(db *gorm.DB) {
	userCount, err := models.GetTotalUsers()
	if err != nil {
		fmt.Println("Failed to check total users: " + err.Error())
	}

	if userCount == 0 {
		fmt.Println("\n🔰 FIRST RUN DETECTED")
		fmt.Println("------------------------------------------------")
		fmt.Println("No users found. Please create the ROOT account.")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Username [root]: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			username = "root"
		}

		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("\n❌ Failed to read password: %v", err)
		}
		fmt.Println()
		password := string(bytePassword)

		if len(password) < 6 {
			log.Fatal("❌ Password must be at least 6 characters.")
		}

		hashedPwd, _ := HashPassword(password)
		initialKey := uuid.New().String()

		rootUser := models.User{
			ID:             uuid.New().String(),
			Username:       username,
			PasswordHash:   hashedPwd,
			Key:            initialKey,
			Role:           "admin",
			AllowedBuckets: "*",
		}

		if err := db.Create(&rootUser).Error; err != nil {
			log.Fatalf("❌ DB Error: %v", err)
		}

		fmt.Println("✅ Root user created successfully!")
		fmt.Printf("🔑 Initial API Key: %s\n", initialKey)
		fmt.Println("------------------------------------------------")
	}
}
