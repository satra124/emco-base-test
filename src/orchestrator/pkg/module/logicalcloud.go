// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

// FIXME:
// Temporarily copied from DCM as of 20201102 to continue development.
// Updated on 20201125 to reflect code changes in DCM.
// Need solution to cyclic dependency: emcolib?

package module

import (
	"encoding/json"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/appcontext"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	pkgerrors "github.com/pkg/errors"
)

// LogicalCloud contains the parameters needed for a Logical Cloud
type LogicalCloud struct {
	MetaData      MetaDataList `json:"metadata"`
	Specification Spec         `json:"spec"`
}

// MetaData contains the parameters needed for metadata
type MetaDataList struct {
	LogicalCloudName string `json:"name"`
	Description      string `json:"description"`
	UserData1        string `json:"userData1"`
	UserData2        string `json:"userData2"`
}

// Spec contains the parameters needed for spec
type Spec struct {
	NameSpace string   `json:"namespace"`
	Level     string   `json:"level"`
	User      UserData `json:"user"`
}

// UserData contains the parameters needed for user
type UserData struct {
	UserName        string     `json:"user-name"`
	Type            string     `json:"type"`
	UserPermissions []UserPerm `json:"user-permissions"`
}

//  UserPerm contains the parameters needed for user permissions
type UserPerm struct {
	PermName  string   `json:"permission-name"`
	APIGroups []string `json:"apiGroups"`
	Resources []string `json:"resources"`
	Verbs     []string `json:"verbs"`
}

// LogicalCloudKey is the key structure that is used in the database
type LogicalCloudKey struct {
	Project          string `json:"project"`
	LogicalCloudName string `json:"logicalCloud"`
}

// AppContextKey is an alternative key to access logical clouds
type AppContextKey struct {
	LCContext string `json:"logicalCloudContext"`
}

// LogicalCloudManager is an interface that exposes the connection
// functionality
type LogicalCloudManager interface {
	Create(project string, c LogicalCloud) (LogicalCloud, error)
	Get(project, name string) (LogicalCloud, error)
	GetAll(project string) ([]LogicalCloud, error)
	Delete(project, name string) error
	Update(project, name string, c LogicalCloud) (LogicalCloud, error)
}

// Interface facilitates unit testing by mocking functions
type Utility interface {
	DBInsert(storeName string, key db.Key, query interface{}, meta string, c interface{}) error
	DBFind(storeName string, key db.Key, meta string) ([][]byte, error)
	DBUnmarshal(value []byte, out interface{}) error
	DBRemove(storeName string, key db.Key) error
	CheckProject(project string) error
	CheckLogicalCloud(project, logicalCloud string) error
	GetLogicalCloudContext(storeName string, key db.Key, meta string, project string, name string) (appcontext.AppContext, string, error)
	GetLogicalCloudFromContext(storeName string, appContextId string) (string, string, error)
	GetAppContextStatus(ac appcontext.AppContext) (*appcontext.AppContextStatus, error)
}

// LogicalCloudClient implements the LogicalCloudManager
// It will also be used to maintain some localized state
type LogicalCloudClient struct {
	storeName  string
	tagMeta    string
	tagContext string
	util       Utility
}

// Added for unit testing; implements Utility interface
type DBService struct{}

// LogicalCloudClient returns an instance of the LogicalCloudClient
// which implements the LogicalCloudManager
func NewLogicalCloudClient() *LogicalCloudClient {
	service := DBService{}
	return &LogicalCloudClient{
		storeName:  "resources",
		tagMeta:    "data",
		tagContext: "logicalCloudContext",
		util:       service,
	}
}

// Create entry for the logical cloud resource in the database
func (v *LogicalCloudClient) Create(project string, c LogicalCloud) (LogicalCloud, error) {

	//Construct key consisting of name
	key := LogicalCloudKey{
		Project:          project,
		LogicalCloudName: c.MetaData.LogicalCloudName,
	}

	//Check if project exists
	err := v.util.CheckProject(project)
	if err != nil {
		return LogicalCloud{}, err
	}

	//Check if this Logical Cloud already exists
	_, err = v.Get(project, c.MetaData.LogicalCloudName)
	if err == nil {
		return LogicalCloud{}, pkgerrors.New("Logical Cloud already exists")
	}

	// if Logical Cloud Level is not specified, it defaults to 1:
	if c.Specification.Level == "" {
		c.Specification.Level = "1"
	}

	err = v.util.DBInsert(v.storeName, key, nil, v.tagMeta, c)
	if err != nil {
		return LogicalCloud{}, pkgerrors.Wrap(err, "Create DB entry error")
	}

	return c, nil
}

// Get returns Logical Cloud corresponding to logical cloud name
func (v *LogicalCloudClient) Get(project, logicalCloudName string) (LogicalCloud, error) {

	//Construct the composite key to select the entry
	key := LogicalCloudKey{
		Project:          project,
		LogicalCloudName: logicalCloudName,
	}
	value, err := v.util.DBFind(v.storeName, key, v.tagMeta)
	if err != nil {
		return LogicalCloud{}, err
	}

	if len(value) == 0 {
		return LogicalCloud{}, pkgerrors.New("LogicalCloud not found")
	}

	//value is a byte array
	if value != nil {
		lc := LogicalCloud{}
		err = v.util.DBUnmarshal(value[0], &lc)
		if err != nil {
			return LogicalCloud{}, pkgerrors.Wrap(err, "Error unmarshalling value")
		}
		return lc, nil
	}

	return LogicalCloud{}, pkgerrors.New("Unknown Error")
}

// GetAll returns Logical Clouds in the project
func (v *LogicalCloudClient) GetAll(project string) ([]LogicalCloud, error) {

	//Construct the composite key to select the entry
	key := LogicalCloudKey{
		Project:          project,
		LogicalCloudName: "",
	}

	var resp []LogicalCloud
	values, err := v.util.DBFind(v.storeName, key, v.tagMeta)
	if err != nil {
		return []LogicalCloud{}, err
	}

	for _, value := range values {
		lc := LogicalCloud{}
		err = v.util.DBUnmarshal(value, &lc)
		if err != nil {
			return []LogicalCloud{}, err
		}
		resp = append(resp, lc)
	}

	return resp, nil
}

// Delete the Logical Cloud entry from database
func (v *LogicalCloudClient) Delete(project, logicalCloudName string) error {

	//Construct the composite key to select the entry
	key := LogicalCloudKey{
		Project:          project,
		LogicalCloudName: logicalCloudName,
	}
	//Check if this Logical Cloud exists
	_, err := v.Get(project, logicalCloudName)
	if err != nil {
		return err
	}

	context, _, err := v.util.GetLogicalCloudContext(v.storeName, key, v.tagContext, project, logicalCloudName)
	// If there's no context for Logical Cloud, just go ahead and delete it now
	if err != nil {
		err = v.util.DBRemove(v.storeName, key)
		if err != nil {
			return pkgerrors.Wrap(err, "Error when deleting Logical Cloud (scenario with no context)")
		}
		return nil
	}

	// Make sure rsync status for this logical cloud is Terminated,
	// otherwise we can't remove appcontext yet
	acStatus, err := v.util.GetAppContextStatus(context)
	if err != nil {
		return err
	}
	switch acStatus.Status {
	case appcontext.AppContextStatusEnum.Terminated:
		// remove the appcontext
		err := context.DeleteCompositeApp()
		if err != nil {
			log.Error("Error deleting AppContext CompositeApp Logical Cloud.", log.Fields{"logicalcloud": logicalCloudName})
			return pkgerrors.Wrap(err, "Error deleting AppContext CompositeApp Logical Cloud.")
		}

		err = v.util.DBRemove(v.storeName, key)
		if err != nil {
			log.Error("Error when deleting Logical Cloud (scenario with Terminated status).", log.Fields{"logicalcloud": logicalCloudName})
			return pkgerrors.Wrap(err, "Error when deleting Logical Cloud (scenario with Terminated status).")
		}
		log.Info("Deleted Logical Cloud.", log.Fields{"logicalcloud": logicalCloudName})
		return nil
	case appcontext.AppContextStatusEnum.Terminating:
		log.Error("The Logical Cloud can't be deleted yet, it is being terminated.", log.Fields{"logicalcloud": logicalCloudName})
		return pkgerrors.New("The Logical Cloud can't be deleted yet, it is being terminated.")
	case appcontext.AppContextStatusEnum.Instantiated:
		log.Error("The Logical Cloud is instantiated, please terminate first.", log.Fields{"logicalcloud": logicalCloudName})
		return pkgerrors.New("The Logical Cloud is instantiated, please terminate first.")
	case appcontext.AppContextStatusEnum.Instantiating:
		log.Error("The Logical Cloud is instantiating, please wait and then terminate.", log.Fields{"logicalcloud": logicalCloudName})
		return pkgerrors.New("The Logical Cloud is instantiating, please wait and then terminate.")
	case appcontext.AppContextStatusEnum.TerminateFailed:
		log.Error("The Logical Cloud has failed terminating, please correct the situation and try again.", log.Fields{"logicalcloud": logicalCloudName})
		return pkgerrors.New("The Logical Cloud has failed terminating, please correct the situation and try again.")
	case appcontext.AppContextStatusEnum.InstantiateFailed:
		log.Error("The Logical Cloud has failed instantiating, for safety please terminate and try again.", log.Fields{"logicalcloud": logicalCloudName})
		return pkgerrors.New("The Logical Cloud has failed instantiating, for safety please terminate and try again.")
	default:
		log.Error("The Logical Cloud isn't in an expected status so not taking any action.", log.Fields{"logicalcloud": logicalCloudName, "status": acStatus.Status})
		return pkgerrors.New("The Logical Cloud isn't in an expected status so not taking any action.")
	}
}

// Update an entry for the Logical Cloud in the database
func (v *LogicalCloudClient) Update(project, logicalCloudName string, c LogicalCloud) (LogicalCloud, error) {

	key := LogicalCloudKey{
		Project:          project,
		LogicalCloudName: logicalCloudName,
	}
	// Check for mismatch, logicalCloudName and payload logical cloud name
	if c.MetaData.LogicalCloudName != logicalCloudName {
		return LogicalCloud{}, pkgerrors.New("Logical Cloud name mismatch")
	}
	//Check if this Logical Cloud exists
	_, err := v.Get(project, logicalCloudName)
	if err != nil {
		return LogicalCloud{}, err
	}
	err = v.util.DBInsert(v.storeName, key, nil, v.tagMeta, c)
	if err != nil {
		return LogicalCloud{}, pkgerrors.Wrap(err, "Updating DB Entry")
	}
	return c, nil
}

// GetLogicalCloudContext returns the AppContext for corresponding provider and name
func (d DBService) GetLogicalCloudContext(storeName string, key db.Key, meta string, project string, name string) (appcontext.AppContext, string, error) {

	value, err := d.DBFind(storeName, key, meta)
	if err != nil {
		return appcontext.AppContext{}, "", err
	}

	if len(value) == 0 {
		return appcontext.AppContext{}, "", pkgerrors.New("Logical Cloud AppContext not found")
	}

	//value is a [][]byte
	if value != nil {
		ctxVal := string(value[0])
		var lcc appcontext.AppContext
		_, err = lcc.LoadAppContext(ctxVal)
		if err != nil {
			return appcontext.AppContext{}, "", pkgerrors.Wrap(err, "Reinitializing Logical Cloud AppContext")
		}
		return lcc, ctxVal, nil
	}

	return appcontext.AppContext{}, "", pkgerrors.New("Unknown Error")
}

// GetLogicalCloudFromContext returns the pair (project, logical cloud name) for a given AppContext ID
func (d DBService) GetLogicalCloudFromContext(storeName string, appContextId string) (string, string, error) {
	key := AppContextKey{
		LCContext: appContextId,
	}
	log.Info("GetLogicalCloudFromContext", log.Fields{"appContextId": appContextId})

	values, err := d.DBFind(storeName, key, "logicalCloud")
	if err != nil {
		log.Error("Couldn't fetch logical cloud", log.Fields{"err": err})
		return "", "", err
	}

	if len(values) == 0 {
		return "", "", pkgerrors.New("LogicalCloud not found")
	}

	logicalCloudName := string(values[0])
	log.Info("", log.Fields{"logicalCloudName": logicalCloudName})

	values, err = d.DBFind(storeName, key, "project")
	if err != nil {
		log.Error("Couldn't fetch project", log.Fields{"err": err})
		return "", "", err
	}

	if len(values) == 0 {
		return "", "", pkgerrors.New("Project not found")
	}

	project := string(values[0])
	log.Info("", log.Fields{"project": project})

	return project, logicalCloudName, nil
}

func (d DBService) DBInsert(storeName string, key db.Key, query interface{}, meta string, c interface{}) error {

	err := db.DBconn.Insert(storeName, key, nil, meta, c)
	if err != nil {
		return pkgerrors.Wrap(err, "Create DB entry error")
	}

	return nil
}

func (d DBService) DBFind(storeName string, key db.Key, meta string) ([][]byte, error) {

	value, err := db.DBconn.Find(storeName, key, meta)
	if err != nil {
		return [][]byte{}, err
	}

	return value, nil
}

func (d DBService) DBUnmarshal(value []byte, out interface{}) error {

	err := db.DBconn.Unmarshal(value, out)
	if err != nil {
		return err
	}

	return nil
}

func (d DBService) DBRemove(storeName string, key db.Key) error {

	err := db.DBconn.Remove(storeName, key)
	return err
}

func (d DBService) CheckProject(project string) error {
	// Check if project exists
	_, err := NewProjectClient().GetProject(project)
	return err
}

func (d DBService) CheckLogicalCloud(project, logicalCloud string) error {
	// Check if logical cloud exists
	_, err := NewLogicalCloudClient().Get(project, logicalCloud)
	return err
}

func (d DBService) GetAppContextStatus(ac appcontext.AppContext) (*appcontext.AppContextStatus, error) {

	h, err := ac.GetCompositeAppHandle()
	if err != nil {
		return nil, err
	}
	sh, err := ac.GetLevelHandle(h, "status")
	if err != nil {
		return nil, err
	}
	s, err := ac.GetValue(sh)
	if err != nil {
		return nil, err
	}
	acStatus := appcontext.AppContextStatus{}
	js, _ := json.Marshal(s)
	json.Unmarshal(js, &acStatus)

	return &acStatus, nil

}
