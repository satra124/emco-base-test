// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package module

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
)

func TestCreateDeploymentIntentGroup(t *testing.T) {
	testCases := []struct {
		label                    string
		inputDeploymentIntentGrp DeploymentIntentGroup
		inputProject             string
		inputCompositeApp        string
		inputCompositeAppVersion string
		expectedError            string
		mockdb                   *db.MockDB
		expected                 DeploymentIntentGroup
	}{
		{
			label: "Create DeploymentIntentGroup",
			inputDeploymentIntentGrp: DeploymentIntentGroup{
				MetaData: DepMetaData{
					Name:        "testDeploymentIntentGroup",
					Description: "DescriptionTestDeploymentIntentGroup",
					UserData1:   "userData1",
					UserData2:   "userData2",
				},
				Spec: DepSpecData{
					Profile: "Testprofile",
					Version: "version of deployment",
					OverrideValuesObj: []OverrideValues{
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
					},
					LogicalCloud: "cloud1",
				},
			},
			inputProject:             "testProject",
			inputCompositeApp:        "testCompositeApp",
			inputCompositeAppVersion: "testCompositeAppVersion",
			expected: DeploymentIntentGroup{
				MetaData: DepMetaData{
					Name:        "testDeploymentIntentGroup",
					Description: "DescriptionTestDeploymentIntentGroup",
					UserData1:   "userData1",
					UserData2:   "userData2",
				},
				Spec: DepSpecData{
					Profile: "Testprofile",
					Version: "version of deployment",
					OverrideValuesObj: []OverrideValues{
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
					},
					LogicalCloud: "cloud1",
				},
			},
			expectedError: "",
			mockdb: &db.MockDB{
				Items: []map[string]map[string][]byte{
					{
						ProjectKey{ProjectName: "testProject"}.String(): {
							"data": []byte(
								"{\"project-name\":\"testProject\"," +
									"\"description\":\"Test project for unit testing\"}"),
						},
						CompositeAppKey{CompositeAppName: "testCompositeApp",
							Version: "testCompositeAppVersion", Project: "testProject"}.String(): {
							"data": []byte(
								"{\"metadata\":{" +
									"\"name\":\"testCompositeApp\"," +
									"\"description\":\"description\"," +
									"\"userData1\":\"user data\"," +
									"\"userData2\":\"user data\"" +
									"}," +
									"\"spec\":{" +
									"\"compositeAppVersion\":\"version of the composite app\"}}"),
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			db.DBconn = testCase.mockdb
			depIntentCli := NewDeploymentIntentGroupClient()
			got, err := depIntentCli.CreateDeploymentIntentGroup(testCase.inputDeploymentIntentGrp, testCase.inputProject, testCase.inputCompositeApp, testCase.inputCompositeAppVersion)
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("CreateDeploymentIntentGroup returned an unexpected error %s, ", err)
				}
				if strings.Contains(err.Error(), testCase.expectedError) == false {
					t.Fatalf("CreateDeploymentIntentGroup returned an unexpected error %s", err)
				}
			} else {
				if reflect.DeepEqual(testCase.expected, got) == false {
					t.Errorf("CreateDeploymentIntentGroup returned unexpected body: got %v; "+" expected %v", got, testCase.expected)
				}
			}
		})
	}
}

func TestGetDeploymentIntentGroup(t *testing.T) {
	testCases := []struct {
		label                    string
		inputDeploymentIntentGrp string
		inputProject             string
		inputCompositeApp        string
		inputCompositeAppVersion string
		expected                 DeploymentIntentGroup
		expectedError            string
		mockdb                   *db.MockDB
	}{
		{
			label:                    "Get DeploymentIntentGroup",
			inputDeploymentIntentGrp: "testDeploymentIntentGroup",
			inputProject:             "testProject",
			inputCompositeApp:        "testCompositeApp",
			inputCompositeAppVersion: "testCompositeAppVersion",
			expected: DeploymentIntentGroup{
				MetaData: DepMetaData{
					Name:        "testDeploymentIntentGroup",
					Description: "DescriptionTestDeploymentIntentGroup",
					UserData1:   "userData1",
					UserData2:   "userData2",
				},
				Spec: DepSpecData{
					Profile: "Testprofile",
					Version: "version of deployment",
					OverrideValuesObj: []OverrideValues{
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
						{AppName: "TestAppName",
							ValuesObj: map[string]string{
								"imageRepository": "registry.hub.docker.com",
							}},
					},
					LogicalCloud: "cloud1",
				},
			},
			expectedError: "",
			mockdb: &db.MockDB{
				Items: []map[string]map[string][]byte{
					{
						DeploymentIntentGroupKey{
							Name:         "testDeploymentIntentGroup",
							Project:      "testProject",
							CompositeApp: "testCompositeApp",
							Version:      "testCompositeAppVersion",
						}.String(): {
							"data": []byte(
								"{\"metadata\":{\"name\":\"testDeploymentIntentGroup\"," +
									"\"description\":\"DescriptionTestDeploymentIntentGroup\"," +
									"\"userData1\": \"userData1\"," +
									"\"userData2\": \"userData2\"}," +
									"\"spec\":{\"compositeProfile\": \"Testprofile\"," +
									"\"version\": \"version of deployment\"," +
									"\"overrideValues\":[" +
									"{" +
									"\"app\": \"TestAppName\"," +
									"\"values\": " +
									"{" +
									"\"imageRepository\":\"registry.hub.docker.com\"" +
									"}" +
									"}," +
									"{" +
									"\"app\": \"TestAppName\"," +
									"\"values\": " +
									"{" +
									"\"imageRepository\":\"registry.hub.docker.com\"" +
									"}" +
									"}" +
									"]," +
									"\"logicalCloud\": \"cloud1\"" +
									"}" +
									"}"),
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			db.DBconn = testCase.mockdb
			depIntentCli := NewDeploymentIntentGroupClient()
			got, err := depIntentCli.GetDeploymentIntentGroup(testCase.inputDeploymentIntentGrp, testCase.inputProject, testCase.inputCompositeApp, testCase.inputCompositeAppVersion)
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("GetDeploymentIntentGroup returned an unexpected error: %s", err)
				}
				if strings.Contains(err.Error(), testCase.expectedError) == false {
					t.Fatalf("GetDeploymentIntentGroup returned an unexpected error: %s", err)
				}
			} else {
				if reflect.DeepEqual(testCase.expected, got) == false {
					t.Errorf("GetDeploymentIntentGroup returned unexpected body: got %v;"+
						" expected %v", got, testCase.expected)
				}
			}
		})
	}
}
