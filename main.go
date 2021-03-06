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

	_ "github.com/lib/pq"
)

var oauthConfig oauth2.Config
var port = "8080"
var debug = true
var addr = "localhost"
var configFileLocation = "config.toml"
var psqlInfo string

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

	if config.Get("db_host") != nil &&
		config.Get("db_port") != nil &&
		config.Get("db_user") != nil &&
		config.Get("db_password") != nil &&
		config.Get("db_name") != nil {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Get("db_host"), config.Get("db_port"), config.Get("db_user"), config.Get("db_password"), config.Get("db_name"))
	} else {
		log.Fatal(fmt.Sprintf("failed retrieving db config: %v", err.Error()))
	}
}

func main() {
	cnn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer cnn.Close()

	tokenRepo := infrastructure.NewTokenRepository(cnn)
	authzRepo := infrastructure.NewAuthzInfoRepositoryImpl(cnn)

	// assuming administrative requests from frontend
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodGet, http.MethodPost})

	r := mux.NewRouter()
	// https://tools.ietf.org/html/rfc6749#section-3.1
	// URI  MUST NOT include a fragment component.
	// MUST support the use of the HTTP "GET" method
	r.HandleFunc("/authorize", handler.NewAuthzHandler(authzRepo).RequestAuthz).Methods(http.MethodGet)
	r.HandleFunc("/token", handler.NewTokenHandler(tokenRepo).RequestToken).Methods(http.MethodPost)
	srv := &http.Server{
		Addr:    addr + ":" + port,
		Handler: handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r),
	}
	log.Fatal(srv.ListenAndServe())

	//bcrypt.CompareHashAndPassword()
}
