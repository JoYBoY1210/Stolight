package main

import (
	"fmt"
	"os"

	"github.com/joyboy1210/stolight/cli"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	command := os.Args[1]
	switch command {
	case "login":
		cli.HandleLogin()
	case "create-project":
		cli.HandleCreateProject()
	case "mb":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sto mb <bucket>")
			return
		}
		cli.HandleMakeBucket(os.Args[2])
	case "cp":

		if len(os.Args) < 4 {
			fmt.Println("Usage: sto cp <file> <bucket/path>")
			return
		}
		cli.HandleUpload(os.Args[2], os.Args[3])

	case "get":
		if len(os.Args) < 4 {
			fmt.Println("Usage: sto get <bucket/filename> <local-file>")
			return
		}
		cli.HandleDownload(os.Args[2], os.Args[3])
	case "rm":

		if len(os.Args) < 2 {
			fmt.Println("Usage: sto rm <bucket/path>")
			return
		}
		cli.HandleDelete(os.Args[2])

	case "ls":

		if len(os.Args) < 2 {
			fmt.Println("Usage: sto ls <bucket>")
			return
		}
		cli.HandleList(os.Args[2])
	case "update-project":

		if len(os.Args) < 4 {
			fmt.Println("Usage: sto update-project <project-name> <allowed-buckets>")
			return
		}
		cli.HandleUpdateProject(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("\n☁️  StoLight CLI (sto)")
	fmt.Println("------------------------------------------------")
	fmt.Println("  login                         -> Log in as Root Admin")
	fmt.Println("  create-project                -> Create a new API User")
	fmt.Println("  update-project                -> Update an existing project permissions")
	fmt.Println("  mb <bucket>                   -> Create a new bucket")
	fmt.Println("  ls <bucket>                   -> List files in a bucket")
	fmt.Println("  cp <local-file> <bucket/path> -> Upload a file")
	fmt.Println("  get <bucket/filename> <local-file> -> Download a file")
	fmt.Println("  rm <bucket/path>              -> Delete a file")
	fmt.Println("------------------------------------------------")
}
