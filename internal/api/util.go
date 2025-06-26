package api

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/marcusprice/twitter-clone/internal/util"
)

var ACCEPTED_IMAGE_FORMATS = []string{"jpg", "jpeg", "gif", "png", "webp"}

type InvalidFileTypeError struct {
	filename string
}

func (w InvalidFileTypeError) Error() string {
	return fmt.Sprintf("%s is not an accepted file type", w.filename)
}

func getUploadPath(fileName string) string {
	return fmt.Sprintf(
		"http://%s:%s/%s/%s",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		UPLOADS_PREFIX,
		fileName,
	)
}

func handleImageUpload(file multipart.File, header *multipart.FileHeader) (filename string, err error) {
	if !validImageFormat(header.Filename) {
		return "", InvalidFileTypeError{header.Filename}
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	path, err := getImageStoragePath()
	if err != nil {
		return "", err
	}

	filename, err = generateUniqueFilename(header.Filename)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s", path, filename), fileBytes, 0755)
	if err != nil {
		return "", err
	}

	return filename, err
}

func getImageStoragePath() (string, error) {
	imageStoragePath := os.Getenv("IMAGE_STORAGE_PATH")

	if imageStoragePath == "" {
		projectRoot, err := util.ProjectRoot()
		if err != nil {
			return "", err
		}
		imageStoragePath = fmt.Sprintf("%s/uploads", projectRoot)
	}

	err := os.MkdirAll(imageStoragePath, 0755)
	if err != nil {
		return "", err
	}

	return imageStoragePath, nil
}

func generateUniqueFilename(filename string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	newFilename := santizeFilename(fmt.Sprintf("%s-%s", id, filename))

	return newFilename, nil
}

func getMaxUploadMemory() int64 {
	var defaultMaxUploadMemory int64 = 2097152
	maxUploadMemory, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_MEMORY"), 10, 64)
	if err != nil {
		maxUploadMemory = defaultMaxUploadMemory
	}

	return maxUploadMemory
}

func validImageFormat(filename string) bool {
	fileType := strings.Split(filename, ".")

	if len(fileType) == 0 {
		return false
	}

	return slices.Contains(
		ACCEPTED_IMAGE_FORMATS,
		strings.ToLower(fileType[len(fileType)-1]))
}

func requestBodyTooLarge(err error) bool {
	return (errors.Is(err, http.ErrBodyReadAfterClose) ||
		strings.Contains(err.Error(), "http: request body too large"))
}

func santizeFilename(str string) (sanitized string) {
	return strings.Map(func(char rune) rune {
		if unicode.IsSpace(char) {
			return '-'
		}

		return char
	}, str)
}
