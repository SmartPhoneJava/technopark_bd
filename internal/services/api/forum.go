package api

import (
	"escapade/internal/models"
	re "escapade/internal/return_errors"
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
		if err.Error() == re.ErrorForumUserNotExist().Error() {
			rw.WriteHeader(http.StatusNotFound)
			sendErrorJSON(rw, err, place)
		} else {
			rw.WriteHeader(http.StatusConflict)
			sendSuccessJSON(rw, forum, place)
		}
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, forum, place)
	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetForum(rw http.ResponseWriter, r *http.Request) {
	const place = "GetForum"
	var (
		forum models.Forum
		slug  string
		err   error
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = h.getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if forum, err = h.DB.GetForum(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, forum, place)
	printResult(err, http.StatusOK, place)
	return
}
