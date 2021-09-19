package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var envFileMask = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*$`)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	vars := make(Environment)
	for _, file := range files {
		ename := file.Name()
		if !envFileMask.MatchString(ename) {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			continue
		}

		if info.Size() == 0 {
			vars[ename] = EnvValue{NeedRemove: true}
			continue
		}

		val, err := getEnvVal(filepath.Join(dir, ename))
		if err != nil {
			return nil, err
		}
		vars[ename] = EnvValue{Value: val}
	}
	return vars, nil
}

func getEnvVal(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := bytes.ReplaceAll(scanner.Bytes(), []byte{0}, []byte("\n"))
		val := string(b)
		val = strings.TrimRight(val, "\t\n\v\f\r ")
		return val, nil
	}

	return "", scanner.Err()
}
