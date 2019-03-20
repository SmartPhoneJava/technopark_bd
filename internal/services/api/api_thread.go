package api

import (
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
	"time"
)

func (h *Handler) CreateThread(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateThread"
	var (
		thread models.Thread
		err    error
	)

	rw.Header().Set("Content-Type", "application/json")

	if thread, err = getThread(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, err = h.DB.CreateThread(&thread); err != nil {
		if err.Error() == re.ErrorThreadConflict().Error() {
			rw.WriteHeader(http.StatusConflict)
			sendSuccessJSON(rw, thread, place)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			sendErrorJSON(rw, err, place)

		}
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetThreads(rw http.ResponseWriter, r *http.Request) {
	const place = "GetThreads"
	var (
		threads    []models.Thread
		slug       string
		limit      int
		t          time.Time
		err        error
		existLimit bool
		existTime  bool
		desc       bool
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existLimit, limit, err = getThreadLimit(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existTime, t, err = getThreadTime(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if desc, err = getThreadDesc(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if threads, err = h.DB.GetThreads(slug, limit, existLimit, t, existTime, desc); err != nil {
		//if err.Error() == re.ErrorForumUserNotExist().Error() {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		// } else {
		// 	rw.WriteHeader(http.StatusConflict)
		// 	sendSuccessJSON(rw, forum, place)
		// }
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, threads, place)
	printResult(err, http.StatusOK, place)
	return
}
