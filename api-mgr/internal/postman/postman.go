package postman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
)

type Postman interface {
	FindCollection(collection string) (map[string]interface{}, error)
	CreateCollection()
	Call()
	FetchSingleCollection() (map[string]interface{}, error)
}

type PostmanClient struct {
}

func prepareRequest(method, route string, payload interface{}, conf *config.Config) (*http.Request, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("failed to prepare request")
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", conf.URI, route), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request %s", err)
	}

	req.Header.Add("X-API-Key", conf.Token)
	return req, nil
}
