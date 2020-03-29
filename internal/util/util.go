package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// SplisLibPath returns the path to the splis standard library which is
// located at '$(SPLIS_HOME)/lib'.
func SplisLibPath() string {
	return path.Join(SplisHomePath(), "lib")
}

// SplisHomePath reads the environment variable SPLIS_HOME.
func SplisHomePath() string {
	if path, ok := os.LookupEnv("SPLIS_HOME"); ok {
		return path
	}
	panic("Unable to read environment variable [SPLIS_HOME].")
}

func FilePathWithoutExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func ReadStatic(file *os.File) []byte {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return b
}

func WriteStatic(code []byte, file *os.File) {
	w := bufio.NewWriter(file)
	_, err := w.Write(code)
	if err != nil {
		panic(err)
	}
	w.Flush()
}
