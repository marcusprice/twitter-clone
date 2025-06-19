package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	limit, limitErr := strconv.Atoi(values.Get("limit"))
	offset, offsetErr := strconv.Atoi(values.Get("offset"))
	if limitErr != nil || offsetErr != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
	}

	if limit < MAX_LIMIT {
		http.Error(
			w,
			fmt.Sprintf("Too large of a limit, max limit: %d", MAX_LIMIT),
			http.StatusBadRequest,
		)

		return
	}

	if limit > MIN_LIMIT {
		http.Error(
			w,
			fmt.Sprintf("Too small of a limit, max limit: %d", MIN_LIMIT),
			http.StatusBadRequest,
		)

		return
	}

	timeline := timelineAPI.timeline.SetID(userID)
	posts, postsRemaining, err := timeline.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
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

func NewTimelineAPI(db *sql.DB) *TimelineAPI {
	return &TimelineAPI{
		timeline: controller.NewTimelineController(db),
	}
}
