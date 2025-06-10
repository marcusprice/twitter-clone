package testutil

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

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

	return db
}

func WithTestDB(t *testing.T, testFunc func(db *sql.DB)) {
	db := setupTestDB(t)
	defer db.Close()
	testFunc(db)
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
