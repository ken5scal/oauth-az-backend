package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	// Following configs read params from Environment Vars & Config File
	// Environment Vars are prioritized because
	// Dockerfile uses and set those variables as Environment Vars
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else if config.Get("port") != nil {
		port = strconv.FormatInt(config.Get("port").(int64), 10)
	}

	if os.Getenv("ADDR") != "" {
		addr = os.Getenv("ADDR")
	} else if config.Get("ADDR") != nil {
		addr = config.Get("ADDR").(string)
	}

	// Following params prioritize the config file
	// because they are the app settings
	if config.Get("debug") != nil {
		debug = config.Get("debug").(bool)
	} else if os.Getenv("debug") != "" {
		debug, err = strconv.ParseBool(os.Getenv("debug"))
		if err != nil {
			log.Fatal(fmt.Sprintf("failed parsing env var: %v", err.Error()))
		}
	}
}

func main() {
	// assuming administrative requests from frontend
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	// Options is for CORS preflight
	allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodGet, http.MethodPost})

	r := mux.NewRouter()
	r.HandleFunc("/", fuga)
	srv := &http.Server{
		Addr:    addr + ":" + port,
		Handler: handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r),
	}
	log.Fatal(srv.ListenAndServe())
}

func fuga(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fuga")
}
