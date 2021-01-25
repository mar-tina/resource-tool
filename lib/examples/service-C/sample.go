package main

import (
	"encoding/json"
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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]string)
	data["service-name"] = "service-b"
	data["health"] = "I am really healthy"

	jsonBytes, _ := json.Marshal(data)
	w.Write(jsonBytes)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Missing .env file")
	}

	resource := os.Getenv("RESOURCE")
	mgr := apimgr.New(resource, "service_c").WithBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	mgr.NewEnv("service-a").EnvVar("service_c", "http://localhost:7758")

	coll := mgr.NewCollection("service_c", apimgr.DefaultCollOpts())
	coll.HandleFunc("c-endpoint", "{{service_c}}", "/health", "POST", HealthHandler, TestPayload{
		Test: "",
		Word: "",
	}).Header("Authorization", "Bearer token")

	coll.HandleFunc("bye-endpoint", "{{service_c}}", "/bye", "POST", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bye"))
	}, TestPayload{
		Test: "",
		Word: "",
	}).Header("Authorization", "Basic auth")

	collection := coll.Requests("sample-requests-c")
	err = mgr.Create(&collection)
	if err != nil {
		log.Printf("failed to create collection %s", err)
	}

	_, err = mgr.Use("service-a.service_a", "bye-endpoint")
	if err != nil {
		log.Printf("failed to fetch resource %s", err)
	}

	_, err = mgr.Use("service-b.service_b", "health-check")
	if err != nil {
		log.Printf("failed to fetch resource %s", err)
	}

	err = http.ListenAndServe(":7758", coll.GetRouter())
	if err != nil {
		log.Printf("something went wrong starting server %s", err)
	}
}
