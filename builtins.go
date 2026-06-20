package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

func handleBuiltInCommand(name string, args ...string) {
	switch name {
		case "exit":
			handleExit()
		case "echo":
			handleEcho(strings.Join(args, " "))
		case "type":
			handleType(strings.Join(args, " "))
		case "pwd":
			handlePwd(args...)
		case "cd":
			handleCd(args...)
		default:
			fmt.Println(name + ": command not found")
	}
}

func handleExit() {
	os.Exit(0)
}

func handleEcho(arg string) {
	fmt.Println(arg)
}

func handleType(arg string) {
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

func handlePwd(args ...string) {
	if len(args) > 0 {
		fmt.Println("pwd: too many arguments")
		return
	}
	wd, err := os.Getwd();
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(wd)
}

func handleCd(args ...string) {
	if len(args) == 0 {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Println("cd: HOME not set")
			return
		}
		os.Chdir(home)
		return
	} 
	path := args[0]
	if err := os.Chdir(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("cd: " + path + ": No such file or directory")
		} else {
			fmt.Println("error: ", err)
		}
	}
} 