package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	apimgr "github.com/mar-tina/resource-tool/lib/pkg"
)

type TestPayload struct {
	Test string `json:"test"`
	Word string `json:"word"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Missing .env file")
	}

	resource := os.Getenv("RESOURCE")
	mgr := apimgr.New(resource, "simple-service").WithBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	mgr.NewEnv("service-a").EnvVar("service_a", "http://localhost:7756")

	coll := mgr.NewCollection("service_a", apimgr.DefaultCollOpts())
	coll.HandleFunc("hello-endpoint", "{{service_a}}", "/hello", "POST", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}, TestPayload{
		Test: "",
		Word: "",
	}).Header("Authorization", "Bearer token")

	coll.HandleFunc("bye-endpoint", "{{service_a}}", "/bye", "POST", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bye"))
	}, TestPayload{
		Test: "",
		Word: "",
	}).Header("Authorization", "Basic auth")

	collection := coll.Requests("sample-requests")
	err = mgr.Create(&collection)
	if err != nil {
		log.Printf("failed to create collection %s", err)
	}

	_, err = mgr.Use("service-b.service_b", "bye-endpoint")
	if err != nil {
		log.Printf("failed to fetch resource %s", err)
	}

	_, err = mgr.Use("service-c.service_c", "bye-endpoint")
	if err != nil {
		log.Printf("failed to fetch resource %s", err)
	}

	err = http.ListenAndServe(":7756", coll.GetRouter())
	if err != nil {
		log.Printf("something went wrong starting server %s", err)
	}
}
