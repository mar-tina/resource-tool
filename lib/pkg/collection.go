package manager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
	"github.com/mar-tina/resource-tool/lib/schema"
)

type CollOpts struct {
	Schema      string `json:"schema"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
}

func DefaultCollOpts() *CollOpts {
	return &CollOpts{
		Description: "no descritpion provided",
		Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		Mode:        "raw",
	}
}

type Coll struct {
	requests       map[string]*schema.Request
	collection     schema.Collection
	Router         *mux.Router
	currentrequest string
}

func (mgr *ApiMgr) NewCollection(name string, opts *CollOpts) *Coll {
	return &Coll{
		collection: schema.Collection{
			ID:   shortuuid.New(),
			Name: name,
			Info: schema.Info{
				Name:        name,
				Schema:      opts.Schema,
				Description: opts.Description,
			},
			Item: make([]schema.ItemEntry, 1, 10),
		},
		requests: make(map[string]*schema.Request),
		Router:   mux.NewRouter(),
	}
}

func (c *Coll) HandleFunc(name, host, path, method string, f func(http.ResponseWriter, *http.Request), payload interface{}) *Coll {
	if name == "" {
		name = path
	}

	c.requests[name] = createNewRequest(name, host, path, method, payload)
	c.Router.HandleFunc(path, f)
	c.currentrequest = name
	return c
}

func (c *Coll) Header(key, value string) {
	c.requests[c.currentrequest].Header = append(c.requests[c.currentrequest].Header, schema.Header{
		Key:   key,
		Value: value,
	})
}

func (c *Coll) GetRouter() *mux.Router {
	return c.Router
}

func (c *Coll) Requests(name string) schema.Collection {
	var collectnested []schema.NestedItem
	for key, val := range c.requests {
		temp := schema.NestedItem{
			Name:    key,
			Request: *val,
		}

		collectnested = append(collectnested, temp)
	}

	c.collection.Item[0].Item = collectnested
	c.collection.Item[0].Name = name

	// c.collection.Item = *items
	return c.collection
}

func createNewRequest(name, host, route, method string, payload interface{}) *schema.Request {
	jsonBytes, _ := json.Marshal(payload)
	return &schema.Request{
		URL:    fmt.Sprintf("%s%s", host, route),
		Method: method,
		Header: []schema.Header{},
		Body:   schema.Body{Mode: "raw", Raw: string(jsonBytes)},
	}
}
