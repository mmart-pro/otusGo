package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type std struct {
	out string
	err string
}

func TestStdOutErr(t *testing.T) {
	t.Run("Test error code and stderr", func(t *testing.T) {
		exitCode, std := stdWrapper(t, []string{"rm", "."}, Environment{})
		require.Equal(t, 1, exitCode)
		require.Equal(t, `rm: "." and ".." may not be removed`+"\n", std.err)
	})
	t.Run("Test stdout", func(t *testing.T) {
		exitCode, std := stdWrapper(t, []string{"echo", "test"}, Environment{})
		require.Equal(t, 0, exitCode)
		require.Equal(t, "test\n", std.out)
	})
}

func TestEnv(t *testing.T) {
	t.Run("Replace env", func(t *testing.T) {
		env := Environment{"HOME": EnvValue{Value: "REPLACED"}}
		exitCode, std := stdWrapper(t, []string{"/bin/bash", "-c", "echo $HOME"}, env)
		require.Equal(t, 0, exitCode)
		require.Equal(t, "REPLACED\n", std.out)
	})
	t.Run("Insert env", func(t *testing.T) {
		env := Environment{"BAR": EnvValue{Value: "INSERTED"}}
		exitCode, std := stdWrapper(t, []string{"/bin/bash", "-c", "echo $BAR"}, env)
		require.Equal(t, 0, exitCode)
		require.Equal(t, "INSERTED\n", std.out)
	})
	t.Run("Empty env", func(t *testing.T) {
		env := Environment{"USER": EnvValue{}}
		exitCode, std := stdWrapper(t, []string{"/bin/bash", "-c", "echo $USER"}, env)
		require.Equal(t, 0, exitCode)
		require.Equal(t, "\n", std.out)
	})
	t.Run("Complex", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err, "Ошибка чтения каталога ./testdata/env")
		_, std := stdWrapper(t, []string{"/bin/bash", "-c", "echo $BAR"}, env)
		require.Equal(t, "bar\n", std.out)
		_, std = stdWrapper(t, []string{"/bin/bash", "-c", "echo $EMPTY"}, env)
		require.Equal(t, "\n", std.out)
		_, std = stdWrapper(t, []string{"/bin/bash", "-c", "echo $FOO"}, env)
		require.Equal(t, "foo with new line\n", std.out)
		_, std = stdWrapper(t, []string{"/bin/bash", "-c", "echo $HELLO"}, env)
		require.Equal(t, "\"hello\"\n", std.out)
		_, std = stdWrapper(t, []string{"/bin/bash", "-c", "echo $UNSET"}, env)
		require.Equal(t, "\n", std.out)
	})
}

// перехват stdout и stderr
func stdWrapper(t *testing.T, cmd []string, env Environment) (int, std) {
	t.Helper()

	oldStdout := os.Stdout
	tempOut, err := os.CreateTemp("", "hw08_*.out")
	require.NoError(t, err, "Не удалось создать временный файл для stdout")
	defer os.Remove(tempOut.Name())
	os.Stdout = tempOut

	oldStderr := os.Stderr
	tempErr, err := os.CreateTemp("", "hw08_*.err")
	require.NoError(t, err, "Не удалось создать временный файл для stderr")
	defer os.Remove(tempErr.Name())
	os.Stderr = tempErr

	defer func() {
		os.Stderr = oldStderr
		os.Stdout = oldStdout
	}()

	exitCode := RunCmd(cmd, env)
	outBuf, err := os.ReadFile(tempOut.Name())
	require.NoError(t, err, "Ошибка чтения stdout")
	errBuf, err := os.ReadFile(tempErr.Name())
	require.NoError(t, err, "Ошибка чтения stderr")
	return exitCode, std{out: string(outBuf), err: string(errBuf)}
}
