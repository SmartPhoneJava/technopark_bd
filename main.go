package main

import (
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

// ./swag init

// @title Escapade API
// @version 1.0
// @description Documentation

// @host https://escapade-backend.herokuapp.com
// @BasePath /api/v1
func main() {
	const confPath = "conf.json"

	API, conf, err := api.GetHandler(confPath) // init.go
	if err != nil {
		fmt.Println("Some error happened with configuration file or database" + err.Error())
		return
	}

	r := mux.NewRouter()

	//var v = r.PathPrefix("/api").Subrouter()

	//v.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//var v1 = r.PathPrefix("/v1").Subrouter()

	/*
		tx := self.db.MustBegin()
			defer tx.Rollback()

			query := "INSERT INTO tasks (description, completed) VALUES"
			args := []interface{}{}
			for idx, item := range params.Body {
				if idx > 0 {
					query += ","
				}
				query += fmt.Sprintf(" ($%d, $%d)", len(args)+1, len(args)+2)
				args = append(args, item.Description, item.Completed)
	*/

	r.HandleFunc("/api/user/{name}/create", API.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{name}/profile", API.GetProfile).Methods("GET")

	// r.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	// r.HandleFunc("/user", mi.CORS(conf.Cors)(API.GetMyProfile)).Methods("GET")
	// r.HandleFunc("/user", mi.CORS(conf.Cors)(API.CreateUser)).Methods("POST")
	// r.HandleFunc("/user", mi.CORS(conf.Cors)(API.DeleteUser)).Methods("DELETE")
	// r.HandleFunc("/user", mi.CORS(conf.Cors)(API.UpdateProfile)).Methods("PUT")
	// r.HandleFunc("/user", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	// r.HandleFunc("/session", mi.CORS(conf.Cors)(API.Logout)).Methods("DELETE")
	// r.HandleFunc("/session", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	// r.HandleFunc("/session", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	// r.HandleFunc("/avatar", mi.CORS(conf.Cors)(API.GetImage)).Methods("GET")
	// r.HandleFunc("/avatar", mi.CORS(conf.Cors)(API.PostImage)).Methods("POST")
	// r.HandleFunc("/avatar", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	// r.HandleFunc("/users", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	// r.HandleFunc("/users/pages/{page}", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	// r.HandleFunc("/users/pages_amount", mi.CORS(conf.Cors)(API.GetUsersPageAmount)).Methods("GET")

	// r.HandleFunc("/users/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// r.HandleFunc("/users/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// r.HandleFunc("/users/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	fmt.Println("launched, look at us on " + conf.Server.Host + os.Getenv("PORT"))

	if err = http.ListenAndServe(os.Getenv("PORT"), r); err != nil { //os.Getenv("PORT"), r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
