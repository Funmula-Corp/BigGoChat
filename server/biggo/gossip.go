package biggo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

type GossipRequestEventType string
type GossipResponseEventType string

const (
	ClusterGossipEventRequestGetLogs            GossipRequestEventType  = model.ClusterGossipEventRequestGetLogs
	ClusterGossipEventResponseGetLogs           GossipResponseEventType = model.ClusterGossipEventResponseGetLogs
	ClusterGossipEventRequestGetClusterStats    GossipRequestEventType  = model.ClusterGossipEventRequestGetClusterStats
	ClusterGossipEventResponseGetClusterStats   GossipResponseEventType = model.ClusterGossipEventResponseGetClusterStats
	ClusterGossipEventRequestGetPluginStatuses  GossipRequestEventType  = model.ClusterGossipEventRequestGetPluginStatuses
	ClusterGossipEventResponseGetPluginStatuses GossipResponseEventType = model.ClusterGossipEventResponseGetPluginStatuses
	ClusterGossipEventRequestSaveConfig         GossipRequestEventType  = model.ClusterGossipEventRequestSaveConfig
	ClusterGossipEventResponseSaveConfig        GossipResponseEventType = model.ClusterGossipEventResponseSaveConfig
	ClusterGossipEventRequestWebConnCount       GossipRequestEventType  = model.ClusterGossipEventRequestWebConnCount
	ClusterGossipEventResponseWebConnCount      GossipResponseEventType = model.ClusterGossipEventResponseWebConnCount

	ClusterGossipEventRequestInfo     GossipRequestEventType  = "gossip_request_cluster_info"
	ClusterGossipEventResponseInfo    GossipResponseEventType = "gossip_response_cluster_info"
	ClusterGossipEventRequestMessage  GossipRequestEventType  = "gossip_request_cluster_message"
	ClusterGossipEventResponseMessage GossipResponseEventType = "gossip_response_cluster_message"
)

type GossipService struct {
	c *BiggoCluster

	gCtx    context.Context
	gCancel context.CancelFunc
}

type GossipServiceRequest struct {
	Type   GossipRequestEventType
	Buffer []byte
}

type GossipServiceResponse struct {
	Type   GossipResponseEventType
	Buffer []byte
}

func (gs *GossipService) RecvGossip(request *GossipServiceRequest, response *GossipServiceResponse) (err error) {
	switch request.Type {
	case ClusterGossipEventRequestMessage:
		response.Type = ClusterGossipEventResponseMessage
		msg := new(model.ClusterMessage)
		if err = json.Unmarshal(request.Buffer, msg); err == nil {
			if cb, ok := gs.c.cbMap[msg.Event]; ok {
				go cb(msg)
			} else {
				mlog.Warn("No callback registered for ClusterMessage Event:", logr.Any("ClusterMessageEvent", msg.Event))
			}
		}
	case ClusterGossipEventRequestSaveConfig:
		response.Type = ClusterGossipEventResponseSaveConfig
		cfg := new(model.Config)
		if err = json.Unmarshal(request.Buffer, cfg); err == nil {
			gs.c.ConfigChanged(nil, cfg, false)
		}
	case ClusterGossipEventRequestInfo:
		response.Type = ClusterGossipEventResponseInfo
		response.Buffer, _ = json.Marshal(gs.c.GetMyClusterInfo())
	default:
		err = errors.New("unknown gossip event type")
	}
	return
}

func (gs *GossipService) SendGossip(host string, eventType GossipRequestEventType, value interface{}) (result interface{}, err error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, *gs.c.ps.Config().ClusterSettings.GossipPort))
	if err != nil {
		mlog.Error("Gossip IPC Connection Error", logr.Err(err))
		return
	}
	defer client.Close()

	var buffer []byte
	if value != nil {
		buffer, _ = json.Marshal(value)
	}
	request := &GossipServiceRequest{
		Type:   eventType,
		Buffer: buffer,
	}

	response := &GossipServiceResponse{}
	if err = client.Call("GossipService.RecvGossip", request, response); err != nil {
		mlog.Error("Gossip IPC Call Error", logr.Err(err))
	} else {
		if len(response.Buffer) > 0 {
			if verifyResponseType(request.Type, response.Type) {
				err = errors.New("request-response event type mismatch")
				mlog.Error("Gossip IPC Response Error", logr.Err(err))
				return
			}

			switch response.Type {
			case ClusterGossipEventResponseInfo:
				var info model.ClusterInfo
				err = json.Unmarshal(response.Buffer, &info)
				result = &info
			default:
				err = errors.New("unknown response event type")
				mlog.Error("Gossip IPC Response Error", logr.Err(err))
			}
		}
	}
	return
}

func verifyResponseType(reqType GossipRequestEventType, resType GossipResponseEventType) bool {
	return (reqType == ClusterGossipEventRequestGetLogs && resType == ClusterGossipEventResponseGetLogs) ||
		(reqType == ClusterGossipEventRequestGetClusterStats && resType == ClusterGossipEventResponseGetClusterStats) ||
		(reqType == ClusterGossipEventRequestGetPluginStatuses && resType == ClusterGossipEventResponseGetPluginStatuses) ||
		(reqType == ClusterGossipEventRequestSaveConfig && resType == ClusterGossipEventResponseSaveConfig) ||
		(reqType == ClusterGossipEventRequestWebConnCount && resType == ClusterGossipEventResponseWebConnCount) ||
		(reqType == ClusterGossipEventRequestMessage && resType == ClusterGossipEventResponseMessage)
}

func (gs *GossipService) start() {
	if gs.gCtx == nil && gs.gCancel == nil {
		gs.gCtx, gs.gCancel = context.WithCancel(context.Background())

		go func(ctx context.Context) {
			defer gs.stop()

			if *gs.c.ps.Config().ClusterSettings.OverrideHostname != "" {
				gs.c.cds.ClusterDiscovery.Hostname = *gs.c.ps.Config().ClusterSettings.OverrideHostname
			} else {
				gs.c.cds.ClusterDiscovery.Hostname = ""
			}

			if *gs.c.ps.Config().ClusterSettings.UseIPAddress {
				gs.c.cds.ClusterDiscovery.AutoFillIPAddress(
					*gs.c.ps.Config().ClusterSettings.NetworkInterface,
					*gs.c.ps.Config().ClusterSettings.AdvertiseAddress,
				)
			} else {
				gs.c.cds.ClusterDiscovery.AutoFillHostname()
			}

			gs.c.cds.ClusterDiscovery.ClusterName = *gs.c.ps.Config().ClusterSettings.ClusterName
			gs.c.cds.ClusterDiscovery.GossipPort = (int32)(*gs.c.ps.Config().ClusterSettings.GossipPort)
			gs.c.cds.ClusterDiscovery.Type = model.GetServiceEnvironment()

			if srv, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.c.cds.ClusterDiscovery.GossipPort)); err != nil {
				mlog.Error("Gossip IPC Socket Error", logr.Err(err))
				return
			} else {
				defer srv.Close()
				defer gs.stop()
				defer gs.c.cds.Stop()
				gs.c.cds.Start()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						conn, err := srv.Accept()
						if err != nil {
							mlog.Error("Gossip IPC Accept Error", logr.Err(err))
							continue
						}
						go rpc.ServeConn(conn)
					}
				}
			}
		}(gs.gCtx)
	}
}

func (gs *GossipService) stop() {
	if gs.gCancel != nil {
		gs.gCancel()
	}
}
