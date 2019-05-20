package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Repo struct {
	id   int
	name string
}

var repos []Repo

func ReposHelper(w http.ResponseWriter, req *http.Request) {
	all_repos := fetchResource(os.Getenv("REPOS_URL"))
	json.NewEncoder(w).Encode(all_repos)
}

func GeneralinfoHelper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	url := fmt.Sprintf("%s/%s", os.Getenv("GENERALINFO_URL"), params["repo"])
	general_info := fetchResource(url)
	json.NewEncoder(w).Encode(general_info)

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
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	port := os.Getenv("PORT")

	fmt.Println(fmt.Sprintf("The server is running on port %s", port))

	router := mux.NewRouter()
	router.HandleFunc("/repos", ReposHelper).Methods("GET")
	router.HandleFunc("/generalinfo/{repo}", GeneralinfoHelper).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
