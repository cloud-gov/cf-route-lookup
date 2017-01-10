package main

import (
	"flag"
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
	log.SetFlags(0)

	flags := flag.NewFlagSet(CMD, flag.ContinueOnError)
	target := flags.Bool("t", false, "Target the org / space containing this route")
	flags.Parse(args[1:])

	if len(flags.Args()) != 1 {
		log.Fatal("Please specify the domain to look up.")
	}

	hostname := flags.Args()[0]

	apps, err := getApps(cliConnection, hostname)
	if err != nil {
		log.Fatal("Error retrieving apps: ", err)
	}

	if len(apps) == 0 {
		log.Println("Not bound to any applications.")
		return
	}
	log.Println("Bound to:")

	for _, app := range apps {
		space, err := app.GetSpace(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving space: ", err)
		}
		org, err := space.GetOrg(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving org: ", err)
		}
		log.Println(org.Entity.Name + "/" + space.Entity.Name + "/" + app.Entity.Name)
	}

	if *target {
		_, err = apps[0].Target(cliConnection)
		if err != nil {
			log.Fatal("Error targeting app: ", err)
		}

		space, err := apps[0].GetSpace(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving space: ", err)
		}
		org, err := space.GetOrg(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving org: ", err)
		}

		log.Printf("Changed target to %s/%s.\n", org.Entity.Name, space.Entity.Name)
	}

}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "route-lookup",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
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
					Usage: "\n   cf " + CMD + " [-t] <some.domain.com>",
					Options: map[string]string{
						"t": "Target the org / space containing the route",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BasicPlugin))
}
