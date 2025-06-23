package util

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
)

// TODO: add support for files
func GenerateMultipartForm(fields map[string]string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	multipartWriter := multipart.NewWriter(&b)
	for key, value := range fields {
		field, err := multipartWriter.CreateFormField(key)
		if err != nil {
			return &bytes.Buffer{}, err
		}

		_, err = io.Copy(field, strings.NewReader(value))
		if err != nil {
			return &bytes.Buffer{}, err
		}
	}
	multipartWriter.Close()

	return &b, nil
}
