package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("test exec run command", func(t *testing.T) {
		params := []string{
			"testdata/echo.sh",
			"arg1=1",
			"arg2=2",
		}
		env := Environment{
			"arg1": EnvValue{"1", false},
			"arg2": EnvValue{"2", false},
		}
		code := RunCmd(params, env)
		assert.Equal(t, 0, code)
	})
}
