// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package module_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/project-emco/core/emco-base/src/dtc/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	pkgerrors "github.com/pkg/errors"
)

var _ = Describe("Trafficgroupintent", func() {

	var (
		TGI    module.TrafficGroupIntent
		TGIDBC *module.TrafficGroupIntentDbClient
		mdb    *db.MockDB
	)

	BeforeEach(func() {
		TGIDBC = module.NewTrafficGroupIntentClient()
		TGI = module.TrafficGroupIntent{
			Metadata: module.Metadata{
				Name:        "testtgi",
				Description: "traffic group intent",
				UserData1:   "user data1",
				UserData2:   "user data2",
			},
		}
		mdb = new(db.MockDB)
		mdb.Err = nil
		db.DBconn = mdb

	})

	Describe("Create traffic intent", func() {
		It("should return nil", func() {
			_, err := (*TGIDBC).CreateTrafficGroupIntent(TGI, "test", "capp1", "v1", "dig", false)
			Expect(err).To(BeNil())
		})
		It("followed by create again should return error", func() {
			_, err := (*TGIDBC).CreateTrafficGroupIntent(TGI, "test", "capp1", "v1", "dig", false)
			Expect(err).To(BeNil())
			_, err = (*TGIDBC).CreateTrafficGroupIntent(TGI, "test", "capp1", "v1", "dig", false)
			Expect(err).To(HaveOccurred())
		})

		It("followed by get should return nil", func() {
			_, err := (*TGIDBC).CreateTrafficGroupIntent(TGI, "test", "capp1", "v1", "dig", false)
			Expect(err).To(BeNil())
			tgi, err := (*TGIDBC).GetTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(BeNil())
			Expect(tgi).Should(Equal(TGI))
		})
		It("followed by delete should return nil", func() {
			_, err := (*TGIDBC).CreateTrafficGroupIntent(TGI, "test", "capp1", "v1", "dig", false)
			Expect(err).To(BeNil())
			err = (*TGIDBC).DeleteTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(BeNil())
		})

	})
	Describe("Get traffic intent", func() {
		It("should return error for non-existing record", func() {
			_, err := (*TGIDBC).GetTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(HaveOccurred())
		})

	})
	Describe("Get traffic intents", func() {
		It("should return error for non-existing record", func() {
			mdb.Err = pkgerrors.New("Traffic group intent not found")
			_, err := (*TGIDBC).GetTrafficGroupIntents("test", "capp1", "v1", "dig")
			Expect(err).To(HaveOccurred())
		})

	})
	Describe("Delete traffic intent", func() {
		It("should return error for non-existing record", func() {
			mdb.Err = pkgerrors.New("db Remove resource not found")
			err := (*TGIDBC).DeleteTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(HaveOccurred())
		})
		It("should return error for deleting parent without deleting child", func() {
			mdb.Err = pkgerrors.New("Cannot delete parent without deleting child references first")
			err := (*TGIDBC).DeleteTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(HaveOccurred())
		})
		It("should return error for general db error", func() {
			mdb.Err = pkgerrors.New("db Remove error")
			err := (*TGIDBC).DeleteTrafficGroupIntent("testtgi", "test", "capp1", "v1", "dig")
			Expect(err).To(HaveOccurred())
		})

	})
})
