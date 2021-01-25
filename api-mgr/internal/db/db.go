package db

import "github.com/mar-tina/resource-tool/api-mgr/internal/schema"

type Repository interface {
	// Names of the resource that's being accessed and the name of the service that is accessing it.
	StoreCollection(collection *schema.Collection) error
	FetchCollection(collectionName string) (*schema.Collection, error)
	CreateEnvironment(env *schema.Environment) error
	FetchKeyInEnvironment(env, key string) (string, error)
	FetchCollectionResource(servicename, resource string) (*schema.Request, error)
	UpdateDescendants(parent, child, resource string) error
	UpdateAncestors(child, parent, resource string) error
	AllCollections(limit, skip int) ([]*schema.Collection, error)
	UpdateEnvironment(env *schema.Environment) error
	// AddEndpointAncestors(child, parent, resource string) error
}

var impl Repository

func SetImpl(repo Repository) {
	impl = repo
}

func CreateEnvironment(env *schema.Environment) error {
	return impl.CreateEnvironment(env)
}

func StoreCollection(collection *schema.Collection) error {
	return impl.StoreCollection(collection)
}

func FetchCollection(collectionName string) (*schema.Collection, error) {
	return impl.FetchCollection(collectionName)
}

func FetchKeyInEnvironment(env, key string) (string, error) {
	return impl.FetchKeyInEnvironment(env, key)
}

func FetchCollectionResource(servicename, resource string) (*schema.Request, error) {
	return impl.FetchCollectionResource(servicename, resource)
}

func UpdateDescendants(parent, child, resource string) error {
	return impl.UpdateDescendants(parent, child, resource)
}

func AllCollections(limit, skip int) ([]*schema.Collection, error) {
	return impl.AllCollections(limit, skip)
}

func UpdateEnvironment(env *schema.Environment) error {
	return impl.UpdateEnvironment(env)
}

func UpdateAncestors(child, parent, resource string) error {
	return impl.UpdateAncestors(child, parent, resource)
}
