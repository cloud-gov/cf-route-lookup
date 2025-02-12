# Cloud Foundry Route Lookup Plugin

This is a Cloud Foundry CLI plugin to find the application a given hostname/domain/route is bound to. Note this will only show applications in orgs/spaces that the logged-in user has permissions to view.

## End-of-Life

This plugin is no longer supported. Equivalent functionality is provided by the [CloudFoundry cf-lookup-route plugin](https://github.com/cloudfoundry/cf-lookup-route).

## Installation

1. Download the appropriate binary from [the Releases page](https://github.com/18F/cf-route-lookup/releases).
1. Run

    ```sh
    cf install-plugin -r CF-Community route-lookup
    ```

## Usage

```
$ cf lookup-route <my.example.com>
Bound to:
<org>/<space>/<app>

# use -t to target the org/space containing the route
$ cf lookup-route -t <my.example.com>
Bound to:
<org>/<space>/<app>
Changed target to: <org>/<space>

$ cf lookup-route <unknown.example.com>
Error retrieving apps: Route not found.
```
