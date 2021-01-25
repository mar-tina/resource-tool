package postman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
	"github.com/mar-tina/resource-tool/api-mgr/internal/schema"
)

type envFuncMap map[string]func(k, v string, env *schema.Environment) *schema.Environment

var FuncEnvMap envFuncMap

func InitFuncMap() {
	FuncEnvMap["create"] = AddNewKey
	FuncEnvMap["delete"] = RemoveKey
	FuncEnvMap["update"] = UpdateKey
}

type Postman interface {
	FindCollection(collection string) (map[string]interface{}, error)
	CreateCollection()
	Call()
	FetchSingleCollection() (map[string]interface{}, error)
}

type PostmanClient struct {
	conf *config.Config
}

func InitPostmanClient(conf *config.Config) *PostmanClient {
	return &PostmanClient{
		conf: conf,
	}
}

func (pc *PostmanClient) prepareRequest(method, route string, payload interface{}) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("failed to prepare request")
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", pc.conf.URI, route), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request %s", err)
	}

	req.Header.Add("X-API-Key", pc.conf.Token)

	client := http.Client{}
	return client.Do(req)
}

func (pc *PostmanClient) CreateCollection(collection schema.CollectionPayload) error {
	resp, err := pc.prepareRequest("POST", "/collections", collection)
	if err != nil {
		return fmt.Errorf("failed to create collection %s", err)
	}

	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)

	if body["error"] != nil {
		return fmt.Errorf("failed to create collection %s", (body["error"].(map[string]interface{})["message"]))
	}

	return nil
}

func (pc *PostmanClient) FetchSingleEnv(envID string) (map[string]interface{}, error) {
	resp, err := pc.prepareRequest("GET", fmt.Sprintf("/environments/%s", envID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch environment %s", err)
	}

	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)

	if body["error"] != nil {
		return nil, fmt.Errorf("failed to fetch environment %s", (body["error"].(map[string]interface{})["message"]))
	}

	return body, nil
}

func (pc *PostmanClient) CreateServiceEnv(env schema.EnvPayload) error {
	resp, err := pc.prepareRequest("POST", "/environments", env)
	if err != nil {
		return fmt.Errorf("failed to create environment %s", err)
	}

	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)

	if body["error"] != nil {
		return fmt.Errorf("failed to create environment %s", (body["error"].(map[string]interface{})["message"]))
	}

	return nil
}

func (pc *PostmanClient) FetchAllEnv() (map[string]interface{}, error) {
	resp, err := pc.prepareRequest("GET", "/environments", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all environments %s", err)
	}

	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)

	if body["error"] != nil {
		return nil, fmt.Errorf("failed to fetch environments %s", (body["error"].(map[string]interface{})["message"]))
	}

	log.Printf("body %v", body)

	return body, nil
}

func (pc *PostmanClient) UpdateEnvironment(envID string, env map[string]interface{}) error {

	resp, err := pc.prepareRequest("PUT", fmt.Sprintf("/environments/%s", envID), env)
	if err != nil {
		return fmt.Errorf("failed to update environment %s", err)
	}

	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)

	log.Printf("Bud %v", body)
	if body["error"] != nil {
		return fmt.Errorf("failed to update environment %s", (body["error"].(map[string]interface{})["message"]))
	}
	return nil
}

func RemoveKey(key, value string, env *schema.Environment) *schema.Environment {
	for idx, el := range env.Values {
		if el.Key == key {
			env.Values[idx] = env.Values[len(env.Values)-1]
			env.Values[len(env.Values)-1] = nil
			env.Values = env.Values[:len(env.Values)-1]
		}
	}

	return env
}

func UpdateKey(key, value string, env *schema.Environment) *schema.Environment {
	for idx, el := range env.Values {
		if el.Key == key {
			env.Values[idx].Value = value
		}
	}

	return env
}

func AddNewKey(key, value string, env *schema.Environment) *schema.Environment {
	env.Values = append(env.Values, &schema.EnvValues{
		Key:   key,
		Value: value,
	})

	return env
}

func ChooseEnvAction(action string) func(k, v string, env *schema.Environment) *schema.Environment {
	return FuncEnvMap[action]
}
