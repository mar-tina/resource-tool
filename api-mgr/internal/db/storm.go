package db

import (
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

// func (s *StormDB) AddEndpointDescendants(collectionName, resourceName, serviceName string) error {
// 	coll := schema.Collection{}
// 	//fetch the collection
// 	err := s.db.One("Collection", collectionName, &coll)
// 	if err != nil {
// 		log.Printf("addEndpointDescendants failed %s", err)
// 		return errors.New("addEndpointDescendants failed")
// 	}

// 	var serviceExists bool
// 	// check if service already exists in the descendants list
// 	for idx, srvc := range coll.Descendants {
// 		//if exists update the list of resource names it is being consumed by
// 		if srvc.Name == serviceName {
// 			for _, rsc := range srvc.Resources {
// 				if rsc != resourceName {
// 					continue
// 				}
// 				srvc.Resources = append(srvc.Resources, resourceName)
// 			}
// 			serviceExists = true
// 			coll.Descendants[idx] = srvc
// 			s.db.Save(coll)
// 			return nil
// 		}
// 	}

// 	if !serviceExists {
// 		srvc := schema.Service{}
// 		srvc.Name = serviceName
// 		srvc.Resources = append(srvc.Resources, resourceName)
// 		coll.Descendants = append(coll.Descendants, srvc)
// 	}

// 	return nil
// }

// func (s *StormDB) AddEndpointAncestors(collectionName, resourceName, serviceName string) error {
// 	log.Printf
// 	coll := schema.Collection{}
// 	//fetch the collection
// 	err := s.db.One("Collection", collectionName, &coll)
// 	if err != nil {
// 		log.Printf("addEndpointAncestors failed %s", err)
// 		return errors.New("addEndpointAncestors failed")
// 	}

// 	var serviceExists bool
// 	// check if service already exists in the descendants list
// 	for idx, srvc := range coll.Ancestors {
// 		//if exists update the list of resource names it is being consumed by
// 		if srvc.Name == serviceName {
// 			for _, rsc := range srvc.Resources {
// 				if rsc != resourceName {
// 					continue
// 				}
// 				srvc.Resources = append(srvc.Resources, resourceName)
// 			}
// 			serviceExists = true
// 			coll.Ancestors[idx] = srvc
// 			s.db.Save(coll)
// 			return nil
// 		}
// 	}

// 	if !serviceExists {
// 		srvc := schema.Service{}
// 		srvc.Name = serviceName
// 		srvc.Resources = append(srvc.Resources, resourceName)
// 		coll.Descendants = append(coll.Ancestors, srvc)
// 	}

// 	return nil
// }

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

//CreateEnvironment TODO: hash the values of the keys
func (s *StormDB) CreateEnvironment(env *schema.Environment) error {
	err := s.db.Save(env)
	if err != nil {
		log.Printf("createEnvironment failed %s", err)
		return fmt.Errorf("createEnvironment failed %s", err)
	}

	return nil
}

func (s *StormDB) UpdateEnvironment(env *schema.Environment) error {
	err := s.db.Update(&schema.Environment{ID: env.ID, Values: env.Values, Name: env.Name})
	if err != nil {
		log.Printf("updateEnvironment failed %s", err)
		return fmt.Errorf("updateEnvironment failed %s", err)
	}

	return nil
}

func (s *StormDB) FetchKeyInEnvironment(env, key string) (string, error) {
	var enviro schema.Environment
	err := s.db.One("Name", env, &enviro)
	if err != nil {
		log.Printf("fetchKeyInEnvironment failed %s", err)
		return "", fmt.Errorf("fetchKeyInEnvironment failed %s", err)
	}

	var res string

	if &enviro != nil {
		for _, el := range enviro.Values {
			if el.Key == key {
				res = el.Value
				break
			}
		}
	}

	if res != "" {
		return res, nil
	}

	return "", fmt.Errorf("could not find specified key")
}

func (s *StormDB) FetchCollectionResource(servicename, resource string) (*schema.Request, error) {
	var coll schema.Collection
	err := s.db.One("Name", servicename, &coll)
	if err != nil {
		log.Printf("failed to find specified resource %s", err)
		return nil, fmt.Errorf("failed to find specified resource %s", err)
	}

	for _, item := range coll.Item {
		for _, inner := range item.Item {
			if inner.Name == resource {
				return &inner.Request, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find specified resource")
}

func (s *StormDB) UpdateDescendants(parent, child, resource string) error {
	var coll schema.Collection
	err := s.db.One("Name", parent, &coll)
	if err != nil {
		log.Printf("updateDescendants failed  %s", err)
		return fmt.Errorf("updateDescendants %s", err)
	}

	var exists bool
	var isChild bool

	for _, desc := range coll.Descendants {
		if desc.Name == child {
			isChild = true
			if isChild {
				for _, item := range desc.Resources {
					if resource == item {
						exists = true
						break
					}
				}
			}

			if !exists && isChild {
				desc.Resources = append(desc.Resources, resource)
				coll.Descendants = append(coll.Descendants, desc)
				err := s.db.Update(&schema.Collection{ID: coll.ID, Descendants: coll.Descendants})
				if err != nil {
					log.Printf("db failed 1 %s", err)
					return fmt.Errorf("updateDescendants %s", err)
				}
				return nil
			}
		}
	}

	if !exists && !isChild {
		srvc := schema.Service{}
		srvc.Name = child
		srvc.Resources = append(srvc.Resources, resource)
		coll.Descendants = append(coll.Descendants, srvc)
		err := s.db.Update(&schema.Collection{ID: coll.ID, Descendants: coll.Descendants})
		if err != nil {
			log.Printf("db failed  %s", err)
			return fmt.Errorf("updateDescendants here %s", err)
		}

	}

	return nil
}

func (s *StormDB) UpdateAncestors(child, parent, resource string) error {
	log.Printf("THE VALUES parent: %s child: %s", parent, child)
	var coll schema.Collection
	err := s.db.One("Name", child, &coll)
	if err != nil {
		log.Printf("UpdateAncestors failed  %s", err)
		return fmt.Errorf("UpdateAncestors %s", err)
	}

	var exists bool
	var isParent bool

	for _, desc := range coll.Ancestors {
		if desc.Name == parent {
			isParent = true
			if isParent {
				for _, item := range desc.Resources {
					if resource == item {
						exists = true
						break
					}
				}
			}

			if !exists && isParent {
				desc.Resources = append(desc.Resources, resource)
				coll.Descendants = append(coll.Ancestors, desc)
				err := s.db.Update(&schema.Collection{ID: coll.ID, Ancestors: coll.Ancestors})
				if err != nil {
					log.Printf("db failed 1 %s", err)
					return fmt.Errorf("UpdateAncestors %s", err)
				}
				return nil
			}
		}
	}

	if !exists && !isParent {
		srvc := schema.Service{}
		srvc.Name = child
		srvc.Resources = append(srvc.Resources, resource)
		coll.Descendants = append(coll.Ancestors, srvc)
		err := s.db.Update(&schema.Collection{ID: coll.ID, Ancestors: coll.Ancestors})
		if err != nil {
			log.Printf("db failed  %s", err)
			return fmt.Errorf("UpdateAncestors here %s", err)
		}

	}

	return nil
}

func (s *StormDB) AllCollections(limit, skip int) ([]*schema.Collection, error) {
	var colls []*schema.Collection
	err := s.db.All(&colls, storm.Limit(limit), storm.Skip(skip), storm.Reverse())
	if err != nil {
		log.Printf("could not find collections")
		return nil, fmt.Errorf("failed to fetch %s", err)
	}

	return colls, nil
}
