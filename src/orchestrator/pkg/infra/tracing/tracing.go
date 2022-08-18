// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package tracing

import (
	"net"
	"net/http"
	"os"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/config"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

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

func InitializeTracer() error {
	tp, err := createTracerProvider()
	if err != nil {
		return err
	}
	otel.SetTracerProvider(tp)
	// Istio uses b3 propagation
	otel.SetTextMapPropagator(b3.New())
	return nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		tracer := otel.Tracer("orchestrator")
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()
		next.ServeHTTP(w, r.Clone(ctx))
	})
}
