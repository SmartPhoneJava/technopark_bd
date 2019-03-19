package api

import (
	"encoding/json"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
	"strconv"
	"time"

	//"reflect"

	"github.com/gorilla/mux"
)

// bd
func getNickname(r *http.Request) (nickname string, err error) {
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

func getSlug(r *http.Request) (slug string, err error) {
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

func getUser(r *http.Request) (user models.User, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}

func getThreadLimit(r *http.Request) (exist bool, limit int, err error) {
	exist = true
	str := r.FormValue("limit")
	if str == "" {
		exist = false
		return
	}
	limit, err = strconv.Atoi(str)

	return
}

func getThreadTime(r *http.Request) (exist bool, t time.Time, err error) {

	exist = true
	str := r.FormValue("since")
	if str == "" {
		exist = false
		return
	}

	t, err = time.Parse("2006-01-02T15:04:05.000+03:00", str)

	return
}

func getThreadDesc(r *http.Request) (desc bool, t time.Time, err error) {

	str := r.FormValue("desc")
	if str == "" {
		desc = false
		return
	}
	if str == "desc" {
		desc = true
	}
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

func getThread(r *http.Request) (user models.Thread, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}
