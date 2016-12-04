package main

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
)

type DomainsResponse struct {
	NextUrl   string        `json:"next_url"`
	Resources []ccv2.Domain `json:"resources"`
}

type RoutesResponse struct {
	NextUrl   string       `json:"next_url"`
	Resources []ccv2.Route `json:"resources"`
}

type Mapping struct {
	Entity struct {
		AppGUID string `json:"app_guid"`
		AppURL  string `json:"app_url"`
	} `json:"entity"`
}

type MappingsResponse struct {
	NextUrl   string    `json:"next_url"`
	Resources []Mapping `json:"resources"`
}
