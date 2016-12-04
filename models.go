package main

import "code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"

type DomainsResponse struct {
	NextUrl   string        `json:"next_url"`
	Resources []ccv2.Domain `json:"resources"`
}

type RoutesResponse struct {
	NextUrl   string       `json:"next_url"`
	Resources []ccv2.Route `json:"resources"`
}

type App struct {
	Entity struct {
		Name      string `json:"name"`
		SpaceGUID string `json:"space_guid"`
		SpaceURL  string `json:"space_url"`
	} `json:"entity"`
}

type AppsResponse struct {
	NextUrl   string `json:"next_url"`
	Resources []App  `json:"resources"`
}
