package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirFailed(t *testing.T) {
	t.Run("Error if not exists", func(t *testing.T) {
		_, err := ReadDir("testdata/-env")
		require.Error(t, err)
	})
}

func TestReadDir(t *testing.T) {
	t.Run("Testdata expected result", func(t *testing.T) {
		result, err := ReadDir("testdata/env")
		require.NoError(t, err)
		expected := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{Value: ""},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: `"hello"`},
			"UNSET": EnvValue{NeedRemove: true, Value: ""},
		}
		require.Equal(t, expected, result)
	})
}
