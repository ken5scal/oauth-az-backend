package main // import "github.com/you/hello"

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var oauthConfig oauth2.Config
var port string

func init() {
	tomlInBytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed reading a toml file: %v", err.Error()))
	}

	config, err := toml.LoadBytes(tomlInBytes)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed parsing a toml file: %v", err.Error()))
	}

	port = strconv.FormatInt(config.Get("port").(int64), 10)
}

func main() {
	// assuming administrative requests from frontend
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	// Options is for CORS preflight
	allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodGet, http.MethodPost})

	r := mux.NewRouter()
	//r.HandleFunc("/hoge", fuga)
	srv := &http.Server {
		Addr: "localhost:" + port,
		Handler: handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r),
	}
	log.Fatal(srv.ListenAndServe())
}
