package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	"code.cloudfoundry.org/cli/plugin"
)

var CMD = "basic-plugin-command"

type BasicPlugin struct{}

type domainsResponse struct {
	NextUrl   string        `json:"next_url"`
	Resources []ccv2.Domain `json:"resources"`
}

type routesResponse struct {
	NextUrl   string       `json:"next_url"`
	Resources []ccv2.Route `json:"resources"`
}

// possibleDomains returns all domain levels, down to the second-level domain (SLD), in order.
func getPossibleDomains(hostname string) []string {
	parts := strings.Split(hostname, ".")
	numCombinations := len(parts) - 1
	possibleDomains := make([]string, numCombinations)
	for i := 0; i < numCombinations; i++ {
		possibleDomains[i] = strings.Join(parts[i:], ".")
	}
	return possibleDomains
}

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

			var data domainsResponse
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

		var data routesResponse
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

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != CMD {
		return
	}

	fmt.Println("Running the " + CMD)

	if len(args) != 2 {
		log.Fatal("Please specify the domain to look up.")
	}

	hostname := args[1]
	route, routeFound, err := getRoute(cliConnection, hostname)
	if err != nil {
		log.Fatal("Error finding route. ", err)
	}
	if !routeFound {
		log.Fatal("Route not found.")
	}
	fmt.Println("Route found! GUID:", route.GUID)
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "MyBasicPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     CMD,
				HelpText: "Basic plugin command's help text",
				UsageDetails: plugin.Usage{
					Usage: CMD + "\n   cf " + CMD,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BasicPlugin))
}
