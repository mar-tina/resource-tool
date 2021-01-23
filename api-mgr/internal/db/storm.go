package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/asdine/storm"
	"github.com/mar-tina/resource-tool/api-mgr/internal/schema"
)

type StormDB struct {
	db *storm.DB
}

func InitStormDB() *StormDB {
	db, err := storm.Open("rs-tool.db")
	if err != nil {
		log.Printf("failed to init db %s", err)
	}

	return &StormDB{
		db: db,
	}
}

func (s *StormDB) AddEndpointDescendants(collectionName, resourceName, serviceName string) error {
	coll := schema.Collection{}
	//fetch the collection
	err := s.db.One("Collection", collectionName, &coll)
	if err != nil {
		log.Printf("addEndpointDescendants failed %s", err)
		return errors.New("addEndpointDescendants failed")
	}

	var serviceExists bool
	// check if service already exists in the descendants list
	for idx, srvc := range coll.Descendants {
		//if exists update the list of resource names it is being consumed by
		if srvc.Name != serviceName {
			for _, rsc := range srvc.Resources {
				if rsc != resourceName {
					continue
				}
				srvc.Resources = append(srvc.Resources, resourceName)
			}
			serviceExists = true
			coll.Descendants[idx] = srvc
			s.db.Save(coll)
			return nil
		}
	}

	if !serviceExists {
		srvc := schema.Service{}
		srvc.Name = serviceName
		srvc.Resources = append(srvc.Resources, resourceName)
		coll.Descendants = append(coll.Descendants, srvc)
	}

	return nil
}

func (s *StormDB) AddEndpointAncestors(collectionName, resourceName, serviceName string) error {
	coll := schema.Collection{}
	//fetch the collection
	err := s.db.One("Collection", collectionName, &coll)
	if err != nil {
		log.Printf("addEndpointAncestors failed %s", err)
		return errors.New("addEndpointAncestors failed")
	}

	var serviceExists bool
	// check if service already exists in the descendants list
	for idx, srvc := range coll.Ancestors {
		//if exists update the list of resource names it is being consumed by
		if srvc.Name != serviceName {
			for _, rsc := range srvc.Resources {
				if rsc != resourceName {
					continue
				}
				srvc.Resources = append(srvc.Resources, resourceName)
			}
			serviceExists = true
			coll.Ancestors[idx] = srvc
			s.db.Save(coll)
			return nil
		}
	}

	if !serviceExists {
		srvc := schema.Service{}
		srvc.Name = serviceName
		srvc.Resources = append(srvc.Resources, resourceName)
		coll.Descendants = append(coll.Ancestors, srvc)
	}

	return nil
}

func (s *StormDB) StoreCollection(collection *schema.Collection) error {
	err := s.db.Save(collection)
	if err != nil {
		log.Printf("storeCollection failed %s", err)
		return fmt.Errorf("StoreCollection failed %s", err)
	}

	return nil
}

func (s *StormDB) FetchCollection(collectionName string) (*schema.Collection, error) {
	coll := schema.Collection{}
	err := s.db.One("Collection", collectionName, &coll)
	if err != nil {
		log.Printf("fetchCollection failed %s", err)
		return nil, fmt.Errorf("fetchCollection failed: %s", err)
	}

	return &coll, nil
}
