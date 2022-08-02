// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package main

import (
	"context"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	register "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/config"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	installpb "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/installapp"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/installappserver"
	readynotifypb "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/readynotify"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/readynotifyserver"
	updatepb "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/updateapp"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/updateappserver"

	contextDb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/contextdb"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	con "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/context"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
)

func RegisterRsyncServices(grpcServer *grpc.Server, srv interface{}) {
	installpb.RegisterInstallappServer(grpcServer, installappserver.NewInstallAppServer())
	readynotifypb.RegisterReadyNotifyServer(grpcServer, readynotifyserver.NewReadyNotifyServer())
	updatepb.RegisterUpdateappServer(grpcServer, updateappserver.NewUpdateAppServer())
}

func createTracerProvider() (*tracesdk.TracerProvider, error) {
	endpoint := "http://" + net.JoinHostPort(config.GetConfiguration().ZipkinIP, "9411") + "/api/v2/spans"
	exp, err := zipkin.New(endpoint)
	if err != nil {
		return nil, err
	}
	name := "unknown"
	name, _ = os.LookupEnv("APP_NAME")
	namespace := "default"
	namespace, _ = os.LookupEnv("POD_NAMESPACE")
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name+"."+namespace),
		)),
	)
	return tp, nil
}

func main() {

	rand.Seed(time.Now().UnixNano())

	tp, err := createTracerProvider()
	if err != nil {
		log.Error("Unable to initialize tracing provider", log.Fields{"Error": err})
		os.Exit(1)
	}
	otel.SetTracerProvider(tp)

	// Istio uses b3 propagation
	otel.SetTextMapPropagator(b3.New())

	ctx := context.Background()

	// Initialize the mongodb
	err = db.InitializeDatabaseConnection(ctx, "emco")
	if err != nil {
		log.Error("Unable to initialize mongo database connection", log.Fields{"Error": err})
		os.Exit(1)
	}

	// Initialize contextdb
	err = contextDb.InitializeContextDatabase()
	if err != nil {
		log.Error("Unable to initialize etcd database connection", log.Fields{"Error": err})
		os.Exit(1)
	}

	go func() {
		err := register.StartGrpcServer("rsync", "RSYNC_NAME", 9031,
			RegisterRsyncServices, nil)
		if err != nil {
			log.Error("GRPC server failed to start", log.Fields{"Error": err})
			os.Exit(1)
		}
	}()

	err = con.RestoreActiveContext(ctx)
	if err != nil {
		log.Error("RestoreActiveContext failed", log.Fields{"Error": err})
	}

	connectionsClose := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(connectionsClose)

}
