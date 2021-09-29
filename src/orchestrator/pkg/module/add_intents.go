// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package module

/*
This files deals with the backend implementation of adding
genericPlacementIntents to deployementIntentGroup
*/

import (
	"encoding/json"
	"reflect"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"

	pkgerrors "github.com/pkg/errors"
)

// Intent shall have 2 fields - MetaData and Spec
type Intent struct {
	MetaData IntentMetaData `json:"metadata"`
	Spec     IntentSpecData `json:"spec"`
}

// IntentMetaData has Name, Description, userdata1, userdata2
type IntentMetaData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserData1   string `json:"userData1"`
	UserData2   string `json:"userData2"`
}

// IntentSpecData has Intent
type IntentSpecData struct {
	Intent map[string]string `json:"intent"`
}

// ListOfIntents is a list of intents
type ListOfIntents struct {
	ListOfIntents []map[string]string `json:"intent"`
}

// IntentManager is an interface which exposes the IntentManager functionality
type IntentManager interface {
	AddIntent(a Intent, p string, ca string, v string, di string) (Intent, error)
	GetIntent(i string, p string, ca string, v string, di string) (Intent, error)
	GetAllIntents(p, ca, v, di string) (ListOfIntents, error)
	GetIntentByName(i, p, ca, v, di string) (IntentSpecData, error)
	DeleteIntent(i string, p string, ca string, v string, di string) error
}

// IntentKey consists of Name if the intent, Project name, CompositeApp name,
// CompositeApp version
type IntentKey struct {
	Name                  string `json:"groupIntent"`
	Project               string `json:"project"`
	CompositeApp          string `json:"compositeApp"`
	Version               string `json:"compositeAppVersion"`
	DeploymentIntentGroup string `json:"deploymentIntentGroup"`
}

// We will use json marshalling to convert to string to
// preserve the underlying structure.
func (ik IntentKey) String() string {
	out, err := json.Marshal(ik)
	if err != nil {
		return ""
	}
	return string(out)
}

// IntentClient implements the AddIntentManager interface
type IntentClient struct {
	storeName   string
	tagMetaData string
}

// NewIntentClient returns an instance of AddIntentClient
func NewIntentClient() *IntentClient {
	return &IntentClient{
		storeName:   "resources",
		tagMetaData: "data",
	}
}

/*
AddIntent adds a given intent to the deployment-intent-group and stores in the db.
Other input parameters for it - projectName, compositeAppName, version, DeploymentIntentgroupName
*/
func (c *IntentClient) AddIntent(a Intent, p string, ca string, v string, di string) (Intent, error) {

	//Check for the AddIntent already exists here.
	res, err := c.GetIntent(a.MetaData.Name, p, ca, v, di)
	if !reflect.DeepEqual(res, Intent{}) {
		return Intent{}, pkgerrors.New("Intent already exists")
	}

	akey := IntentKey{
		Name:                  a.MetaData.Name,
		Project:               p,
		CompositeApp:          ca,
		Version:               v,
		DeploymentIntentGroup: di,
	}

	err = db.DBconn.Insert(c.storeName, akey, nil, c.tagMetaData, a)
	if err != nil {
		return Intent{}, pkgerrors.Wrap(err, "Create DB entry error")
	}
	return a, nil
}

/*
GetIntent takes in an IntentName, ProjectName, CompositeAppName, Version and DeploymentIntentGroup.
It returns the Intent.
*/
func (c *IntentClient) GetIntent(i string, p string, ca string, v string, di string) (Intent, error) {

	k := IntentKey{
		Name:                  i,
		Project:               p,
		CompositeApp:          ca,
		Version:               v,
		DeploymentIntentGroup: di,
	}

	result, err := db.DBconn.Find(c.storeName, k, c.tagMetaData)
	if err != nil {
		return Intent{}, err
	}

	if len(result) == 0 {
		return Intent{}, pkgerrors.New("Intent not found")
	}

	if result != nil {
		a := Intent{}
		err = db.DBconn.Unmarshal(result[0], &a)
		if err != nil {
			return Intent{}, err
		}
		return a, nil

	}
	return Intent{}, pkgerrors.New("Unknown Error")
}

/*
GetIntentByName takes in IntentName, projectName, CompositeAppName, CompositeAppVersion
and deploymentIntentGroupName returns the list of intents under the IntentName.
*/
func (c IntentClient) GetIntentByName(i string, p string, ca string, v string, di string) (IntentSpecData, error) {
	k := IntentKey{
		Name:                  i,
		Project:               p,
		CompositeApp:          ca,
		Version:               v,
		DeploymentIntentGroup: di,
	}
	result, err := db.DBconn.Find(c.storeName, k, c.tagMetaData)
	if err != nil {
		return IntentSpecData{}, err
	}
	var a Intent
	err = db.DBconn.Unmarshal(result[0], &a)
	if err != nil {
		return IntentSpecData{}, err
	}
	return a.Spec, nil
}

/*
GetAllIntents takes in projectName, CompositeAppName, CompositeAppVersion,
DeploymentIntentName . It returns ListOfIntents.
*/
func (c IntentClient) GetAllIntents(p string, ca string, v string, di string) (ListOfIntents, error) {
	k := IntentKey{
		Name:                  "",
		Project:               p,
		CompositeApp:          ca,
		Version:               v,
		DeploymentIntentGroup: di,
	}

	result, err := db.DBconn.Find(c.storeName, k, c.tagMetaData)
	if err != nil {
		return ListOfIntents{}, err
	}
	var a Intent
	var listOfMapOfIntents []map[string]string

	if len(result) != 0 {
		for i := range result {
			a = Intent{}
			err = db.DBconn.Unmarshal(result[i], &a)
			if err != nil {
				return ListOfIntents{}, err
			}
			listOfMapOfIntents = append(listOfMapOfIntents, a.Spec.Intent)
		}
		return ListOfIntents{listOfMapOfIntents}, nil
	}
	return ListOfIntents{}, err
}

// DeleteIntent deletes a given intent tied to project, composite app and deployment intent group
func (c IntentClient) DeleteIntent(i string, p string, ca string, v string, di string) error {
	k := IntentKey{
		Name:                  i,
		Project:               p,
		CompositeApp:          ca,
		Version:               v,
		DeploymentIntentGroup: di,
	}

	err := db.DBconn.Remove(c.storeName, k)
	return err
}
