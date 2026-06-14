package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = strings.TrimSpace(command)
		if command == "exit" {
			break
		}
		if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
			continue
		}
		if strings.HasPrefix(command, "type ") {
			arg := command[5:]
			executeTypeCommand(arg)
			continue
		}
		fmt.Println(command + ": command not found")
	}
}

func executeTypeCommand(arg string) {
	builtins := map[string]bool{
		"echo": true,
		"exit": true,
		"type": true,
	}
	if builtins[arg] {
		fmt.Println(arg + " is a shell builtin")
		return
	} 

	path, err := exec.LookPath(arg)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println(arg + ": not found")
		} else {
			fmt.Println("error:", err)
		}
		return;
	}
	fmt.Println(arg + " is " + path)
}

