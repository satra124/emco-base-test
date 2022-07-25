// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"gitlab.com/project-emco/core/emco-base/src/clm/api"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/config"
	contextDb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/contextdb"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	err := db.InitializeDatabaseConnection(ctx, "emco")
	if err != nil {
		log.Error("Unable to initialize mongo database connection", log.Fields{"Error": err})
		os.Exit(1)
	}
	err = contextDb.InitializeContextDatabase()
	if err != nil {
		log.Error("Unable to initialize etcd database connection", log.Fields{"Error": err})
		os.Exit(1)
	}

	httpRouter := api.NewRouter(nil)
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)
	log.Info("Starting Cluster Manager", log.Fields{"Port": config.GetConfiguration().ServicePort})

	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	connectionsClose := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		httpServer.Shutdown(ctx)
		close(connectionsClose)
	}()

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Error("HTTP server failed", log.Fields{"Error": err})
	}
}
