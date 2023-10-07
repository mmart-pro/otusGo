package main

import (
	"bufio"
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
	dirContent, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)
	for _, entry := range dirContent {
		if entry.IsDir() || strings.Contains(entry.Name(), "=") {
			continue
		}

		val, err := readOnce(dir + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		result[entry.Name()] = EnvValue{
			NeedRemove: len(val) == 0,
			Value:      strings.ReplaceAll(strings.TrimRight(val, " "), "\x00", "\n"),
		}
	}

	return result, nil
}

// чтение первой строки из файла
func readOnce(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	} else {
		err := scanner.Err()
		return "", err
	}
}
