package main

import (
	"fmt"
	"os"

	_jvm "github.com/Stolkerve/go-jvm/jvm"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "java class file expected")
		os.Exit(1)
	}

	jvm, err := _jvm.NewJvm(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_jvm.RunJvm(jvm)
}

// type Constant
