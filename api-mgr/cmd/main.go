package main

import (
	"log"
	"net/http"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
	"github.com/mar-tina/resource-tool/api-mgr/internal/db"
	"github.com/mar-tina/resource-tool/api-mgr/internal/postman"
	"github.com/mar-tina/resource-tool/api-mgr/internal/service"
	"github.com/rs/cors"
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

	pc := postman.InitPostmanClient(conf)

	handler := service.InitHandler(conf, pc)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	})

	router := mux.NewRouter()
	router.HandleFunc("/resource/exists", service.BasicAuth(http.HandlerFunc(handler.CheckIFResourceExists), conf.USERNAME, conf.PASSWORD, "username and password please"))
	router.HandleFunc("/collection/create", service.BasicAuth(http.HandlerFunc(handler.CreateCollection), conf.USERNAME, conf.PASSWORD, "username and password please"))
	router.HandleFunc("/env/create", service.BasicAuth(http.HandlerFunc(handler.CreateEnv), conf.USERNAME, conf.PASSWORD, "username and password please"))
	router.HandleFunc("/resource/fetch", service.BasicAuth(http.HandlerFunc(handler.FetchResoure), conf.USERNAME, conf.PASSWORD, "username and password please"))
	router.HandleFunc("/env/fetch", handler.FetchSingleEnv)
	router.HandleFunc("/env/all", handler.FetchAllEnvHandler)
	router.HandleFunc("/collections/all", handler.FetchAllCollections)
	router.HandleFunc("/env/update", handler.UpdateEnvHandler)

	r := c.Handler(router)
	log.Printf("starting app on PORT: %s", conf.PORT)
	err = http.ListenAndServe(":"+conf.PORT, r)
	if err != nil {
		log.Fatal(err)
	}
}
