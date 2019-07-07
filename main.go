package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ken5scal/oauth-az/handler"
	"github.com/ken5scal/oauth-az/infrastructure"
	"github.com/pelletier/go-toml"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var oauthConfig oauth2.Config
var port = "8080"
var debug = true
var addr = "localhost"
var configFileLocation = "config.toml"

func init() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		configFileLocation = os.Args[1]
	}

	tomlInBytes, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed reading a toml file: %v", err.Error()))
	}

	config, err := toml.LoadBytes(tomlInBytes)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed parsing a toml file: %v", err.Error()))
	}

	if config.Get("port") != nil {
		port = strconv.FormatInt(config.Get("port").(int64), 10)
	}

	if config.Get("addr") != nil {
		addr = config.Get("addr").(string)
	}

	if config.Get("debug") != nil {
		debug = config.Get("debug").(bool)
	} else if os.Getenv("DEBUG") != "" {
		debug, err = strconv.ParseBool(os.Getenv("debug"))
		if err != nil {
			log.Fatal(fmt.Sprintf("failed parsing env var: %v", err.Error()))
		}
	}
}

func main() {
	cnn, _ := sql.Open("mysql", "dataSourceName")
	repo := infrastructure.NewTokenRepository(cnn)
	service := handler.NewService(repo)
	ctrl := handler.NewHandler(service)

	// assuming administrative requests from frontend
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	// Options is for CORS preflight
	allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodGet, http.MethodPost})

	r := mux.NewRouter()
	r.HandleFunc("/", fuga)
	r.HandleFunc("/hoge", ctrl.RequestToken)
	srv := &http.Server{
		Addr:    addr + ":" + port,
		Handler: handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r),
	}
	log.Fatal(srv.ListenAndServe())
}

func fuga(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fuga")
}
