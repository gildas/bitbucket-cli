package profile

import (
	"context"

	"github.com/gildas/go-logger"
)

type PaginatedResources[T any] struct {
	Values   []T `json:"values"`
	Page     int `json:"page"`
	PageSize int `json:"pagelen"`
	Size     int `json:"size"`
}

// GetAllResources gets all resources using the given profile
func GetAllResources[T any](context context.Context, profile *Profile) (resources []T, err error) {
	log := logger.Must(logger.FromContext(context, Log)).Child(nil, "getall")

	log.Infof("Getting all resources for profile %s", profile.Name)

	return resources, nil
}

func getResourcesAtPage[T any](context context.Context, profile *Profile, page int) (resources []T, err error) {
	log := logger.Must(logger.FromContext(context, Log)).Child(nil, "getall")

	log.Infof("Getting all resources for profile %s", profile.Name)

	return resources, nil
}
