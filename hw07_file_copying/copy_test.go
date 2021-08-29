package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	target := "/tmp/go-cp-out.txt"
	defer os.Remove(target)

	err := Copy("", target, 0, 0)
	require.Equal(t, ErrFromPathEmpty, err)

	err = Copy("/tmp/1.txt", "", 0, 0)
	require.Equal(t, ErrToPathEmpty, err)

	err = Copy(target, target, 0, 0)
	require.Equal(t, ErrFromPathToPathIdentical, err)

	err = Copy("testdata", target, 0, 0)
	require.Equal(t, ErrUnsupportedFile, err)

	err = Copy("testdata/empty.txt", target, 100, 0)
	require.Equal(t, ErrOffsetExceedsFileSize, err)

	err = Copy("testdata/not_exists", target, 0, 0)
	require.NotNil(t, err)
}

func TestWithSubDirCreation(t *testing.T) {
	dir := "/tmp/go_cp_not_exists_folder"
	target := path.Join(dir, "go-cp-out.txt")
	defer os.RemoveAll(dir)

	err := Copy("testdata/input.txt", target, 0, 0)
	require.Nil(t, err)
	_, err = os.Stat(target)
	require.Nil(t, err)
}

func TestCopy(t *testing.T) {
	source := "testdata/input.txt"
	target := "/tmp/go-cp-out.txt"
	defer os.Remove(target)

	err := Copy(source, target, 0, 0)
	require.Nil(t, err)
	requireFilesEqual(t, source, target)

	err = Copy(source, target, 0, 1000)
	require.Nil(t, err)
	requireFilesEqual(t, "testdata/out_offset0_limit1000.txt", target)

	err = Copy(source, target, 6000, 1000)
	require.Nil(t, err)
	requireFilesEqual(t, "testdata/out_offset6000_limit1000.txt", target)
}

func requireFilesEqual(t *testing.T, file1, file2 string) {
	t.Helper()

	c1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatal(err)
	}
	if len(c1) != len(c2) {
		t.Fatalf("file %q is not equals to file %q", file1, file2)
	}
	for i := 0; i < len(c1); i++ {
		if c1[i] != c2[i] {
			t.Fatalf("file %q is not equals to file %q", file1, file2)
		}
	}
}
