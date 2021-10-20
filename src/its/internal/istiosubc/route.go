// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package istiosubc

type HTTPRoute struct {
	Route   []HTTPRouteDestination   `yaml:"route,omitempty"`
}

type HTTPRouteDestination struct {
	Dest Destination `yaml:"destination,omitempty"`
	Weight int32     `yaml:"weight,omitempty"`
}

type Destination struct {
	Host string `yaml:"host"`
	Subset string `yaml:"subset,omitempty"`
	Port   PortSelector `yaml:"port,omitempty"`
}

type PortSelector struct {
	Number uint32 `yaml:"number,omitempty"`
}

func createPortSelector(number uint32)(PortSelector) {
	var portSelector = PortSelector {
		Number: number,
	}
	return portSelector
}

func createDestination(host string, port PortSelector)(Destination) {
	var dest = Destination {
		Host: host,
		Port: port,
	}
	return dest

}

func createHTTPRouteDestination(dest Destination, weight int32)(HTTPRouteDestination) {
	var httpRoutedest = HTTPRouteDestination {
		Dest:   dest,
		Weight: weight,
	}
	return httpRoutedest
}

func createHTTPRoute(route []HTTPRouteDestination)(HTTPRoute) {
	var httpRoute = HTTPRoute {
		Route: route,
	}

	return httpRoute
}
