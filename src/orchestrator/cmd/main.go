// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/api"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/auth"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/config"
	contextDb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/contextdb"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/rpc"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module/controller"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	err := db.InitializeDatabaseConnection("emco")
	if err != nil {
		log.Println("Unable to initialize mongo database connection...")
		log.Println(err)
		log.Fatalln("Exiting...")
	}
	err = contextDb.InitializeContextDatabase()
	if err != nil {
		log.Println("Unable to initialize etcd database connection...")
		log.Println(err)
		log.Fatalln("Exiting...")
	}

	httpRouter := api.NewRouter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)
	log.Println("Starting Kubernetes Multicloud API")

	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	controller.NewControllerClient("resources", "data", "orchestrator").InitControllers()

	connectionsClose := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		httpServer.Shutdown(context.Background())
		rpc.CloseAllRpcConn()
		close(connectionsClose)
	}()

	tlsConfig, err := auth.GetTLSConfig("ca.cert", "server.cert", "server.key")
	if err != nil {
		log.Println("WARNING :: Getting TLS Configuration failed. Starting without TLS...")
		log.Fatal(httpServer.ListenAndServe())
	} else {
		httpServer.TLSConfig = tlsConfig
		// empty strings because tlsconfig already has this information
		err = httpServer.ListenAndServeTLS("", "")
	}
}
