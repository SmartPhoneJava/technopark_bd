package api

import (
	"encoding/json"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"

	//"reflect"

	"github.com/gorilla/mux"
)

// bd
func (h *Handler) getNickname(r *http.Request) (nickname string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if nickname = vars["name"]; nickname == "" {
		err = re.ErrorInvalidName()
		return
	}
	return
}

func (h *Handler) getSlug(r *http.Request) (slug string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if slug = vars["slug"]; slug == "" {
		err = re.ErrorForumSlugInvalid()
		return
	}
	return
}

// bd
func getUser(r *http.Request) (user models.User, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}

func getForum(r *http.Request) (user models.Forum, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}
