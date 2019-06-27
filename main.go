package main

import (
	"flag"
	"log"
	"sort"
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

	if len(args) == 1 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		log.Println("Lack of features? Go to the https://github.com/18F/cf-route-lookup/issues")
		return
	}

	if len(args) >= 1 && args[0] != CMD {
		log.Println("Type --help...")
		return
	}

	flags := flag.NewFlagSet(CMD, flag.ContinueOnError)
	target := flags.Bool("t", false, "Target the org / space containing this route")
	_ = flags.Parse(args[1:])

	if len(flags.Args()) != 1 {
		log.Fatal("Please specify the domain to look up.")
	}

	hostname := flags.Args()[0]

	appsExt, err := getApps(cliConnection, hostname)
	if err != nil {
		log.Fatal("Error retrieving apps: ", err)
	}

	if len(appsExt) == 0 {
		log.Println("Not bound to any applications.")
		return
	}
	log.Println("Bound to:")

	var firstApp App
	found := false

	sort.Slice(appsExt, func(i, j int) bool {
		return appsExt[i].route.route.SpaceGUID < appsExt[j].route.route.SpaceGUID
	})

	spaceGuid := ""
	for _, appExt := range appsExt {
		if appExt.app != (App{}) {
			if !found {
				firstApp = appExt.app
				found = true
			}
		}

		if spaceGuid != appExt.route.route.SpaceGUID {
			spaceGuid = appExt.route.route.SpaceGUID
			org, space, err := getTargetBySpaceGuid(cliConnection, spaceGuid)
			if err != nil {
				log.Println("Cannot find the space with GUID=" + spaceGuid)
				return
			}

			log.Printf("\n> cf target -o %v -s %v\n", org.Entity.Name, space.Entity.Name)
		}
		if appExt.app == (App{}) {
			log.Printf("  > # unbounded route: %v (%v)\n", appExt.route.hostname, appExt.route.route.GUID)
		} else {
			log.Printf("  > cf app %v # route: %v (%v)\n", appExt.app.Entity.Name, appExt.route.hostname, appExt.route.route.GUID)
		}
	}

	if *target {
		_, err = firstApp.Target(cliConnection)
		if err != nil {
			log.Fatal("Error targeting app: ", err)
		}

		space, err := firstApp.GetSpace(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving space: ", err)
		}
		org, err := space.GetOrg(cliConnection)
		if err != nil {
			log.Fatal("Error retrieving org: ", err)
		}

		log.Println("Changed target to:", org.Entity.Name+"/"+space.Entity.Name)
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
