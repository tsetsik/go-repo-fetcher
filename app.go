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

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Response struct {
	Error      bool        `json:"error"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	statusCode int
}

const TOKEN_HEADER = "X-JWT-TOKEN"

type Token struct {
	token *jwt.Token
	str   string
}

func (t *Token) initiateFromStr() error {
	token, err := jwt.Parse(t.str, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err == nil {
		t.token = token
	}

	return err
}

func (t Token) Valid() bool {
	return t.token.Valid
}

func (t *Token) Create() string {
	t.token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"ExpiresAt": time.Now().Add(48 * time.Hour)})
	t.str, _ = t.token.SignedString(jwtKey)

	return t.str
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
		w.WriteHeader(http.StatusUnauthorized)
		response.Error = true
		response.Message = "Unatuhorized request"
	} else {
		token := &Token{}
		w.Header().Set(TOKEN_HEADER, token.Create())

		response.Data = map[string]interface{}{"username": username}
	}

	json.NewEncoder(w).Encode(response)
}

func respondWithError(url string, w http.ResponseWriter) {
	resp, result := fetchResource(url)

	respond(resp.StatusCode, result, "Something went wrong", w)
}

func respond(status int, data interface{}, msg string, w http.ResponseWriter) {
	response := &Response{Error: false, Data: data, Message: msg, statusCode: http.StatusOK}

	if status >= 400 {
		response.Error = true
		response.statusCode = status
		response.Message = "Something went wrong"
	} else {
		response.Data = data
	}

	w.WriteHeader(response.statusCode)
	w.Header().Set("Content-Type", "application/json")

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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		has_error := false

		if req.URL.Path != "/auth" {
			token := &Token{str: req.Header.Get(TOKEN_HEADER)}
			if err := token.initiateFromStr(); err != nil || !token.Valid() {
				has_error = true
				respond(http.StatusUnauthorized, "", "Unauthorized access", w)
			}
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if has_error == false {
			next.ServeHTTP(w, req)
		}
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	port := os.Getenv("PORT")

	fmt.Println(fmt.Sprintf("The server is running on port %s", port))

	router := mux.NewRouter()

	// Register auth middleware
	router.Use(authMiddleware)

	router.HandleFunc("/repos", ReposHelper).Methods("GET")
	router.HandleFunc("/generalinfo/{repo}", GeneralinfoHelper).Methods("GET")
	router.HandleFunc("/auth", AuthHelper).Methods("POST")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
