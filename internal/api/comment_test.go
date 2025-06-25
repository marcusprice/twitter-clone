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

	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestCreateCommentContentOnly(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		defer tu.CleanTestUploads()
		handler := RegisterHandlers(db)

		testUser := createTestUser(db)
		testUser.Login()
		token, _ := GenerateJWT(testUser.ID())

		formValues := make(map[string]string)
		formValues["content"] = "Cats are awesome"
		formValues["postID"] = "1"
		requestBody, contentType, _ := util.GenerateMultipartForm(formValues)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comment/create", requestBody)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", contentType)
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		var commentPayload CommentPayload
		json.Unmarshal(res.Body.Bytes(), &commentPayload)

		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("Cats are awesome", commentPayload.Content)
		tu.AssertEqual("", commentPayload.Image)
		tu.AssertEqual(0, commentPayload.LikeCount)
		tu.AssertEqual(0, commentPayload.RetweetCount)
		tu.AssertEqual(0, commentPayload.BookmarkCount)
		tu.AssertEqual(0, commentPayload.Impressions)
		tu.AssertEqual("esteban", commentPayload.Author.Username)
		tu.AssertEqual("Bubba", commentPayload.Author.DisplayName)
		tu.AssertEqual("", commentPayload.Author.Avatar)
		tu.AssertTrue(commentPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(commentPayload.CreatedAt.Before(afterRequest))
	})
}

func TestCreateCommentImageOnly(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
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

		postIDField, _ := writer.CreateFormField("postID")
		io.Copy(postIDField, strings.NewReader("1"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comment/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		uploads := testutil.GetTestUploads()
		fileWritten := len(uploads) == 1
		uploadedFileName := uploads[0].Name()
		var commentPayload CommentPayload
		json.Unmarshal(res.Body.Bytes(), &commentPayload)

		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("", commentPayload.Content)
		tu.AssertEqual(0, commentPayload.LikeCount)
		tu.AssertEqual(0, commentPayload.RetweetCount)
		tu.AssertEqual(0, commentPayload.BookmarkCount)
		tu.AssertEqual(0, commentPayload.Impressions)
		tu.AssertEqual("esteban", commentPayload.Author.Username)
		tu.AssertEqual("Bubba", commentPayload.Author.DisplayName)
		tu.AssertEqual("", commentPayload.Author.Avatar)
		tu.AssertEqual(commentPayload.Image, uploadedFileName)
		tu.AssertTrue(fileWritten)
		tu.AssertTrue(strings.Contains(uploadedFileName, "meme.png"))
		tu.AssertTrue(strings.Contains(commentPayload.Image, "meme.png"))
		tu.AssertTrue(commentPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(commentPayload.CreatedAt.Before(afterRequest))
	})
}

func TestCreateCommentContentAndImage(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
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
		postIDField, _ := writer.CreateFormField("postID")
		io.Copy(postIDField, strings.NewReader("1"))
		imgField, _ := writer.CreateFormFile("image", "sunset.jpeg")
		io.Copy(imgField, strings.NewReader(imgData))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/comment/create", &b)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(res, req)
		afterRequest := time.Now().UTC().Add(time.Minute)

		uploads := testutil.GetTestUploads()
		fileWritten := len(uploads) == 1
		uploadedFileName := uploads[0].Name()
		var commentPayload CommentPayload
		json.Unmarshal(res.Body.Bytes(), &commentPayload)
		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual("Check out this gorgeous sunset", commentPayload.Content)
		tu.AssertEqual(0, commentPayload.LikeCount)
		tu.AssertEqual(0, commentPayload.RetweetCount)
		tu.AssertEqual(0, commentPayload.BookmarkCount)
		tu.AssertEqual(0, commentPayload.Impressions)
		tu.AssertEqual("esteban", commentPayload.Author.Username)
		tu.AssertEqual("Bubba", commentPayload.Author.DisplayName)
		tu.AssertEqual("", commentPayload.Author.Avatar)
		tu.AssertEqual(commentPayload.Image, uploadedFileName)
		tu.AssertTrue(fileWritten)
		tu.AssertTrue(strings.Contains(uploadedFileName, "sunset.jpeg"))
		tu.AssertTrue(strings.Contains(commentPayload.Image, "sunset.jpeg"))
		tu.AssertTrue(commentPayload.CreatedAt.After(beforeRequest))
		tu.AssertTrue(commentPayload.CreatedAt.Before(afterRequest))
	})
}
