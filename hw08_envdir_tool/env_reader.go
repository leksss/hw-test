package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))
	for _, file := range files {
		envName := file.Name()
		if !isValidEnvName(envName) {
			continue
		}

		envValue, err := readFileFirstLine(dir, envName)
		if err != nil {
			return nil, err
		}

		_, ok := os.LookupEnv(envName)
		if ok {
			err := os.Unsetenv(envName)
			if err != nil {
				return nil, err
			}
			delete(env, envName)
		}

		if len(envValue) > 0 {
			err = os.Setenv(envName, envValue)
			if err != nil {
				return nil, err
			}
			env[envName] = EnvValue{
				Value: envValue,
			}
		}
	}
	return env, nil
}

func readFileFirstLine(dirName, fileName string) (string, error) {
	file, err := os.Open(dirName + "/" + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return filterText(scanner.Text()), nil
}

func filterText(text string) string {
	text = strings.TrimRight(text, " \t")
	replText := bytes.ReplaceAll([]byte(text), []byte("\x00"), []byte("\n"))
	return string(replText)
}

func isValidEnvName(name string) bool {
	if name == "" {
		return false
	}
	return !strings.Contains(name, "=")
}
