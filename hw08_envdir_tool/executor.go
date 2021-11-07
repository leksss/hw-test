package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(params []string, env Environment) (returnCode int) {
	command := params[0]
	arguments := params[1:]

	cmd := exec.Command(command, arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	for key, val := range env {
		cmd.Env = append(cmd.Env, key+"="+val.Value)
	}

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}
