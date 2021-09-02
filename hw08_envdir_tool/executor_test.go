package main

import (
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/example/stringutil"
)

func TestRunCmd(t *testing.T) {
	t.Run("test_stdin_stdout", func(t *testing.T) {
		restoreStdin := mockStdin(t)
		defer restoreStdin()

		restoreStdout := mockStdout(t)
		defer restoreStdout()

		in := "ABC"
		cmd := "rev"
		out := stringutil.Reverse(in)

		writeStdin(t, in)

		code := RunCmd([]string{cmd}, make(Environment))

		result := readStdout(t)

		require.Equal(t, 0, code)
		require.Equal(t, out, result)
	})

	t.Run("test_env_var", func(t *testing.T) {
		restoreStdout := mockStdout(t)
		defer restoreStdout()

		envs := make(Environment)
		envs["HELLO_WORLD"] = EnvValue{Value: "hello, world!"}

		cmd := []string{"printenv", "HELLO_WORLD"}

		code := RunCmd(cmd, envs)

		result := strings.TrimRight(readStdout(t), "\n")

		require.Equal(t, 0, code)
		require.Equal(t, envs["HELLO_WORLD"].Value, result)
	})
}

func TestReplaceEnvs(t *testing.T) {
	osEnv := []string{
		"ORIGINAL=original",
		"UNSET=unset",
		"REPLACED=replaced",
	}

	newEnv := Environment{
		"UNSET":    EnvValue{NeedRemove: true},
		"REPLACED": EnvValue{Value: "new_value"},
		"NEW_VAR":  EnvValue{Value: "new_var"},
	}

	expected := []string{
		"ORIGINAL=original",
		"REPLACED=new_value",
		"NEW_VAR=new_var",
	}

	result := replaceEnvs(newEnv, osEnv)

	sort.Strings(result)
	sort.Strings(expected)
	require.Equal(t, expected, result)
}

func mockStdout(t *testing.T) (restore func()) {
	t.Helper()

	backup := os.Stdout
	mock, err := os.CreateTemp("testdata", "std-tmp")
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = mock

	return func() {
		os.Stdout = backup
		mock.Close()
		os.Remove(mock.Name())
	}
}

func mockStdin(t *testing.T) (restore func()) {
	t.Helper()

	backup := os.Stdin
	mock, err := os.CreateTemp("testdata", "std-tmp")
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = mock

	return func() {
		os.Stdin = backup
		mock.Close()
		os.Remove(mock.Name())
	}
}

func writeStdin(t *testing.T, content string) {
	t.Helper()

	if _, err := os.Stdin.WriteString(content); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stdin.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
}

func readStdout(t *testing.T) string {
	t.Helper()

	if _, err := os.Stdout.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	result, err := io.ReadAll(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

	return string(result)
}
