// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"gitlab.com/project-emco/core/emco-base/src/ncm/api"
	"gitlab.com/project-emco/core/emco-base/src/ncm/pkg/statusnotify"
	register "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc"
	statusnotifypb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotify"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/auth"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/config"
	contextDb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/contextdb"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

func startGrpcServer() error {
	var tls bool

	if strings.Contains(config.GetConfiguration().GrpcEnableTLS, "enable") {
		tls = true
	} else {
		tls = false
	}
	certFile := config.GetConfiguration().GrpcServerCert
	keyFile := config.GetConfiguration().GrpcServerKey

	_, port := register.GetServerHostPort("ncm", "NCM_NAME", 9082)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Could not listen to port: %v", err)
	}
	var opts []grpc.ServerOption
	if tls {
		if certFile == "" {
			certFile = testdata.Path("server.pem")
		}
		if keyFile == "" {
			keyFile = testdata.Path("server.key")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Could not generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	statusnotifypb.RegisterStatusNotifyServer(grpcServer, statusnotify.StartStatusNotifyServer())

	log.Println("Starting NCM Status Notification gRPC Server")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("NCM Status Notification gRPC server is not serving %v", err)
	}
	return err
}

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

	httpRouter := api.NewRouter(nil)
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)
	log.Println("Starting Network Customization Manager")

	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	go func() {
		err := startGrpcServer()
		if err != nil {
			log.Fatalf("GRPC server failed to start")
		}
	}()

	connectionsClose := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		httpServer.Shutdown(context.Background())
		close(connectionsClose)
	}()

	tlsConfig, err := auth.GetTLSConfig("ca.cert", "server.cert", "server.key")
	if err != nil {
		log.Println("Error Getting TLS Configuration. Starting without TLS...")
		log.Fatal(httpServer.ListenAndServe())
	} else {
		httpServer.TLSConfig = tlsConfig
		// empty strings because tlsconfig already has this information
		err = httpServer.ListenAndServeTLS("", "")
	}
}
