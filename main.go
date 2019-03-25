package main

import (
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

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

	r.HandleFunc("/api/user/{name}/create", API.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{name}/profile", API.GetProfile).Methods("GET")
	r.HandleFunc("/api/user/{name}/profile", API.UpdateProfile).Methods("POST")

	r.HandleFunc("/api/forum/create", API.CreateForum).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", API.GetForum).Methods("GET")

	r.HandleFunc("/api/forum/{slug}/create", API.CreateThread).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/threads", API.GetThreads).Methods("GET")

	r.HandleFunc("/api/thread/{slug}/details", API.GetThreadDetails).Methods("GET")
	r.HandleFunc("/api/thread/{slug}/details", API.UpdateThread).Methods("POST")

	r.HandleFunc("/api/thread/{slug}/create", API.CreatePosts).Methods("POST")
	r.HandleFunc("/api/thread/{slug}/posts", API.GetPosts).Methods("GET")

	r.HandleFunc("/api/thread/{slug}/vote", API.Vote).Methods("POST")

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	fmt.Println("launched, look at us on " + conf.Server.Host + os.Getenv("PORT"))

	if err = http.ListenAndServe(os.Getenv("PORT"), r); err != nil { //os.Getenv("PORT"), r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
