package testutil

import (
	"database/sql"
	"os"
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
