package main

import (
	"errors"
	"log"
	"os"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrExitCodeReceived   = errors.New("exit code received")
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal(ErrNotEnoughArguments.Error())
	}

	env, err := ReadDir(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	returnCode := RunCmd(args[1:], env)
	if returnCode != 0 {
		log.Fatal(ErrExitCodeReceived.Error())
	}
}
