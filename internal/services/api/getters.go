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

func getPostID(r *http.Request) (id int, err error) {
	var (
		vars map[string]string
		str  string
	)

	vars = mux.Vars(r)

	if str = vars["id"]; str == "" {
		err = re.ErrorInvalidID()
		return
	}

	if id, err = strconv.Atoi(str); err != nil {
		return
	}
	if id < 0 {
		err = re.ErrorInvalidID()
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

func getLimit(r *http.Request) (exist bool, limit int, err error) {
	exist = true
	str := r.FormValue("limit")
	if str == "" {
		exist = false
		return
	}
	if limit, err = strconv.Atoi(str); err != nil {
		return
	}
	if limit < 0 {
		err = re.ErrorInvalidLimit()
		return
	}

	return
}

func getNickNameMin(r *http.Request) (exist bool, str string) {
	exist = true
	str = r.FormValue("since")
	if str == "" {
		exist = false
	}
	return
}

func getRelated(r *http.Request) (exist bool, str string) {
	exist = true
	str = r.FormValue("related")
	if str == "" {
		exist = false
	}
	return
}

func getIDmin(r *http.Request) (exist bool, since int, err error) {
	exist = true
	str := r.FormValue("since")
	if str == "" {
		exist = false
		return
	}

	if since, err = strconv.Atoi(str); err != nil {
		return
	} else {
		if since < 0 {
			err = re.ErrorInvalidDate()
			return
		}
	}

	return
}

func getTime(r *http.Request) (exist bool, t time.Time, err error) {

	exist = true
	str := r.FormValue("since")
	if str == "" {
		exist = false
		return
	}

	var num int

	if t, err = time.Parse("2006-01-02T15:04:05.000+03:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.000Z", str); err != nil {
			if num, err = strconv.Atoi(str); err != nil {
				return
			} else {
				if num < 0 {
					err = re.ErrorInvalidDate()
					return
				} else if num < 10000 {
					t = time.Date(num, 0, 0, 0, 0, 0, 0, time.UTC)
				} else {
					err = re.ErrorInvalidDate()
					return
				}
			}
		}

	}

	return
}

func getDesc(r *http.Request) (desc bool) {

	str := r.FormValue("desc")
	if str == "" {
		desc = false
		return
	}
	if str == "true" {
		desc = true
	}
	return
}

func getSort(r *http.Request) string {

	return r.FormValue("sort")
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

func getPosts(r *http.Request) (posts []models.Post, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&posts)

	return
}

func getVote(r *http.Request) (vote models.Vote, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&vote)

	return
}
