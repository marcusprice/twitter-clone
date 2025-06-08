package util

import (
	"os"
	"slices"

	"github.com/marcusprice/twitter-clone/internal/constants"
)

func InDevContext() bool {
	env := os.Getenv("ENV")
	devEnvs := []string{constants.DEV_ENV, constants.TEST_ENV}
	return slices.Contains(devEnvs, env)
}
