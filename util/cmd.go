package util

import (
	log "fmt"
	"os/exec"
)

func OpenURL(url string) {
	cmd := exec.Command("cmd", "/c", "start", url)
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open URL: %v\n", err)
	}
}

func ShowMessage(msg string) {
	cmd := exec.Command("cmd", "/c", "msg", "%username%", msg)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to show message: %v\n", err)
	}
}
