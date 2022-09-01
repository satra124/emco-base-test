// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Aarna Networks

package metrics

import (
	"context"
	"os"
	"strconv"
	"time"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	orchModule "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/ovnaction/pkg/module"
	"go.opentelemetry.io/otel"
)

func Start() {
	pollingInterval, _ := strconv.Atoi(os.Getenv("EMCO_METRICS_POLLING_INTERVAL_SECS"))
	if pollingInterval == -1 {
		return
	}
	if pollingInterval == 0 {
		pollingInterval = 15
	}
	go func() {
		tracer := otel.Tracer("ovnaction")
		for {
			ctx, span := tracer.Start(context.Background(), "get-metrics")
			do(ctx)
			span.End()
			time.Sleep(time.Duration(time.Duration(pollingInterval) * time.Second))
		}
	}()
}

func do(ctx context.Context) {
	orchClient := orchModule.NewClient()
	client := module.NewClient()
	fields := log.Fields{"service": "ovnaction"}
	projects, err := orchModule.NewProjectClient().GetAllProjects(ctx)
	if err != nil {
		log.Error(err.Error(), fields)
		return
	}

	for _, proj := range projects {
		fields := fields
		fields["project"] = proj.MetaData.Name

		apps, err := orchClient.CompositeApp.GetAllCompositeApps(ctx, proj.MetaData.Name)
		if err != nil {
			log.Error(err.Error(), fields)
			continue
		}
		for _, app := range apps {
			fields := fields
			fields["composite_app"] = app.Metadata.Name

			digs, err := orchClient.DeploymentIntentGroup.GetAllDeploymentIntentGroups(ctx, proj.MetaData.Name, app.Metadata.Name, app.Spec.Version)
			if err != nil {
				log.Error(err.Error(), fields)
				continue
			}
			for _, dig := range digs {
				fields := fields
				fields["dig"] = dig.MetaData.Name

				ncis, err := client.NetControlIntent.GetNetControlIntents(proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name)
				if err != nil {
					log.Error(err.Error(), fields)
					continue
				}
				for _, nci := range ncis {
					fields := fields
					fields["net_controller_intent"] = nci.Metadata.Name
					NetworkControllerIntentGauge.WithLabelValues(nci.Metadata.Name, proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name).Set(1)
					wis, err := client.WorkloadIntent.GetWorkloadIntents(proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name, nci.Metadata.Name)
					if err != nil {
						log.Error(err.Error(), fields)
						continue
					}
					for _, wi := range wis {
						fields := fields
						fields["workload_intent"] = wi.Metadata.Name
						WorkloadIntentGauge.WithLabelValues(wi.Metadata.Name, proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name, nci.Metadata.Name, wi.Spec.AppName, wi.Spec.WorkloadResource, wi.Spec.Type).Set(1)
						wiifs, err := client.WorkloadIfIntent.GetWorkloadIfIntents(proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name, nci.Metadata.Name, wi.Metadata.Name)
						if err != nil {
							log.Error(err.Error(), fields)
							continue
						}
						for _, wiif := range wiifs {
							WorkloadInterfaceIntentGauge.WithLabelValues(
								wiif.Metadata.Name,
								proj.MetaData.Name, app.Metadata.Name, app.Spec.Version, dig.MetaData.Name, nci.Metadata.Name, wi.Metadata.Name,
								wiif.Spec.IfName,
								wiif.Spec.NetworkName,
								wiif.Spec.DefaultGateway,
								wiif.Spec.IpAddr,
								wiif.Spec.MacAddr,
							).Set(1)
						}
					}
				}
			}
		}
	}
}