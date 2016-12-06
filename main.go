package main

import (
	"fmt"
	"log"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

var CMD = "lookup-route"

type BasicPlugin struct{}

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

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != CMD {
		return
	}

	fmt.Println("Running " + CMD)

	if len(args) != 2 {
		log.Fatal("Please specify the domain to look up.")
	}

	hostname := args[1]
	apps, err := getApps(cliConnection, hostname)
	if err != nil {
		log.Fatal("Error retrieving apps: ", err)
	}

	for _, app := range apps {
		space, err := app.GetSpace(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving space: ", err)
		}
		org, err := space.GetOrg(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving org: ", err)
		}
		fmt.Println(org.Entity.Name + "/" + space.Entity.Name + "/" + app.Entity.Name)
	}
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "route-lookup",
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
				HelpText: "Look up the mapping of a provided route",
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
