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
	repondWithError(os.Getenv("REPOS_URL"), w)
}

func GeneralinfoHelper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	url := fmt.Sprintf("%s/%s", os.Getenv("GENERALINFO_URL"), params["repo"])

	repondWithError(url, w)
}

func repondWithError(url string, w http.ResponseWriter) {
	resp, result := fetchResource(url)

	w.WriteHeader(resp.StatusCode)

	str_response := map[string]interface{}{"error": false, "data": nil, "message": ""}

	if resp.StatusCode >= 400 {
		str_response["error"] = true
		str_response["message"] = "Something went wrong"
	} else {
		str_response["data"] = result
	}

	json.NewEncoder(w).Encode(str_response)

}

func fetchResource(url string) (*http.Response, interface{}) {
	resp, error := http.Get(url)
	if error != nil {
		panic(error.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var data interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err.Error())
	}

	return resp, data
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	port := os.Getenv("PORT")

	fmt.Println(fmt.Sprintf("The server is running on port %s", port))

	router := mux.NewRouter()
	router.HandleFunc("/repos", ReposHelper).Methods("GET")
	router.HandleFunc("/generalinfo/{repo}", GeneralinfoHelper).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
