package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lithammer/shortuuid/v3"
	"github.com/mar-tina/resource-tool/api-mgr/internal/config"
	"github.com/mar-tina/resource-tool/api-mgr/internal/db"
	"github.com/mar-tina/resource-tool/api-mgr/internal/postman"
	"github.com/mar-tina/resource-tool/api-mgr/internal/schema"
)

type Payload struct {
	Collection string `json:"collection"`
	Resource   string `json:"resource"`
}

type Handler struct {
	conf *config.Config
	pc   *postman.PostmanClient
}

func InitHandler(conf *config.Config, pc *postman.PostmanClient) *Handler {
	return &Handler{
		conf: conf,
		pc:   pc,
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
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var temp schema.Collection = dst
	dst.ID = ""
	dst.Name = ""
	dst.Owner = ""
	dst.UID = ""

	var payload schema.CollectionPayload
	payload.Collection = &dst
	err = h.pc.CreateCollection(payload)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.StoreCollection(&temp)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, "collection created", http.StatusOK)
}

func (h *Handler) CreateEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst schema.Environment
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var payload schema.EnvPayload
	payload.Environment = &dst
	err = h.pc.CreateServiceEnv(payload)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload.Environment.ID = shortuuid.New()
	err = db.CreateEnvironment(payload.Environment)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, "environment created", http.StatusOK)
}

func (h *Handler) FetchResoure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst schema.FetchResoure
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	split := strings.Split(dst.Key, ".")
	env := split[0]
	srvc := split[1]

	key, err := db.FetchKeyInEnvironment(env, srvc)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resource, err := db.FetchCollectionResource(split[1], dst.Resource)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	replaced := TrimURL(resource.URL)
	resource.URL = fmt.Sprintf("%s%s", key, replaced)

	err = db.UpdateDescendants(srvc, dst.ServiceID, dst.Resource)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// err = db.UpdateAncestors(dst.ServiceID, srvc, dst.Resource)
	// if err != nil {
	// 	ResponseError(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	ResponseOk(w, resource, http.StatusOK)
}

func (h *Handler) FetchSingleEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst map[string]interface{}
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := h.pc.FetchSingleEnv(dst["id"].(string))
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, resp, http.StatusOK)
}

func (h *Handler) FetchAllEnvHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := h.pc.FetchAllEnv()
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, resp, http.StatusOK)
}

func (h *Handler) UpdateEnvHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst schema.UpdateEnv
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Printf("HEREIEDINED")
	payload := make(map[string]interface{})
	payload["environment"] = dst.Env

	log.Printf("%v", payload)

	err = h.pc.UpdateEnvironment(dst.Env.ID, payload)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.UpdateEnvironment(&dst.Env)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, "environment updated", http.StatusOK)
}

func (h *Handler) FetchAllCollections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dst schema.FetchAllCollections
	err := DecodeJsonBody(w, r, &dst)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			ResponseError(w, mr.msg, mr.status)
			return
		}

		log.Println(err.Error())
		ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := db.AllCollections(dst.Limit, dst.Skip)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseOk(w, resp, http.StatusOK)
}

func TrimURL(str string) string {
	initial := 0
	close := 0

	for i, x := range str {
		if x == '{' {
			initial += i
		}

		if x == '}' {
			close = i
			break
		}
	}

	initial = initial - 1
	close = close + 1

	return str[close+1:]
}
