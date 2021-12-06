// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020-2021 Intel Corporation

package statusnotifyserver

import (
	"context"
	"sync"

	pkgerrors "github.com/pkg/errors"
	inc "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/installappclient"
	pb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotify"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/rpc"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module/controller"
	readynotifypb "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/readynotify"
)

// StatusNotifyServerHelpers is an interface supported by the specific microservices that need to provide status notification
type StatusNotifyServerHelpers interface {
	GetAppContextId(reg *pb.StatusRegistration) (string, error)
	PrepareStatusNotification(reg *pb.StatusRegistration) *pb.StatusNotification
}

// streamInfo contains information about a given status notification stream, including:
// the stream server, and information about the type of notifications desired and
// filter and output details.
type streamInfo struct {
	stream       pb.StatusNotify_StatusRegisterServer
	reg          *pb.StatusRegistration
	appContextID string
}

type appContextInfo struct {
	readyNotifyStream readynotifypb.ReadyNotify_AlertClient
	statusClientIDs   map[string]struct{}
}

// StatusNotifyServer will be initialized by NewStatusNotifyServer() and
// its lifecycle is valid until all the clients unsubscribed the stream notification channel
type StatusNotifyServer struct {
	name string
	// clientId is expected to be unique.  Registering with a clientId that is already in the
	// map will be rejected - i.e. the new stream will close immediately with an error
	appContexts       map[string]appContextInfo // map[appcontextid]appContextInfo
	statusClients     map[string]streamInfo     // map[clientId]streamInfo
	streamChannels    map[pb.StatusNotify_StatusRegisterServer]chan int
	sh                StatusNotifyServerHelpers
	readyNotifyClient readynotifypb.ReadyNotifyClient
	mutex             sync.Mutex
}

var notifServer *StatusNotifyServer

// StatusRegister gets notified when a client registers for a status notification stream for a given resource
func (s *StatusNotifyServer) StatusRegister(reg *pb.StatusRegistration, stream pb.StatusNotify_StatusRegisterServer) error {

	// Check if the clientId is already in use, return error if yes
	clientId := reg.GetClientId()
	if len(clientId) == 0 {
		logutils.Info("[StatusNotify gRPC] Recieved a status notification registration with invalid client ID", logutils.Fields{})
		return pkgerrors.New("Invalid client ID")
	}
	if _, ok := s.statusClients[clientId]; ok {
		logutils.Info("[StatusNotify gRPC] Recieved a duplicate status notification registration", logutils.Fields{"client": clientId})
		return pkgerrors.New("Duplicate client ID: " + clientId)
	}
	appContextID, err := s.sh.GetAppContextId(reg)
	if err != nil {
		logutils.Info("[StatusNotify gRPC] Could not get appContextID for status notification registration", logutils.Fields{"client": clientId, "AppContextID": appContextID})
		return err
	}

	logutils.Info("[StatusNotify gRPC] Recieved a status notification registration", logutils.Fields{"client": clientId, "appContextID": appContextID})

	// Add the client info to the statusnotify server maps
	needReadyNotifyStream := false
	s.mutex.Lock()

	// update appContexts
	if _, ok := s.appContexts[appContextID]; !ok {
		s.appContexts[appContextID] = appContextInfo{
			readyNotifyStream: nil,
			statusClientIDs:   make(map[string]struct{}),
		}
		needReadyNotifyStream = true
	}
	s.appContexts[appContextID].statusClientIDs[clientId] = struct{}{}

	// update statusClients
	s.statusClients[clientId] = streamInfo{
		stream:       stream,
		reg:          reg,
		appContextID: appContextID,
	}

	// update streamChannels
	s.streamChannels[stream] = make(chan int)
	c := s.streamChannels[stream]
	ctx := stream.Context()

	if needReadyNotifyStream {
		if s.readyNotifyClient == nil {
			s.readyNotifyClient = newReadyNotifyClient()
			if s.readyNotifyClient == nil {
				s.mutex.Unlock()
				logutils.Error("[StatusNotify gRPC] Could not get ReadyNotify Client", logutils.Fields{"appContextID": appContextID, "client": clientId})
				return pkgerrors.Errorf("Unable to get ReadyNotifyClient for StatusNotifyServer: %v, %v", appContextID, clientId)
			}
		}
		readyNotifyStream, err := s.readyNotifyClient.Alert(context.Background(), &readynotifypb.Topic{ClientName: s.name, AppContext: appContextID})
		if err != nil {
			s.mutex.Unlock()
			logutils.Error("[StatusNotify gRPC] Could not get ReadyNotify Stream", logutils.Fields{"appContextID": appContextID, "client": clientId})
			return err
		}

		// set the readyNotifyStream for this appContextID - check for a race - no need to set if it's already been set
		acInfo := s.appContexts[appContextID]
		acInfo.readyNotifyStream = readyNotifyStream
		s.mutex.Unlock()

		logutils.Info("[StatusNotify gRPC] ready to start sending status notifications", logutils.Fields{"client": clientId})
		go sendStatusNotifications(readyNotifyStream)
	} else {
		s.mutex.Unlock()
	}

	// Keep stream open
	for {
		select {
		case <-ctx.Done():
			logutils.Info("[StatusNotify gRPC] Client has disconnected", logutils.Fields{"client": clientId})
			cleanup(clientId)
			// need to clean up here ?
			return nil
		case <-c:
			logutils.Info("[StatusNotify gRPC] Stop channel has been triggered for client", logutils.Fields{"client": clientId})
			return nil
		default:
		}
	}
}

// SendStatusNotification sends a status notification message to the subscriber
func sendStatusNotifications(stream readynotifypb.ReadyNotify_AlertClient) error {
	var appContextID string

	// TESTING loop - send notifications to everyone every 10 seconds
	notifServer.mutex.Lock()
	for appContextID, _ := range notifServer.appContexts {
		acInfo, _ := notifServer.appContexts[appContextID]
		for clientId, _ := range acInfo.statusClientIDs {
			si := notifServer.statusClients[clientId]
			err := si.stream.Send(notifServer.sh.PrepareStatusNotification(si.reg))
			if err != nil {
				logutils.Error("[StatusNotify gRPC] TEST Status notification failed to be sent", logutils.Fields{"clientId": clientId, "err": err})
				// return pkgerrors.New("Notification failed")
			} else {
				logutils.Info("[StatusNotify gRPC] TEST Status notification was sent", logutils.Fields{"clientId": clientId, "appContextID": appContextID})
			}
		}
	}
	notifServer.mutex.Unlock()

	for true {

		resp, err := stream.Recv()
		// TODO: check for io.EOF
		// TODO: some kind of throttle mechanism needed? (in the case of receiving many events at once)
		if err != nil {
			logutils.Error("[ReadyNotify gRPC] Failed to receive notification", logutils.Fields{"err": err})
			return err
		}

		notifServer.mutex.Lock()
		acInfo, ok := notifServer.appContexts[resp.AppContext]
		if !ok {
			notifServer.mutex.Unlock()
			logutils.Warn("[StatusNotify gRPC] Received a ReadyNotify alert from rsync for missing appContext", logutils.Fields{"appContextID": appContextID})
			continue
		}

		// For a given alert, send a status notification to each status client watching the appcontextId
		for clientId, _ := range acInfo.statusClientIDs {
			si := notifServer.statusClients[clientId]
			err := si.stream.Send(notifServer.sh.PrepareStatusNotification(si.reg))
			if err != nil {
				logutils.Error("[StatusNotify gRPC] Status notification failed to be sent", logutils.Fields{"clientId": clientId, "err": err})
			} else {
				logutils.Info("[StatusNotify gRPC] Status notification was sent", logutils.Fields{"clientId": clientId, "appContextID": appContextID})
			}
		}
		notifServer.mutex.Unlock()

	}
	return nil
}

// cleanup will be called when the subscriber wants to terminate the stream
func cleanup(clientId string) {
	notifServer.mutex.Lock()

	// get the clientId entry, stop the stream
	si, ok := notifServer.statusClients[clientId]
	if !ok {
		// client already deregistered, or never existed
		notifServer.mutex.Unlock()
		return
	}
	/*
		notifServer.streamChannels[si.stream] <- 1
		logutils.Info("[StatusNotify gRPC]   cleanup sent channel message", logutils.Fields{})
	*/
	delete(notifServer.statusClients, clientId)
	delete(notifServer.streamChannels, si.stream)

	// remove the clientId from the appContext Info
	acInfo, ok := notifServer.appContexts[si.appContextID]
	if !ok {
		// this should not occur, but appcontext is clear already
		notifServer.mutex.Unlock()
		return
	}
	delete(acInfo.statusClientIDs, clientId)
	if len(acInfo.statusClientIDs) == 0 {
		// if no clients are left for the app context
		//acInfo.readyNotifyStream.CloseSend()  ? - this crashed orchestrator
		delete(notifServer.appContexts, si.appContextID)
	}

	logutils.Info("[StatusNotify gRPC] Cleaned up", logutils.Fields{"clientId": clientId, "appContextID": si.appContextID})
	notifServer.mutex.Unlock()
	return
}

// StatusDeregister will be called when the subscriber wants to terminate the stream
func (s *StatusNotifyServer) StatusDeregister(ctx context.Context, dereg *pb.StatusDeregistration) (*pb.StatusDeregistrationResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// get the clientId entry, stop the stream
	si, ok := s.statusClients[dereg.ClientId]
	if !ok {
		// client already deregistered, or never existed
		return &pb.StatusDeregistrationResponse{}, nil
	}
	s.streamChannels[si.stream] <- 1
	delete(s.statusClients, dereg.ClientId)
	delete(s.streamChannels, si.stream)

	// remove the clientId from the appContext Info
	acInfo, ok := s.appContexts[si.appContextID]
	if !ok {
		// this should not occur, but appcontext is clear already
		return &pb.StatusDeregistrationResponse{}, nil
	}
	delete(acInfo.statusClientIDs, dereg.ClientId)
	if len(acInfo.statusClientIDs) == 0 {
		// if no clients are left for the app context
		acInfo.readyNotifyStream.CloseSend()
		delete(s.appContexts, si.appContextID)
	}

	return &pb.StatusDeregistrationResponse{}, nil
}

// NewStatusNotifyServer will create a new StatusNotifyServer and destroys the previous one
func NewStatusNotifyServer(readyNotifyClientID string, sh StatusNotifyServerHelpers) *StatusNotifyServer {

	s := &StatusNotifyServer{
		name:              readyNotifyClientID,
		appContexts:       make(map[string]appContextInfo),
		statusClients:     make(map[string]streamInfo),
		streamChannels:    make(map[pb.StatusNotify_StatusRegisterServer]chan int),
		sh:                sh,
		readyNotifyClient: nil, // initialize this later on Registration call
	}
	notifServer = s
	return s
}

const rsyncName = "rsync"

func queryDBAndInitRsync() error {
	client := controller.NewControllerClient("resources", "data", "orchestrator")
	vals, _ := client.GetControllers()
	for _, v := range vals {
		if v.Metadata.Name == rsyncName {
			logutils.Info("Initializing RPC connection to resource synchronizer", logutils.Fields{
				"Controller": v.Metadata.Name,
			})
			inc.NewRsyncInfo(v.Metadata.Name, v.Spec.Host, v.Spec.Port)
			inc.InitRsyncClient()
			return nil
		}
	}
	return pkgerrors.Errorf("queryRsyncInfoInMCODB Failed - Could not get find rsync by name : %v", rsyncName)
}

func newReadyNotifyClient() readynotifypb.ReadyNotifyClient {
	conn := rpc.GetRpcConn(rsyncName)
	if conn == nil {
		if !inc.InitRsyncClient() {
			err := queryDBAndInitRsync()
			if err != nil {
				logutils.Warn("[ReadyNotify gRPC] Failed to initialize get ReadyNotifyClient", logutils.Fields{})
				return nil
			}
		}
		conn = rpc.GetRpcConn(rsyncName)
	}

	if conn != nil {
		return readynotifypb.NewReadyNotifyClient(conn)
	} else {
		logutils.Warn("[ReadyNotify gRPC] Failed to get ReadyNotifyClient", logutils.Fields{})
		return nil
	}
}

// GetStatusParameters retrieves the status query parameters from the StatusRegistration
func GetStatusParameters(reg *pb.StatusRegistration) (string, string, []string, []string, []string) {
	var output, statusType string

	switch reg.StatusType {
	case pb.StatusValue_READY:
		statusType = "cluster"
	case pb.StatusValue_DEPLOYED:
		statusType = "rsync"
	}

	switch reg.Output {
	case pb.OutputType_SUMMARY:
		output = "cluster"
	case pb.OutputType_ALL:
		output = "all"
	}

	return statusType, output, reg.Apps, reg.Clusters, reg.Resources
}
