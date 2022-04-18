// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"gitlab.com/project-emco/core/emco-base/src/workflowmgr/api/mocks"
	tmpl "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/emcotemporalapi"
	moduleLib "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/module"
	cl "go.temporal.io/sdk/client"
)

var _ = Describe("workflowIntenthandler", func() {

	type testCase struct {
		inputName    string
		inputReader  io.Reader
		inStruct     moduleLib.WorkflowIntent
		mockError    error
		mockVal      moduleLib.WorkflowIntent
		mockVals     []moduleLib.WorkflowIntent
		expectedCode int
		client       *mocks.WorkflowIntentManager
	}

	DescribeTable("Create workflowIntent tests",
		func(t testCase) {
			// set up client mock responses
			t.client.On("CreateWorkflowIntent", t.inStruct, "test-project", "test-compositeapp", "v1", "test-dig", false).Return(t.mockVal, t.mockError)

			// make HTTP request
			request := httptest.NewRequest("POST", "/v2/projects/test-project/composite-apps/test-compositeapp/v1/deployment-intent-groups/test-dig/temporal-workflow-intents", t.inputReader)
			resp := executeRequest(request, NewRouter(t.client))

			//Check returned code
			Expect(resp.StatusCode).To(Equal(t.expectedCode))

			//Check returned body
			got := moduleLib.WorkflowIntent{}
			json.NewDecoder(resp.Body).Decode(&got)
			Expect(got).To(Equal(t.mockVal))
		},

		Entry("successful create", testCase{
			expectedCode: http.StatusCreated,
			inputReader: bytes.NewBuffer([]byte(`{
				"metadata": {
					"name": "testworkflowintent",
					"description": "test workflow intent",
					"userData1": "some user data 1",
					"userData2": "some user data 2"
				},
				"spec": {
					"workflowClient": {
					"clientEndpointName": "endpoint1",
					"clientEndpointPort": 9091
					},
					"temporal": {
						"workflowClientName": "client1",
						"workflowStartOptions": {
							"id":  "xyz",
							"taskQueue": "my-task-queue"
						}
					}
				}
			}`)),
			inStruct: moduleLib.WorkflowIntent{
				Metadata: moduleLib.Metadata{
					Name:        "testworkflowintent",
					Description: "test workflow intent",
					UserData1:   "some user data 1",
					UserData2:   "some user data 2",
				},
				Spec: moduleLib.WorkflowIntentSpec{
					WfClientSpec: moduleLib.WfClientSpec{
						WfClientEndpointName: "endpoint1",
						WfClientEndpointPort: 9091,
					},
					WfTemporalSpec: tmpl.WfTemporalSpec{
						WfClientName: "client1",
						WfStartOpts: cl.StartWorkflowOptions{
							ID:        "xyz",
							TaskQueue: "my-task-queue",
						},
					},
				},
			},
			mockError: nil,
			mockVal: moduleLib.WorkflowIntent{
				Metadata: moduleLib.Metadata{
					Name:        "testworkflowintent",
					Description: "test workflow intent",
					UserData1:   "some user data 1",
					UserData2:   "some user data 2",
				},
				Spec: moduleLib.WorkflowIntentSpec{
					WfClientSpec: moduleLib.WfClientSpec{
						WfClientEndpointName: "endpoint1",
						WfClientEndpointPort: 9091,
					},
					WfTemporalSpec: tmpl.WfTemporalSpec{
						WfClientName: "client1",
						WfStartOpts: cl.StartWorkflowOptions{
							ID:        "xyz",
							TaskQueue: "my-task-queue",
						},
					},
				},
			},
			client: &mocks.WorkflowIntentManager{},
		}),

		Entry("fails due to empty body", testCase{
			expectedCode: http.StatusBadRequest,
			inStruct:     moduleLib.WorkflowIntent{},
			mockError:    nil,
			mockVal:      moduleLib.WorkflowIntent{},
			client:       &mocks.WorkflowIntentManager{},
		}),
	)
})
