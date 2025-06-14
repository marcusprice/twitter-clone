package util

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marcusprice/twitter-clone/internal/constants"
)

func ParseTime(timestamp string) time.Time {
	parsedTime, err := time.Parse(constants.TIME_LAYOUT, timestamp)
	if err != nil {
		// likely a null/empty value
		return time.Time{}
	}

	return parsedTime
}

func ProjectRoot() (string, error) {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

func LoadEnvVariables() {
	root, err := ProjectRoot()
	if err != nil {
		panic(err)
	}
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
		if len(line) == 0 {
			continue
		}

		if string(line[0]) == "#" {
			continue
		}

		keyValueSplit := strings.SplitN(line, "=", 2)
		if len(keyValueSplit) != 2 {
			continue
		}

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
