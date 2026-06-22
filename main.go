package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd": true,
	"cd": true,
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
		fields := handleSplit(input)

		if len(fields) == 0 {
			fmt.Println("Please enter your command")
			continue
		}

		command := fields[0]
		var args []string
		if len(fields) >= 2 {
			args = fields[1:]
		} else {
			args = []string{}
		}
		isBuiltIn := builtins[command]
		
		if isBuiltIn {
			handleBuiltInCommand(command, args...)
		} else {
			handleExternalCommand(command, args...)
		}
	}
}

