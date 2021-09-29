// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package config

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

// Configuration loads up all the values that are used to configure
// backend implementations
type Configuration struct {
	CAFile              string `json:"ca-file"`
	ServerCert          string `json:"server-cert"`
	ServerKey           string `json:"server-key"`
	Password            string `json:"password"`
	DatabaseAddress     string `json:"database-address"`
	DatabaseType        string `json:"database-type"`
	PluginDir           string `json:"plugin-dir"`
	EtcdIP              string `json:"etcd-ip"`
	EtcdCert            string `json:"etcd-cert"`
	EtcdKey             string `json:"etcd-key"`
	EtcdCAFile          string `json:"etcd-ca-file"`
	ServicePort         string `json:"service-port"`
	KubernetesLabelName string `json:"kubernetes-label-name"`
}

// Config is the structure that stores the configuration
var gConfig *Configuration

// readConfigFile reads the specified smsConfig file to setup some env variables
func readConfigFile(file string) (*Configuration, error) {
	f, err := os.Open(file)
	if err != nil {
		return defaultConfiguration(), err
	}
	defer f.Close()

	// Setup some defaults here
	// If the json file has values in it, the defaults will be overwritten
	conf := defaultConfiguration()

	// Read the configuration from json file
	decoder := json.NewDecoder(f)
	err = decoder.Decode(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func defaultConfiguration() *Configuration {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting cwd. Using .")
		cwd = "."
	}

	return &Configuration{
		CAFile:              "ca.cert",
		ServerCert:          "server.cert",
		ServerKey:           "server.key",
		Password:            "",
		DatabaseAddress:     "127.0.0.1",
		DatabaseType:        "mongo",
		PluginDir:           cwd,
		EtcdIP:              "127.0.0.1",
		EtcdCert:            "",
		EtcdKey:             "",
		EtcdCAFile:          "",
		ServicePort:         "9015",
		KubernetesLabelName: "k8splugin.io/rb-instance-id",
	}
}

// GetConfiguration returns the configuration for the app.
// It will try to load it if it is not already loaded.
func GetConfiguration() *Configuration {
	if gConfig == nil {
		conf, err := readConfigFile("k8sconfig.json")
		if err != nil {
			log.Println("Error loading config file. Using defaults.")
		}
		gConfig = conf
	}

	return gConfig
}

// SetConfigValue sets a value in the configuration
// This is mostly used to customize the application and
// should be used carefully.
func SetConfigValue(key string, value string) *Configuration {
	c := GetConfiguration()
	if value == "" || key == "" {
		return c
	}

	v := reflect.ValueOf(c).Elem()
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(key)
		if f.IsValid() {
			if f.CanSet() {
				if f.Kind() == reflect.String {
					f.SetString(value)
				}
			}
		}
	}
	return c
}
