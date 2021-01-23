package db

import "github.com/mar-tina/resource-tool/api-mgr/internal/schema"

type Repository interface {
	AddEndpointAncestors(collectionName, resourceName, serviceName string) error
	// Names of the resource that's being accessed and the name of the service that is accessing it.
	AddEndpointDescendants(collectionName, resourceName, serviceName string) error
	StoreCollection(collection *schema.Collection) error
	FetchCollection(collectionName string) (*schema.Collection, error)
}

var impl Repository

func SetImpl(repo Repository) {
	impl = repo
}

func AddEndpointDescendants(collectionName, resourceName, serviceName string) error {
	return impl.AddEndpointDescendants(collectionName, resourceName, serviceName)
}

func StoreCollection(collection *schema.Collection) error {
	return impl.StoreCollection(collection)
}

func FetchCollection(collectionName string) (*schema.Collection, error) {
	return impl.FetchCollection(collectionName)
}

func AddEndpointAncestors(collectionName, resourceName, serviceName string) error {
	return impl.AddEndpointAncestors(collectionName, resourceName, serviceName)
}
