package api

import (
	"escapade/internal/models"
	"net/http"
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

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusCreated, place)
	return
}
