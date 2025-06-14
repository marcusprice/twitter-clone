package api

import (
	"fmt"
	"os"

	"github.com/marcusprice/twitter-clone/internal/util"
)

func init() {
	util.LoadEnvVariables()
	testUploadsPath := os.Getenv("TEST_IMAGE_STORAGE_PATH")
	if testUploadsPath == "" {
		panic(fmt.Errorf("need TEST_IMAGE_STORAGE_PATH env variable to be set"))
	}

	os.Setenv("IMAGE_STORAGE_PATH", testUploadsPath)
}
