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

		fields := strings.Fields(command)
		executeCommand(fields[0], fields[1:]...)
	}
}

func executeCommand(name string, args ...string) {
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

