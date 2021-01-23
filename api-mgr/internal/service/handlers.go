package service

import (
	"errors"
	"log"
	"net/http"

	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
	"github.com/mar-tina/resource-tool/api-mgr/internal/db"
	"github.com/mar-tina/resource-tool/api-mgr/internal/schema"
)

type ResourceCheckPayload struct {
	collectionName string `json:"collection"`
	resourceName   string `json:"resource"`
}

type Payload struct {
	Collection string `json:"collection"`
	Resource   string `json:"resource"`
}

type Handler struct {
	conf *config.Config
}

func InitHandler(conf *config.Config) *Handler {
	return &Handler{
		conf: conf,
	}
}

func (h *Handler) CheckIFResourceExists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var payload Payload
	err := DecodeJsonBody(w, r, &payload)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	coll, err := db.FetchCollection(payload.Collection)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, resource := range coll.Resources {
		if resource.Name == payload.Resource {
			ResponseOk(w, "resource found", http.StatusOK)
			return
		}
	}

	ResponseError(w, "resource not found", http.StatusNoContent)
}

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst schema.Collection
	err := SimpleDecodeJSON(w, r, &dst)
	if err != nil {
		ResponseError(w, err, -1)
		return
	}

	err = db.StoreCollection(&dst)
	if err != nil {
		ResponseError(w, err, -1)
		return
	}

	ResponseOk(w, "collection created", 1)
}
