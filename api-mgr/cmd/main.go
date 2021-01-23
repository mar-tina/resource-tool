package main

import (
	"log"
	"net/http"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
	"github.com/mar-tina/resource-tool/api-mgr/internal/db"
	"github.com/mar-tina/resource-tool/api-mgr/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Missing .env file")
	}

	conf := config.DefaultConfig()
	arg.MustParse(conf)

	stormdb := db.InitStormDB()
	db.SetImpl(stormdb)

	handler := service.InitHandler(conf)

	router := mux.NewRouter()
	router.HandleFunc("/resource/exists", service.BasicAuth(http.HandlerFunc(handler.CheckIFResourceExists), conf.USERNAME, conf.PASSWORD, "username and password please"))

	log.Printf("starting app on PORT: %s", conf.PORT)
	err = http.ListenAndServe(":"+conf.PORT, router)
	if err != nil {
		log.Fatal(err)
	}
}
