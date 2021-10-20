// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package depend

import (
	"context"
	"reflect"
	"sync"
	"time"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/types"
)

type DependManager struct {
	acID string
	// Per App Ready channels to notify
	readyCh map[string][]appData
	// Per App Deploy channels to notify
	deployedCh map[string][]appData
	// Per app channels to wait on
	appCh map[string][]chan struct{}
	// Single Resource succeed channels
	resCh map[string][]resData
	sync.Mutex
}

type resData struct {
	res string
	// Chan to report if app meets Criteria
	ch chan struct{}
}

type appData struct {
	app string
	crt types.Criteria
	// Chan to report if app meets Criteria
	ch chan struct{}
}

var dmList map[string]*DependManager

func init() {
	dmList = make(map[string]*DependManager)
}

const SEPARATOR = "+"

// New Manager for acID
func NewDependManager(acID string) *DependManager {
	d := DependManager{
		acID: acID,
	}
	d.deployedCh = make(map[string][]appData)
	d.readyCh = make(map[string][]appData)
	d.appCh = make(map[string][]chan struct{})
	d.resCh = make(map[string][]resData)
	dmList[acID] = &d
	return &d
}

// Functions registers an app for dependency
func (dm *DependManager) AddDependency(app string, dep map[string]*types.Criteria) error {

	dm.Lock()
	defer dm.Unlock()

	log.Info("AddDependency", log.Fields{"app": app, "dep": dep})
	for d, c := range dep {
		depLabel := d
		ch := make(chan struct{}, 1)
		data := appData{app: app, crt: *c, ch: ch}
		if c.OpStatus == types.OpStatusReady {
			dm.readyCh[depLabel] = append(dm.readyCh[depLabel], data)
		} else if c.OpStatus == types.OpStatusDeployed {
			dm.deployedCh[depLabel] = append(dm.deployedCh[depLabel], data)
		} else {
			// Ignore it
			continue
		}
		// Add all channels to per app list
		dm.appCh[app] = append(dm.appCh[app], ch)
	}
	return nil
}

// Wait for all dependecy to be met
func (dm *DependManager) WaitForDependency(ctx context.Context, app string) error {

	var cases []reflect.SelectCase
	if len(dm.appCh[app]) <= 0 {
		return nil
	}
	log.Info("WaitForDependency", log.Fields{"app": app})
	// Add the case for ctx done
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	for _, ch := range dm.appCh[app] {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	num := len(cases) - 1
	// Wait for all channels to be available
	for i := 0; i < num; i++ {
		// Wait for all channels
		index, _, _ := reflect.Select(cases)
		switch index {
		case 0:
			// case <- ctx.Done()
			return nil
		default:
			// Delete the channel from list
			log.Info("WaitForDependency: Coming out of wait", log.Fields{"app": app, "index": index})
			// Some channel is done, remove it from the list
			cases = append(cases[:index], cases[index+1:]...)
			continue
		}
	}
	return nil
}

func (dm *DependManager) NotifyAppliedStatus(app string) {
	log.Info("NotifyAppliedStatus", log.Fields{"app": app})
	for _, d := range dm.deployedCh[app] {
		if d.crt.Wait != 0 {
			time.Sleep(time.Duration(d.crt.Wait) * time.Second)
		}
		d.ch <- struct{}{}
	}
}

func (dm *DependManager) NotifyReadyStatus(app string) {
	log.Info("NotifyReadyStatus", log.Fields{"app": app})
	for _, d := range dm.readyCh[app] {
		if d.crt.Wait != 0 {
			time.Sleep(time.Duration(d.crt.Wait) * time.Second)
		}
		d.ch <- struct{}{}
	}
}

// Update status for the App ready on a cluster and check if app ready on all clusters
func ResourcesReady(acID, app, cluster string) {

	// Check if AppContext has dependency
	dm, ok := dmList[acID]
	if !ok {
		return
	}
	str := app + SEPARATOR + cluster
	// If no app is waiting for ready status of the app
	// Not further processing needed
	if len(dm.readyCh[app]) == 0 && len(dm.resCh[str]) == 0 {
		return
	}
	// Notify the apps waiting for the app to be ready
	dm.NotifyReadyStatus(app)
}
