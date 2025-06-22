package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestTimelineWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		postReq := httptest.NewRequest(http.MethodPost, "/api/v1/timeline", nil)
		postRes := httptest.NewRecorder()
		putReq := httptest.NewRequest(http.MethodPut, "/api/v1/timeline", nil)
		putRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/timeline", nil)
		patchRes := httptest.NewRecorder()
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/timeline", nil)
		deleteRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/timeline", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/timeline", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/timeline", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/timeline", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(postRes, postReq)
		handler.ServeHTTP(putRes, putReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(deleteRes, deleteReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, postRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}
