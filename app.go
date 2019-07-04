package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Repo struct {
	id   int
	name string
}

var repos []Repo

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Response struct {
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ReposHelper(w http.ResponseWriter, req *http.Request) {
	respondWithError(os.Getenv("REPOS_URL"), w)
}

func GeneralinfoHelper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	url := fmt.Sprintf("%s/%s", os.Getenv("GENERALINFO_URL"), params["repo"])

	respondWithError(url, w)
}

func AuthHelper(w http.ResponseWriter, req *http.Request) {
	username, password := req.FormValue("username"), req.FormValue("password")

	response := &Response{Error: false, Data: "", Message: ""}

	if username != os.Getenv("USERNAME") || password != os.Getenv("PASSWORD") {
		response.Error = true
		response.Message = "Wrong username and/or password"
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"ExpiresAt": time.Now().Add(48 * time.Hour)})
		tokenString, _ := token.SignedString(jwtKey)

		fmt.Println("Token is", tokenString)

		response.Data = map[string]interface{}{"token": token}
	}

	json.NewEncoder(w).Encode(response)
}

func respondWithError(url string, w http.ResponseWriter) {
	resp, result := fetchResource(url)

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	response := &Response{Error: false, Data: "", Message: ""}

	if resp.StatusCode >= 400 {
		response.Error = true
		response.Message = "Something went wrong"
	} else {
		response.Data = result
	}

	json.NewEncoder(w).Encode(response)
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
	router.HandleFunc("/auth", AuthHelper).Methods("POST")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
