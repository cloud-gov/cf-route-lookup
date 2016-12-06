package main

import (
	"encoding/json"

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

type Space struct {
	Entity struct {
		Name    string `json:"name"`
		OrgGUID string `json:"organization_guid"`
		OrgURL  string `json:"organization_url"`
	} `json:"entity"`
}

type App struct {
	Entity struct {
		Name      string `json:"name"`
		SpaceGUID string `json:"space_guid"`
		SpaceURL  string `json:"space_url"`
	} `json:"entity"`
}

func (a App) GetSpace(cliConnection plugin.CliConnection) (space Space, err error) {
	body, err := apiCall(cliConnection, a.Entity.SpaceURL)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(body), &space)
	if err != nil {
		return
	}

	return
}

type AppsResponse struct {
	NextUrl   string `json:"next_url"`
	Resources []App  `json:"resources"`
}
