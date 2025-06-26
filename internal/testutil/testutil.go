package testutil

import (
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/constants"
	_ "github.com/mattn/go-sqlite3"
)

type TestUtil struct {
	t *testing.T
}

func (tu TestUtil) AssertEqual(expected, actual interface{}) {
	tu.t.Helper()
	if expected != actual {
		tu.t.Errorf("expected %v, instead got %v", expected, actual)
	}
}

func (tu TestUtil) AssertTrue(value bool) {
	tu.t.Helper()
	if !value {
		tu.t.Error("expected value to be true, instead it was false")
	}
}

func (tu TestUtil) AssertFalse(value bool) {
	tu.t.Helper()
	if value {
		tu.t.Error("expected value to be false, instead it was true")
	}
}

func (tu TestUtil) AssertNotNil(value interface{}) {
	tu.t.Helper()
	if isNil(value) {
		tu.t.Errorf("expected non-nil value, it was nil: %v", value)
	}
}

func (tu TestUtil) AssertNil(value interface{}) {
	tu.t.Helper()
	if !isNil(value) {
		tu.t.Errorf("expected nil value, it was not nil: %v", value)
	}
}

func (tu TestUtil) AssertErrorNotNil(err error) {
	tu.t.Helper()
	if err == nil {
		tu.t.Error("expected error to not be nil, instead it was nil")
	}
}

func (tu TestUtil) AssertErrorNil(err error) {
	tu.t.Helper()
	if err != nil {
		tu.t.Error("expected error to be nil, instead it was not nil")
	}
}

func (tu TestUtil) ShouldPanic() {
	tu.t.Helper()
	if r := recover(); r == nil {
		tu.t.Errorf("expected panic, but function did not panic")
	}
}

func (tu TestUtil) CreateTestUploadsDir() {
	os.MkdirAll(os.Getenv("TEST_IMAGE_STORAGE_PATH"), os.ModePerm)
}

func (tu TestUtil) CleanTestUploads() {
	tu.t.Cleanup(func() {
		deleteTestDir()
	})
}

func NewTestUtil(t *testing.T) TestUtil {
	return TestUtil{t}
}

func setTestEnvironment() {
	os.Setenv("ENV", constants.TEST_ENV)
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("failed to open in memory db:", err)
	}

	schema, err := os.ReadFile("../../sql/schema.sql")
	if err != nil {
		t.Fatal("failed to read schema file:", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatal("failed to execute schema sql", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatal("could not enable foreign keys:", err)
	}

	return db
}

func WithTestDB(t *testing.T, testFunc func(db *sql.DB)) {
	db := setupTestDB(t)
	defer db.Close()
	testFunc(db)
}

func WithTestData(t *testing.T, testFunc func(db *sql.DB, timestamp time.Time)) {
	db := setupTestDB(t)
	defer db.Close()
	seedDataQuery, err := os.ReadFile("../../sql/seed-test-data.sql")
	if err != nil {
		t.Fatal("error reading test seed data file: ", err)
	}
	timestamp := time.Now()
	_, err = db.Exec(string(seedDataQuery))
	if err != nil {
		t.Fatal("error executing test seed data query: ", err)
	}

	testFunc(db, timestamp)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}

func QueryUserPassword(userID int, db *sql.DB) string {
	query := `SELECT password FROM User WHERE id = $1;`

	var password string
	err := db.QueryRow(query, userID).Scan(&password)
	if err != nil {
		panic(err)
	}

	return password
}

func deleteTestDir() {
	os.RemoveAll(os.Getenv("TEST_IMAGE_STORAGE_PATH"))
}

func cleanDir(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range entries {
		err = os.RemoveAll(filepath.Join(path, file.Name()))
		if err != nil {
			panic(err)
		}
	}
}

func GetTestUploads() []fs.DirEntry {
	files, err := os.ReadDir(os.Getenv("TEST_IMAGE_STORAGE_PATH"))
	if err != nil {
		panic(err)
	}

	return files
}
