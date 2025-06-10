package util

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcusprice/twitter-clone/internal/constants"
)

func projectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("project root not found")
		}
		dir = parent
	}
}

func LoadEnvVariables() {
	root := projectRoot()
	envFile, err := os.Open(fmt.Sprintf(`%s/%s`, root, ".env"))
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(envFile)
	if err != nil {
		panic(err)
	}

	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")
	linesWithoutLastEmptyLine := lines[:len(lines)-1]
	for _, line := range linesWithoutLastEmptyLine {
		keyValueSplit := strings.SplitN(line, "=", 2)
		err := os.Setenv(keyValueSplit[0], keyValueSplit[1])
		if err != nil {
			panic(err)
		}
	}
}

func InDevContext() bool {
	env := os.Getenv("ENV")
	return env == constants.DEV_ENV && flag.Lookup("test.v") == nil
}
