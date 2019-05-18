package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Repo struct {
	id   int
	name string
}

var repos []Repo

func ReposHelper(w http.ResponseWriter, req *http.Request) {
	all_repos := fetchResource("https://api.github.com/users/vmware/repos")
	json.NewEncoder(w).Encode(all_repos)
}

func fetchResource(url string) interface{} {
	resp, error := http.Get(url)
	if error != nil {
		panic(error.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/repos", ReposHelper).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", router))
}
