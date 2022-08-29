```
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2022 Intel Corporation
```
# Edge Multi-Cloud Orchestrator (EMCO) Observability
EMCO metrics and tracing builds upon the features provided by Istio sidecars. The first step to enabling observability in EMCO is to setup Istio with its telemetry addons.

## Setup EMCO with Istio
A more complete description of the steps to setup EMCO with Istio can be found in [EMCO Integrity and Access Management](Emco_Integrity_Access_Management.md). The steps listed here include only the minimum required to enable the observability of EMCO and may not be suitable in production.

These steps need to be followed in the Kubernetes Cluster where EMCO is installed.

### Install Istio with addons
Install Istio using your preferred method from the Istio [Installation Guides](https://istio.io/latest/docs/setup/install/).

Install the addons as described in [Telemetry Addons](https://github.com/istio/istio/tree/master/samples/addons).

### Configure Istio sidecar injection for emco namespace
```shell
$ kubectl label namespace emco istio-injection=enabled
```

### Install EMCO in the emco namespace
Use the EMCO Helm chart to install EMCO in the emco namespace. The EMCO services will come up with Istio sidecars.

## Access addon dashboards
The Istio addon dashboards can be accessed at the following URLs in the cluster.

| Dashbaord  | URL                          |
| ---------- | ---------------------------- |
| kiali      | http://localhost:20001/kiali |
| prometheus | http://localhost:9090        |
| jaeger     | http://localhost:16686       |
| grafana    | http://localhost:3000        |

## Jaeger configuration
General instructions are at [Configure tracing using MeshConfig and Pod annotations](https://istio.io/latest/docs/tasks/observability/distributed-tracing/mesh-and-proxy-config/).

To capture all samples, the IstioOperator configuration can be edited as shown below. Pods must be restarted after making this change.
```yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: istio-config
  namespace: istio-system
spec:
  meshConfig:
    defaultConfig:
      tracing:
        sampling: 100
```

## Implementation notes

### Adding tracing to existing services and controllers
The general process is to review the code for any uses of context.Background(). Instead of context.Background(), use a context provided by the caller. Inject the (yet to be added) tracing headers into the outgoing request context.

Then propagate the context up until the incoming server request is received (HTTP, gRPC, etc.). At this point create a derived context that extracts the tracing headers from the request.

For the most part the injection and extraction of tracing headers can be handled by library specific interceptors. This prevents littering the code with inject and extract calls. The common orchestrator packages for mongo and gRPC use interceptors, and tracingMiddleware() in the orchestrator uses the gorilla/mux middleware functionality to intercept incoming HTTP requests. The tracing implementation looks for `zipkin-ip` in config.json for where to post the traces.

Care must be taken when passing the context through to a goroutine. This may result in the context being cancelled in the goroutine when the creator of the goroutine returns. The solution currently used is to create a new context to provide to the goroutine. If the service logs show an unexpected cancel or a trace shows what appear to be orphaned spans, this is likely pointing to an incorrect use of the context.

One last note: if new errors appear in the tests after plumbing the context through then it may be due to the mocks not having the right type signature anymore.

#### Example code flow of tracing through the orchestrator service
Beginning in main.main(), tracing.Middleware() is inserted into the HTTP router. This wraps each API handler with the code needed to setup the tracing context:
```go
	httpRouter := api.NewRouter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	httpRouter.Use(tracing.Middleware)
```

tracing.Middleware() extracts the tracing headers from the request header, creates a span describing the API request, and propagates the context containing the tracing headers to the actual API handler:
```go
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		tracer := otel.Tracer("orchestrator")
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()
		next.ServeHTTP(w, r.Clone(ctx))
```

The actual API handler receives the context from the request and continues passing it down until it eventually nears an exit of the service:
```go
func (h instantiationHandler) approveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iErr := h.client.Approve(ctx, p, ca, v, di)
...
func ...
        result, err := db.DBconn.Find(ctx, c.storeName, key, c.tagMetaData)
```

And finally, the Mongo client has been configured to inject the headers from the passed down context into the outgoing Mongo request using the Monitor client option:
```go
func NewMongoStore(ctx context.Context, name string, store *mongo.Database) (Store, error) {
		clientOptions.Monitor = otelmongo.NewMonitor()
```

This completes propagation of the tracing headers from the incoming API request to the outgoing Mongo, etc. request.
