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

type TestNested struct {
	Test struct {
		Nested    string `json:"nested"`
		SuperNest struct {
			Inner string `json:"inner"`
		} `json:"super_nested"`
	} `json:"test"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]string)
	data["service-name"] = "service-b"
	data["data"] = "this is all the information i have for you"

	jsonBytes, _ := json.Marshal(data)
	w.Write(jsonBytes)
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
	// Init a new manager with basic auth properties.
	// Could be implemented for tokens and/or certs.
	mgr := apimgr.New(resource, "service-b").WithBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	mgr.NewEnv("service-b").EnvVar("service_b", "http://localhost:7757")

	c := mgr.NewCollection("service-b", apimgr.DefaultCollOpts())
	c.HandleFunc("health-check", "{{service_b}}", "/health", "POST", HealthHandler, TestPayload{})
	c.HandleFunc("test", "{{service-b-endpoint}}", "/test", "POST", TestHandler, TestNested{}).Header("Authorization", "Bearer {{token_special}}")

	requests := c.Requests("samples")
	err = mgr.Create(&requests)
	if err != nil {
		log.Printf("req creation failed %s", err)
	}

	//Accessing other endpoints and getting indexed as part of consumers of those endpoints.
	_, err = mgr.Use("service-a.service_a", "bye-endpoint")
	if err != nil {
		log.Printf("resource is unavailable %s", err)
	}

	_, err = mgr.Use("service-c.service_c", "c-endpoint")
	if err != nil {
		log.Printf("resource is unavailable %s", err)
	}

	err = http.ListenAndServe(":7757", c.GetRouter())
	if err != nil {
		log.Printf("something went wrong starting server %s", err)
	}
}
