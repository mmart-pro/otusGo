package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	run := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	run.Stdin = os.Stdin
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	run.Env = mergeEnv(env)

	if err := run.Run(); err != nil {
		fmt.Println(err)
	}
	return run.ProcessState.ExitCode()
}

// мерж массива переменных окружения с переданными в env
func mergeEnv(env Environment) []string {
	src := os.Environ()
	result := make([]string, 0, len(src)+len(env))
	for _, v := range src {
		vName := strings.SplitN(v, "=", 2)[0]
		if _, ok := env[vName]; !ok {
			result = append(result, v)
		}
	}
	// добавить новые значения
	for key, value := range env {
		if !value.NeedRemove {
			result = append(result, key+"="+value.Value)
		}
	}
	return result
}
