package api

import (
	"escapade/internal/models"
	"net/http"
)

func (h *Handler) CreateForum(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateForum"
	var (
		forum models.Forum
		err   error
	)

	rw.Header().Set("Content-Type", "application/json")

	if forum, err = getForum(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if forum, err = h.DB.CreateForum(&forum); err != nil {
		rw.WriteHeader(http.StatusConflict)
		sendSuccessJSON(rw, forum, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, forum, place)
	printResult(err, http.StatusCreated, place)
	return
}
