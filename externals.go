package main

import (
	"fmt"
	"os"
	"os/exec"
)

func handleExternalCommand(name string, args ...string) {
	_, err := exec.LookPath(name)
	if err != nil {
		fmt.Println(name + ": command not found")
		return
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error:", err)
		return;
	}
}