package api

import (
	"escapade/internal/models"
	"net/http"
)

func (h *Handler) CreatePosts(rw http.ResponseWriter, r *http.Request) {
	const place = "CreatePosts"
	var (
		posts []models.Post
		err   error
	)

	rw.Header().Set("Content-Type", "application/json")

	if posts, err = getPosts(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if posts, err = h.DB.CreatePost(posts); err != nil {
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
	sendSuccessJSON(rw, posts, place)
	printResult(err, http.StatusCreated, place)
	return
}
