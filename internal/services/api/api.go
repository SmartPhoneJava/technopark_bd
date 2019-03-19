package api

import (
	database "escapade/internal/database"
	"escapade/internal/models"
	"fmt"
	"net/http"
)

// Handler is struct
type Handler struct {
	DB                    database.DataBase
	PlayersAvatarsStorage string
	FileMode              int
}

// catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK1
// @Success 200 "successfully"
// @Router /user [OPTIONS]
func (h *Handler) Ok(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
	return
}

// CreateUser create new user
// @Summary create new user
// @Description create new user
// @ID Register
// @Success 201 {object} models.Result "Create user successfully"
// @Failure 400 {object} models.Result "Invalid information"
// @Router /user [POST]
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateUser"
	var (
		user  models.User
		err   error
		users *[]models.User
	)

	rw.Header().Set("Content-Type", "application/json")

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user.Nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	users, user, err = h.DB.CreateUser(&user)

	if len(*users) > 0 {
		rw.WriteHeader(http.StatusConflict)
		printResult(err, http.StatusConflict, place)
		sendSuccessJSON(rw, users, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		sendErrorJSON(rw, err, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, user, place)
	printResult(err, http.StatusCreated, place)
	return
}

// GetProfile returns model UserPublicInfo
// @Summary Get some of user fields
// @Description return public information, such as name or best_score
// @ID GetProfile
// @Param name path string false "User name"
// @Success 200 {object} models.UserPublicInfo "Profile found successfully"
// @Failure 400 {object} models.Result "Invalid username"
// @Failure 404 {object} models.Result "User not found"
// @Router /users/{name}/profile [GET]
func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "GetProfile"

	var (
		err      error
		nickname string
		user     models.User
	)

	rw.Header().Set("Content-Type", "application/json")

	if nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, err = h.DB.GetUser(nickname); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	sendSuccessJSON(rw, user, place)
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}

// UpdateProfile updates profile
// @Summary update user information
// @Description update public info
// @ID UpdateProfile
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 401 {object} models.Result "need authorization"
// @Router /user [PUT]
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateProfile"

	var (
		user models.User
		err  error
	)

	rw.Header().Set("Content-Type", "application/json")

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user.Nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	user.Print()

	if user, err = h.DB.UpdateUser(user); err != nil {
		rw.WriteHeader(http.StatusConflict)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusConflict, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, user, place)
	printResult(err, http.StatusOK, place)
	return
}
