package main

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	"code.cloudfoundry.org/cli/plugin"
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

func (m Mapping) GetApp(cliConnection plugin.CliConnection) (App, error) {
	return getApp(cliConnection, m.Entity.AppGUID)
}

type MappingsResponse struct {
	NextUrl   string    `json:"next_url"`
	Resources []Mapping `json:"resources"`
}

type App struct {
	Name      string `json:"name"`
	SpaceGUID string `json:"space_guid"`
}
