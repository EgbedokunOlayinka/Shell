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
			handleEcho(args...)
		case "type":
			handleType(args...)
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

func handleEcho(args ...string) {
	arg := strings.Join(args, " ")
	fmt.Println(arg)
}

func handleType(args ...string) {
	if len(args) == 0 {
		return
	}
	arg := strings.Join(args, " ")
	if builtins[arg] {
		fmt.Println(arg + " is a shell builtin")
		return
	} 
	path, err := exec.LookPath(arg)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Fprintln(os.Stderr, arg + ": not found")
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		return;
	}
	fmt.Println(arg + " is " + path)
}

func handlePwd(args ...string) {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "pwd: too many arguments")
		return
	}
	wd, err := os.Getwd();
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}
	fmt.Println(wd)
}

func handleCd(args ...string) {
	if len(args) == 0 {
		return
	} 
	path := args[0]
	if(path == "~") {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Fprintln(os.Stderr, "cd: HOME not set")
			return
		}
		os.Chdir(home)
		return
	}
	if err := os.Chdir(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "cd: " + path + ": No such file or directory")
		} else {
			fmt.Fprintln(os.Stderr, "error: ", err)
		}
	}
} 