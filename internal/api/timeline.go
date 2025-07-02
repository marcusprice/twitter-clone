package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/marcusprice/twitter-clone/internal/controller"
)

const MAX_LIMIT = 1
const MIN_LIMIT = 40

type TimelineAPI struct {
	timeline *controller.Timeline
}

func (timelineAPI *TimelineAPI) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	values := r.URL.Query()
	limitParam := values.Get("limit")
	offsetParam := values.Get("offset")
	viewParam := values.Get("view")

	limit, offset, err := parseLimitAndOffset(limitParam, offsetParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !slices.Contains(controller.TIMELINE_VIEWS, controller.TimelineView(viewParam)) {
		http.Error(w, BadRequest, http.StatusBadRequest)
		return
	}

	view := controller.TimelineView(viewParam)
	timeline := timelineAPI.timeline.Set(userID, view)
	posts, postsRemaining, err := timeline.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	var postPayloads []PostPayload
	for _, post := range posts {
		postPayload := generatePostPayload(post)
		postPayloads = append(postPayloads, postPayload)
	}

	timelinePayload := TimelinePayload{
		Posts:          postPayloads,
		HasMore:        postsRemaining > 0,
		PostsRemaining: postsRemaining,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timelinePayload)
}

func parseLimitAndOffset(limitParam, offsetParam string) (limit, offset int, err error) {
	limit, limitErr := strconv.Atoi(limitParam)
	if limitParam == "" || limitErr != nil {
		return -1, -1, errors.New("Bad limit value")
	}

	offset, offsetErr := strconv.Atoi(offsetParam)
	if offsetParam == "" || offsetErr != nil {
		return -1, -1, errors.New("Bad offset value")
	}

	if limit < MAX_LIMIT {
		return -1, -1, fmt.Errorf("Too large of a limit, max limit: %d", MAX_LIMIT)
	}

	if limit > MIN_LIMIT {
		return -1, -1, fmt.Errorf("Too small of a limit, max limit: %d", MIN_LIMIT)
	}

	return limit, offset, err
}

func NewTimelineAPI(db *sql.DB) *TimelineAPI {
	return &TimelineAPI{
		timeline: controller.NewTimelineController(db),
	}
}
