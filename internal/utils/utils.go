package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Printf("Error clearing screen: %v", err)
	}
}

func PrintWelcomeMessage() {
	fmt.Println("Welcome to GetCourse Video Downloader!")
	fmt.Println("https://github.com/dzatona/getcourse-video-downloader")
	fmt.Println("=====================================================")
	fmt.Println()
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(f)

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func ClearParts() {
	err := os.RemoveAll("parts")
	if err != nil {
		log.Println("Error removing parts directory:", err)
	}
}
