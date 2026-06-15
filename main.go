package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		fields := strings.Fields(input)

		if len(fields) == 0 {
			fmt.Println("Please enter your command")
			continue
		}

		command := fields[0]
		args := fields[1:]
		isBuiltIn := builtins[command]
		
		if isBuiltIn {
			executeBuiltInCommand(command, args...)
		} else {
			executeExternalCommand(command, args...)
		}
	}
}

func executeBuiltInCommand(name string, args ...string) {
	switch name {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			executeTypeCommand(strings.Join(args, " "))
		default:
			fmt.Println(name + ": command not found")
	}
}

func executeExternalCommand(name string, args ...string) {
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

