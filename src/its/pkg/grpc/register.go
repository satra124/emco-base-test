// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package grpc

import (
	"os"
	"strconv"
	"strings"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

const default_host = "localhost"
const default_port = 9040
const default_its_name = "its"
const ENV_IT_NAME = "ITS"

func GetServerHostPort() (string, int) {

	// expect name of this its program to be in env variable "IT_NAME" - e.g. IT_NAME="its"
	serviceName := os.Getenv(ENV_IT_NAME)
	if serviceName == "" {
		serviceName = default_its_name
		log.Info("Using default name for IT service name", log.Fields{
			"Name": serviceName,
		})
	}

	// expect service name to be in env variable - e.g. IT_SERVICE_HOST
	host := os.Getenv(strings.ToUpper(serviceName) + "_SERVICE_HOST")
	if host == "" {
		host = default_host
		log.Info("Using default host for its gRPC controller", log.Fields{
			"Host": host,
		})
	}

	// expect service port to be in env variable - e.g. IT_SERVICE_PORT
	port, err := strconv.Atoi(os.Getenv(strings.ToUpper(serviceName) + "_SERVICE_PORT"))
	if err != nil || port < 0 {
		port = default_port
		log.Info("Using default port for its gRPC controller", log.Fields{
			"Port": port,
		})
	}
	return host, port
}
