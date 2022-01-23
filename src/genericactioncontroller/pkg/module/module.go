package module

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

// Client for using the services
type Client struct {
	// Add Clients for API's here
	Customization    *CustomizationClient
	GenericK8sIntent *GenericK8sIntentClient
	Resource         *ResourceClient
}

// ClientDbInfo holds the MongoDB collection and attributes info
type ClientDbInfo struct {
	storeName  string // name of the mongodb collection to use for client documents
	tagMeta    string // attribute key name for the json data of a client document
	tagContent string // attribute key name for the file data of a client document
}

// Metadata holds the data
type Metadata struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"-"`
	UserData1   string `json:"userData1" yaml:"-"`
	UserData2   string `json:"userData2" yaml:"-"`
}

// NewClient creates a new client for using the services
func NewClient() *Client {
	c := &Client{}
	c.GenericK8sIntent = NewGenericK8sIntentClient()
	c.Resource = NewResourceClient()
	c.Customization = NewCustomizationClient()
	return c
}
