package coros

import (
	"fmt"
	"strings"
)

// Region is a COROS data region. Each region has its own API host.
type Region struct {
	Name string
	URL  string
}

// Regions lists the supported COROS regions in the order they are offered to the user.
var Regions = []Region{
	{Name: "America", URL: "https://teamapi.coros.com"},
	{Name: "Europe", URL: "https://teameuapi.coros.com"},
	{Name: "China", URL: "https://teamcnapi.coros.com"},
}

// ParseRegion finds a region by name, ignoring case and surrounding spaces.
func ParseRegion(value string) (Region, error) {
	value = strings.TrimSpace(value)

	for _, region := range Regions {
		if strings.EqualFold(region.Name, value) {
			return region, nil
		}
	}

	return Region{}, fmt.Errorf("unknown region %q; expected America, Europe or China", value)
}
