package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	t.Run("test_dir_not_exists", func(t *testing.T) {
		env, err := ReadDir("not_exists")
		assert.Nil(t, env)
		assert.Error(t, err)
	})
	t.Run("test_dir_is_empty", func(t *testing.T) {
		env, err := ReadDir("testdata")
		assert.NotNil(t, env)
		assert.Empty(t, env)
		assert.Nil(t, err)
	})
	t.Run("test_env_file_names", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		assert.NotNil(t, env)
		assert.NotEmpty(t, env)
		assert.Nil(t, err)
		// following filenames is wrong
		assert.NotContains(t, env, "2SKIP")
		assert.NotContains(t, env, "SKIP.txt")
		assert.NotContains(t, env, "SKIP=")
		// following filenames is valid
		assert.Contains(t, env, "_noskip")
		assert.Contains(t, env, "BAR")
		assert.Contains(t, env, "EMPTY")
		assert.Contains(t, env, "FOO")
		assert.Contains(t, env, "HELLO")
		assert.Contains(t, env, "NO_SKIP_09")
		assert.Contains(t, env, "UNSET")
	})

	t.Run("test_env_file_contents", func(t *testing.T) {
		env, _ := ReadDir("testdata/env")

		assert.Equal(t, "", env["_noskip"].Value)
		assert.True(t, env["_noskip"].NeedRemove)

		assert.Equal(t, "bar", env["BAR"].Value)
		assert.False(t, env["BAR"].NeedRemove)

		assert.Equal(t, "", env["EMPTY"].Value)
		assert.False(t, env["EMPTY"].NeedRemove)

		assert.Equal(t, "   foo\nwith new line", env["FOO"].Value)
		assert.False(t, env["FOO"].NeedRemove)

		assert.Equal(t, `"hello"`, env["HELLO"].Value)
		assert.False(t, env["HELLO"].NeedRemove)

		assert.Equal(t, "   no_skip_09", env["NO_SKIP_09"].Value)
		assert.False(t, env["NO_SKIP_09"].NeedRemove)

		assert.Equal(t, "", env["UNSET"].Value)
		assert.True(t, env["UNSET"].NeedRemove)
	})
}
