// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package sdewancc

const NETWORKING_SDEWAN_APIVERSION = "batch.sdewan.akraino.org/v1alpha1"

type Metadata struct {
	Name        string `yaml:"name"`
	Namespace   string `yaml:"namespace,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func createGenericMetadata(name string, namespace string, description string)(Metadata) {
	var meta = Metadata {
		Name:        name,
		Namespace:   namespace,
		Description: description,
	}
	return meta
}
