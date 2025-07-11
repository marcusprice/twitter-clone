package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestCreatePostContentOnly(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)

		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		contenfPart, _ := writer.CreateFormField("content")
		io.Copy(contenfPart, strings.NewReader("Cats are awesome"))
		writer.Close()

		req := httptest.NewRequest(
			http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		var postPayload PostPayload
		json.Unmarshal(res.Body.Bytes(), &postPayload)

		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("Cats are awesome", postPayload.Content)
		tu.AssertEqual("", postPayload.Image)
		tu.AssertEqual(0, postPayload.LikeCount)
		tu.AssertEqual(0, postPayload.RetweetCount)
		tu.AssertEqual(0, postPayload.BookmarkCount)
		tu.AssertEqual(0, postPayload.Impressions)
		tu.AssertEqual("esteban", postPayload.Author.Username)
		tu.AssertEqual("Bubba", postPayload.Author.DisplayName)
		tu.AssertEqual("", postPayload.Author.Avatar)
		tu.AssertTrue(postPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(postPayload.CreatedAt.Before(afterRequest))
	})
}

func TestCreatePostImageOnly(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)

		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		imgField, _ := writer.CreateFormFile("image", "meme.png")
		imgData := generateLargeString(5)
		io.Copy(imgField, strings.NewReader(imgData))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		uploads := testutil.GetTestUploads()
		fileWritten := len(uploads) == 1
		uploadedFileName := uploads[0].Name()
		var postPayload PostPayload
		json.Unmarshal(res.Body.Bytes(), &postPayload)

		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("", postPayload.Content)
		tu.AssertEqual(0, postPayload.LikeCount)
		tu.AssertEqual(0, postPayload.RetweetCount)
		tu.AssertEqual(0, postPayload.BookmarkCount)
		tu.AssertEqual(0, postPayload.Impressions)
		tu.AssertEqual("esteban", postPayload.Author.Username)
		tu.AssertEqual("Bubba", postPayload.Author.DisplayName)
		tu.AssertEqual("", postPayload.Author.Avatar)
		tu.AssertEqual(postPayload.Image, getUploadPath(uploadedFileName))
		tu.AssertTrue(fileWritten)
		tu.AssertTrue(strings.Contains(uploadedFileName, "meme.png"))
		tu.AssertTrue(strings.Contains(postPayload.Image, "meme.png"))
		tu.AssertTrue(postPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(postPayload.CreatedAt.Before(afterRequest))
	})
}

func TestCreatePostContentAndImage(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)
		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		var b bytes.Buffer
		imgData := generateLargeString(5)
		writer := multipart.NewWriter(&b)
		contentField, _ := writer.CreateFormField("content")
		io.Copy(contentField, strings.NewReader("Check out this gorgeous sunset"))
		imgField, _ := writer.CreateFormFile("image", "sunset.jpeg")
		io.Copy(imgField, strings.NewReader(imgData))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		uploads := testutil.GetTestUploads()
		fileWritten := len(uploads) == 1
		uploadedFileName := uploads[0].Name()
		var postPayload PostPayload
		json.Unmarshal(res.Body.Bytes(), &postPayload)
		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("Check out this gorgeous sunset", postPayload.Content)
		tu.AssertEqual(0, postPayload.LikeCount)
		tu.AssertEqual(0, postPayload.RetweetCount)
		tu.AssertEqual(0, postPayload.BookmarkCount)
		tu.AssertEqual(0, postPayload.Impressions)
		tu.AssertEqual("esteban", postPayload.Author.Username)
		tu.AssertEqual("Bubba", postPayload.Author.DisplayName)
		tu.AssertEqual("", postPayload.Author.Avatar)
		tu.AssertEqual(postPayload.Image, getUploadPath(uploadedFileName))
		tu.AssertTrue(fileWritten)
		tu.AssertTrue(strings.Contains(uploadedFileName, "sunset.jpeg"))
		tu.AssertTrue(strings.Contains(postPayload.Image, "sunset.jpeg"))
		tu.AssertTrue(postPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(postPayload.CreatedAt.Before(afterRequest))
	})
}

func TestCreatePostInvalidFileType(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)
		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		fileString := generateLargeString(5)
		imgPart, _ := writer.CreateFormFile("image", "video.mp4")
		io.Copy(imgPart, strings.NewReader(fileString))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded := len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusUnsupportedMediaType, res.Code)
		tu.AssertTrue(fileNotUploaded)

		writer = multipart.NewWriter(&b)
		imgPart, _ = writer.CreateFormFile("image", "video.flac")
		io.Copy(imgPart, strings.NewReader(fileString))
		writer.Close()

		req = httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded = len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusUnsupportedMediaType, res.Code)
		tu.AssertTrue(fileNotUploaded)
	})
}

func TestCreatePostNoContentOrImage(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)
		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		content, _ := writer.CreateFormField("content")
		image, _ := writer.CreateFormField("image")
		io.Copy(content, strings.NewReader(""))
		io.Copy(image, strings.NewReader(""))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded := len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusBadRequest, res.Code)
		tu.AssertTrue(fileNotUploaded)

		b.Reset()
		writer = multipart.NewWriter(&b)
		content, _ = writer.CreateFormField("content")
		image, _ = writer.CreateFormField("image")
		io.Copy(content, strings.NewReader(""))
		io.Copy(image, strings.NewReader("data but not image"))
		writer.Close()

		req = httptest.NewRequest(http.MethodPost, "/api/v1/post/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded = len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusBadRequest, res.Code)
		tu.AssertTrue(fileNotUploaded)
	})
}

func TestCreatePostUploadSizeTooLarge(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())
		tu.CreateTestUploadsDir()
		defer tu.CleanTestUploads()

		handler := RegisterHandlers(db)
		b, contentType := createLargeImgMultipartFormBody(0)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", contentType)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded := len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusRequestEntityTooLarge, res.Code)
		tu.AssertTrue(fileNotUploaded)

		b, contentType = createLargeImgMultipartFormBody(0.5)

		req = httptest.NewRequest(http.MethodPost, "/api/v1/post/create", b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", contentType)
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded = len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusRequestEntityTooLarge, res.Code)
		tu.AssertTrue(fileNotUploaded)

		b, contentType = createLargeImgMultipartFormBody(5)

		req = httptest.NewRequest(http.MethodPost, "/api/v1/post/create", b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", contentType)
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		fileNotUploaded = len(testutil.GetTestUploads()) == 0

		tu.AssertEqual(http.StatusRequestEntityTooLarge, res.Code)
		tu.AssertTrue(fileNotUploaded)
	})
}

func TestCreatePostUnauthorized(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		noAuthHeaderReq := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", nil)
		noAuthHeaderRes := httptest.NewRecorder()
		handler.ServeHTTP(noAuthHeaderRes, noAuthHeaderReq)

		headerNoTokenReq := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", nil)
		headerNoTokenReq.Header.Set("Authorization", "Bearer ")
		headerNoTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(headerNoTokenRes, headerNoTokenReq)

		headerWrongKeywordReq := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", nil)
		headerWrongKeywordReq.Header.Set("Authorization", "Esteban ")
		headerWrongKeywordRes := httptest.NewRecorder()
		handler.ServeHTTP(headerWrongKeywordRes, headerWrongKeywordReq)

		badToken := generateBadToken()
		badTokenReq := httptest.NewRequest(http.MethodPost, "/api/v1/post/create", nil)
		badTokenReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", badToken))
		badTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(badTokenRes, badTokenReq)

		tu.AssertEqual(http.StatusUnauthorized, noAuthHeaderRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerNoTokenRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerWrongKeywordRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, badTokenRes.Code)
	})
}

func TestCreatePostWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/post/create", nil)
		getRes := httptest.NewRecorder()
		putReq := httptest.NewRequest(http.MethodPut, "/api/v1/post/create", nil)
		putRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/post/create", nil)
		patchRes := httptest.NewRecorder()
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/post/create", nil)
		deleteRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/post/create", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/post/create", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/post/create", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/post/create", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(getRes, getReq)
		handler.ServeHTTP(putRes, putReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(deleteRes, deleteReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}

func generateBadToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 42069,
	})

	out, err := token.SignedString([]byte("esteban is a very good cat"))
	if err != nil {
		panic(err)
	}
	return out
}

func generateLargeString(sizeInMB float64) string {
	size := int(sizeInMB * 1024 * 1024)
	b := make([]byte, size)
	for i := range b {
		b[i] = 'A'
	}

	return string(b)
}

func convertBytesToMB(bytes int64) float64 {
	return float64(bytes / 1024 / 1024)
}

func createLargeImgMultipartFormBody(mbOver float64) (*bytes.Buffer, string) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	mb := convertBytesToMB(MAX_POST_UPLOAD_BYTES) + mbOver
	largerString := generateLargeString(mb)
	imgPart, err := writer.CreateFormFile("image", "bigger.jpg")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(imgPart, strings.NewReader(largerString))
	if err != nil {
		panic(err)
	}

	writer.Close()

	return &b, writer.FormDataContentType()
}

func createTestUser(db *sql.DB) *controller.User {
	user := controller.NewUserController(db)
	userInput := dtypes.UserInput{
		Username:    "esteban",
		Email:       "estecat42069@yahoo.com",
		Password:    "password",
		DisplayName: "Bubba",
	}
	user.Set(nil, userInput)
	err := user.Create("password")
	if err != nil {
		panic(err)
	}
	return user
}
