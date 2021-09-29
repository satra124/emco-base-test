// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package module_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	orch "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/sfc/pkg/model"
	"gitlab.com/project-emco/core/emco-base/src/sfc/pkg/module"
)

var _ = Describe("SFCClientSelectorIntent", func() {
	var (
		proj       orch.Project
		projClient *orch.ProjectClient

		ca       orch.CompositeApp
		caClient *orch.CompositeAppClient

		dig       orch.DeploymentIntentGroup
		digClient *orch.DeploymentIntentGroupClient

		sfcIntent               model.SfcIntent
		sfcIntentClient         *module.SfcIntentClient
		sfcClientSelectorIntent model.SfcClientSelectorIntent
		sfcClient               *module.SfcClientSelectorIntentClient

		mdb *db.NewMockDB
	)

	BeforeEach(func() {
		projClient = orch.NewProjectClient()
		proj = orch.Project{
			MetaData: orch.ProjectMetaData{
				Name: "testproject",
			},
		}

		caClient = orch.NewCompositeAppClient()
		ca = orch.CompositeApp{
			Metadata: orch.CompositeAppMetaData{
				Name: "ca",
			},
			Spec: orch.CompositeAppSpec{
				Version: "v1",
			},
		}

		digClient = orch.NewDeploymentIntentGroupClient()
		dig = orch.DeploymentIntentGroup{
			MetaData: orch.DepMetaData{
				Name: "dig",
			},
			Spec: orch.DepSpecData{
				Profile:      "profilename",
				Version:      "testver",
				LogicalCloud: "logCloud",
			},
		}

		sfcIntentClient = module.NewSfcIntentClient()
		sfcIntent = model.SfcIntent{
			Metadata: model.Metadata{
				Name: "sfcIntentName",
			},
		}

		sfcClient = module.NewSfcClientSelectorIntentClient()
		sfcClientSelectorIntent = model.SfcClientSelectorIntent{
			Metadata: model.Metadata{
				Name: "sfcClientSelectorIntentName",
			},
		}

		mdb = new(db.NewMockDB)
		mdb.Err = nil
		db.DBconn = mdb
	})

	Describe("Create SFC client selector intent", func() {
		It("successful creation of sfc client selector intent", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC Client selector intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())
		})
		It("followed by create again should return error", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())
			// test SFC intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(strings.Contains(err.Error(), "SFC Client Selector Intent already exists")).To(Equal(true))
		})
		It("successful creation of sfc client selector intent with update version of call", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC Client Selector intent creation, with update form of call (exists bool == true)
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", true)
			Expect(err).To(BeNil())
		})
		It("successful creation of sfc client selector intent with update version of call", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC Client Selector intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", true)
			Expect(err).To(BeNil())
		})
	})

	Describe("Get all sfc client selector intents", func() {
		It("Parent SFC Intent does not exist - return not found error", func() {
			_, err := (*sfcClient).GetAllSfcClientSelectorIntents("testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "not found")).To(Equal(true))
		})
		It("Parent SFC Intent does exist - No SFC Client Selector Intents - should return empty list", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			list, err := (*sfcClient).GetAllSfcClientSelectorIntents("testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(err).To(BeNil())
			Expect(len(list)).To(Equal(0))
		})
		It("Parent SFC Intent does exist - 2 SFC Intents created - should return list of len 2", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// SFC client selector intent creation - make 2 of them
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", true)
			Expect(err).To(BeNil())
			sfcClientSelectorIntent.Metadata.Name = "2nd_name"
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", true)
			Expect(err).To(BeNil())

			list, err := (*sfcClient).GetAllSfcClientSelectorIntents("testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(len(list)).To(Equal(2))

		})
		It("should return error for general db error", func() {
			mdb.Err = pkgerrors.New("db Find error")
			_, err := (*sfcClient).GetAllSfcClientSelectorIntents("testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "db Find error")).To(Equal(true))
		})
		It("should return error for unmarshalling db error", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			mdb.MarshalErr = pkgerrors.New("Unmarshalling bson")
			_, err = (*sfcClient).GetAllSfcClientSelectorIntents("testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "Unmarshalling bson")).To(Equal(true))
		})
	})

	Describe("Get sfc client selector intent", func() {
		It("Parent SFC Intent does not exist - return not found error", func() {
			_, err := (*sfcClient).GetSfcClientSelectorIntent("sfcIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "not found")).To(Equal(true))
		})
		It("Successful get of sfcClientSelectorIntent", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())

			_, err = (*sfcClient).GetSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(err).To(BeNil())
		})
		It("should return error for general db error", func() {
			mdb.Err = pkgerrors.New("db Find error")
			_, err := (*sfcClient).GetSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "db Find error")).To(Equal(true))
		})
		It("should return error for unmarshalling db error", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())
			mdb.MarshalErr = pkgerrors.New("Unmarshalling bson")
			_, err = (*sfcClient).GetSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "Unmarshalling bson")).To(Equal(true))
		})
	})

	Describe("Delete SFC Client Selector intent", func() {
		It("successful delete", func() {
			// set up prerequisites
			_, err := (*projClient).CreateProject(proj, false)
			Expect(err).To(BeNil())
			_, err = (*caClient).CreateCompositeApp(ca, "testproject", false)
			Expect(err).To(BeNil())
			_, err = (*digClient).CreateDeploymentIntentGroup(dig, "testproject", "ca", "v1")
			Expect(err).To(BeNil())
			_, err = (*sfcIntentClient).CreateSfcIntent(sfcIntent, "testproject", "ca", "v1", "dig", false)
			Expect(err).To(BeNil())

			// test SFC intent creation
			_, err = (*sfcClient).CreateSfcClientSelectorIntent(sfcClientSelectorIntent, "testproject", "ca", "v1", "dig", "sfcIntentName", false)
			Expect(err).To(BeNil())

			err = (*sfcClient).DeleteSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(err).To(BeNil())
		})
		It("should return not found error for non-existing record", func() {
			mdb.Err = pkgerrors.New("db Remove resource not found")
			err := (*sfcClient).DeleteSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "db Remove resource not found")).To(Equal(true))
		})
		It("should return error for deleting parent without deleting child", func() {
			mdb.Err = pkgerrors.New("db Remove parent child constraint")
			err := (*sfcClient).DeleteSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "db Remove parent child constraint")).To(Equal(true))
		})
		It("should return error for general db error", func() {
			mdb.Err = pkgerrors.New("db Remove error")
			err := (*sfcClient).DeleteSfcClientSelectorIntent("sfcClientSelectorIntentName", "testproject", "ca", "v1", "dig", "sfcIntentName")
			Expect(strings.Contains(err.Error(), "db Remove error")).To(Equal(true))
		})
	})
})
