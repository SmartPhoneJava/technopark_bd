package api

import (
	"encoding/json"
	"escapade/internal/models"
	"fmt"
	"net/http"
)

func printResult(catched error, number int, place string) {
	if place != " " {
		return
	}
	if catched != nil {
		fmt.Println("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		fmt.Println("api/"+place+" success(code:", number, ")")
	}
}

func sendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	json.NewEncoder(rw).Encode(result)
}

func sendSuccessJSON(rw http.ResponseWriter, result interface{}, place string) {
	if result == nil {
		result = models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}

	//bytes, _ := json.Marshal(result)
	//fmt.Println("result2:" + string(bytes))
	//rw.Write(bytes)
	json.NewEncoder(rw).Encode(result)
}

// func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) error {

// 	var (
// 		user models.UserPublicInfo
// 		err  error
// 	)

// 	if user, err = h.DB.GetProfile(username); err != nil {
// 		return err
// 	}

// 	sendSuccessJSON(rw, user, place)
// 	return err
// }
