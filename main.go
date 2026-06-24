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
		func() {
			fmt.Print("$ ")
			input, err := reader.ReadString('\n')

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				os.Exit(1)
			}

			input = strings.TrimSpace(input)
			splitResult := handleSplit(input)
			fields := splitResult.fields
			redirectCharIndex := splitResult.redirectCharIndex
			redirectChar := splitResult.redirectChar

			if len(fields) == 0 {
				fmt.Println("Please enter your command")
				return
			}

			command := fields[0]
			var args []string
			if len(fields) >= 2 {
				args = fields[1:]
			} else {
				args = []string{}
			}

			if redirectCharIndex <= 0 { // if ">" is absent in args or it is the first arg
				runCommand(command, args...)
				return
			}

			// if ">" is present in args and is not the first arg
			fileNameIndex := redirectCharIndex + 1
			if len(fields) <= fileNameIndex {
				fmt.Fprintln(os.Stderr, "Incomplete command: Output filename not provided")
				return
			}
			fileName := fields[fileNameIndex]
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to file:", err)
				return
			}
			defer file.Close()
			var stream *os.File
			if redirectChar == "2" {
				stream = os.Stderr
			} else {
				stream = os.Stdout
			}
			originalStream := stream
			defer func() {
				if redirectChar == "2" {
					os.Stderr = originalStream
				} else {
					os.Stdout = originalStream
				}
			}()
			if redirectChar == "2" {
				os.Stderr = file
			} else {
				os.Stdout = file
			}
			args = args[:redirectCharIndex-1]
			runCommand(command, args...)
		}()
	}
}

func runCommand(command string, args ...string) {
	isBuiltIn := builtins[command]
	if isBuiltIn {
		handleBuiltInCommand(command, args...)
	} else {
		handleExternalCommand(command, args...)
	}
}

