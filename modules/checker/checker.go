package checker

import "github.com/TechMinerApps/upmaster/models"

type Checker interface {
}

type checker struct {
	endpointPool []*models.Endpoint
}
