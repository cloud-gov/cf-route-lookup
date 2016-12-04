package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	"code.cloudfoundry.org/cli/plugin"
)

func apiCall(cliConnection plugin.CliConnection, path string) (body string, err error) {
	// based on https://github.com/krujos/cfcurl/blob/320854091a119f220102ba356e507c361562b221/cfcurl.go
	bodyLines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", path)
	if err != nil {
		return
	}
	body = strings.Join(bodyLines, "\n")
	return
}

func getDomains(cliConnection plugin.CliConnection, names []string) (domains []ccv2.Domain, err error) {
	// based on https://github.com/ECSTeam/buildpack-usage/blob/e2f7845f96c021fa7f59d750adfa2f02809e2839/command/buildpack_usage_cmd.go#L161-L167

	domains = make([]ccv2.Domain, 0)

	endpoints := [...]string{"/v2/private_domains", "/v2/shared_domains"}

	params := url.Values{}
	params.Set("q", "name IN "+strings.Join(names, ","))
	params.Set("results-per-page", "100")
	queryString := params.Encode()

	for _, endpoint := range endpoints {
		uri := endpoint + "?" + queryString

		// paginate
		for uri != "" {
			var body string
			body, err = apiCall(cliConnection, uri)
			if err != nil {
				return
			}

			var data DomainsResponse
			err = json.Unmarshal([]byte(body), &data)
			if err != nil {
				return
			}

			domains = append(domains, data.Resources...)
			uri = data.NextUrl
		}
	}

	return
}

func getDomain(cliConnection plugin.CliConnection, hostname string) (matchingDomain ccv2.Domain, found bool, err error) {
	possibleDomains := getPossibleDomains(hostname)

	domains, err := getDomains(cliConnection, possibleDomains)
	if err != nil {
		return
	}
	fmt.Println("Matching domains:", domains)

	for _, possibleDomain := range possibleDomains {
		for _, domain := range domains {
			if domain.Name == possibleDomain {
				found = true
				matchingDomain = domain
				return
			}
		}
	}

	return
}

func getRoutes(cliConnection plugin.CliConnection, domain ccv2.Domain) (routes []ccv2.Route, err error) {
	// based on https://github.com/ECSTeam/buildpack-usage/blob/e2f7845f96c021fa7f59d750adfa2f02809e2839/command/buildpack_usage_cmd.go#L161-L167

	routes = make([]ccv2.Route, 0)

	params := url.Values{}
	// TODO also filter by host
	params.Set("q", "domain_guid:"+domain.GUID)
	params.Set("results-per-page", "100")
	uri := "/v2/routes?" + params.Encode()

	// paginate
	for uri != "" {
		var body string
		body, err = apiCall(cliConnection, uri)
		if err != nil {
			return
		}

		var data RoutesResponse
		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return
		}

		routes = append(routes, data.Resources...)
		uri = data.NextUrl
	}

	return
}

func getRoute(cliConnection plugin.CliConnection, hostname string) (matchingRoute ccv2.Route, found bool, err error) {
	domain, domainFound, err := getDomain(cliConnection, hostname)
	if err != nil {
		return
	}
	if !domainFound {
		err = errors.New("Could not find matching domain.")
		return
	}

	routes, err := getRoutes(cliConnection, domain)
	if err != nil {
		return
	}
	fmt.Println(len(routes), "routes found.")

	for _, route := range routes {
		routeHostname := domain.Name
		if route.Host != "" {
			routeHostname = route.Host + "." + routeHostname
		}
		if routeHostname == hostname {
			found = true
			matchingRoute = route
			return
		}
	}

	return
}

func getApps(cliConnection plugin.CliConnection, hostname string) (apps []App, err error) {
	route, routeFound, err := getRoute(cliConnection, hostname)
	if err != nil {
		return
	}
	if !routeFound {
		err = errors.New("Route not found.")
		return
	}
	fmt.Println("Route found! GUID:", route.GUID)

	apps = make([]App, 0)
	uri := "/v2/routes/" + route.GUID + "/apps"

	// paginate
	for uri != "" {
		var body string
		body, err = apiCall(cliConnection, uri)
		if err != nil {
			return
		}

		var data AppsResponse
		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return
		}

		apps = append(apps, data.Resources...)
		uri = data.NextUrl
	}

	return
}
