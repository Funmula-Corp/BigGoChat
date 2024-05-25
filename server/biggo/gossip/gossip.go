package gossip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGossipService(ci einterfaces.ClusterInterface, ps *platform.PlatformService) *GossipService {
	return &GossipService{
		ci: ci, ps: ps, cbMap: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
	}
}

type GossipService struct {
	UnimplementedClusterServer

	ci  einterfaces.ClusterInterface
	cds *platform.ClusterDiscoveryService
	ps  *platform.PlatformService

	running atomic.Bool

	cbMap map[model.ClusterEvent]einterfaces.ClusterMessageHandler

	grpcSoc net.Listener
	grpcSvr *grpc.Server

	mx sync.Mutex
}

func (g2s *GossipService) GetClusterDiscoveryService() *platform.ClusterDiscoveryService {
	return g2s.cds
}

func (g2s *GossipService) NewClusterDiscoveryService() {
	if g2s.cds != nil {
		g2s.cds.Stop()
	}

	g2s.cds = g2s.ps.NewClusterDiscoveryService()
	if *g2s.ps.Config().ClusterSettings.OverrideHostname != "" {
		g2s.cds.Hostname = *g2s.ps.Config().ClusterSettings.OverrideHostname
	} else if *g2s.ps.Config().ClusterSettings.UseIPAddress {
		g2s.cds.AutoFillIPAddress(
			*g2s.ps.Config().ClusterSettings.NetworkInterface,
			*g2s.ps.Config().ClusterSettings.AdvertiseAddress,
		)
	} else {
		g2s.cds.AutoFillHostname()
	}

	g2s.cds.ClusterName = *g2s.ps.Config().ClusterSettings.ClusterName
	g2s.cds.GossipPort = (int32)(*g2s.ps.Config().ClusterSettings.GossipPort)
	g2s.cds.Type = model.GetServiceEnvironment()
}

func (g2s *GossipService) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	g2s.cbMap[event] = crm
}

func (g2s *GossipService) StartInterNodeCommunication() (err error) {
	g2s.mx.Lock()
	defer g2s.mx.Unlock()

	if !g2s.running.Load() {
		g2s.NewClusterDiscoveryService()
		if g2s.grpcSoc, err = net.Listen("tcp", fmt.Sprintf(":%d", g2s.cds.GossipPort)); err != nil {
			return
		}

		g2s.grpcSvr = grpc.NewServer()
		RegisterClusterServer(g2s.grpcSvr, g2s)

		go func() {
			defer g2s.running.Store(false)
			defer g2s.grpcSoc.Close()
			defer g2s.cds.Stop()
			g2s.running.Store(true)
			g2s.cds.Start()
			g2s.grpcSvr.Serve(g2s.grpcSoc)
		}()
	}
	return
}

func (g2s *GossipService) StopInterNodeCommunication() {
	g2s.mx.Lock()
	defer g2s.mx.Unlock()

	if g2s.running.Load() {
		g2s.grpcSvr.Stop()
	}
}

func (g2s *GossipService) GetMyClusterInfo(ctx context.Context, in *Void) (reply *ClusterInfoReply, err error) {
	cInfo := g2s.ci.GetMyClusterInfo()
	reply = &ClusterInfoReply{
		Id:            cInfo.Id,
		Version:       cInfo.Version,
		SchemaVersion: cInfo.SchemaVersion,
		ConfigHash:    cInfo.ConfigHash,
		IPAddress:     cInfo.IPAddress,
		Hostname:      cInfo.Hostname,
	}
	return
}

func (g2s *GossipService) CallGetMyClusterInfo(addr string) (result *model.ClusterInfo, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *ClusterInfoReply
	client := NewClusterClient(soc)
	if reply, err = client.GetMyClusterInfo(context.Background(), &Void{}); err != nil {
		return
	}

	result = &model.ClusterInfo{
		Id:            reply.Id,
		Version:       reply.Version,
		SchemaVersion: reply.SchemaVersion,
		ConfigHash:    reply.ConfigHash,
		IPAddress:     reply.IPAddress,
		Hostname:      reply.Hostname,
	}
	return
}

func (g2s *GossipService) GetPluginStatuses(ctx context.Context, in *Void) (reply *PluginStatusReply, err error) {
	var (
		appErr *model.AppError
		pStats model.PluginStatuses
	)

	if pStats, appErr = g2s.ps.GetPluginStatuses(); appErr != nil {
		err = appErr.Unwrap()
		return
	}

	reply = &PluginStatusReply{Statuses: []*PluginStatus{}}
	for _, pStat := range pStats {
		reply.Statuses = append(reply.Statuses, &PluginStatus{
			PluginId:    pStat.PluginId,
			ClusterId:   pStat.ClusterId,
			PluginPath:  pStat.PluginPath,
			State:       int32(pStat.State),
			Error:       pStat.Error,
			Name:        pStat.Name,
			Description: pStat.Description,
			Version:     pStat.Version,
		})
		break
	}
	return
}

func (g2s *GossipService) CallGetPluginStatuses(addr string) (result model.PluginStatuses, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *PluginStatusReply
	client := NewClusterClient(soc)
	if reply, err = client.GetPluginStatuses(context.Background(), &Void{}); err != nil {
		return
	}

	result = model.PluginStatuses{}
	for _, pStat := range reply.Statuses {
		result = append(result, &model.PluginStatus{
			PluginId:    pStat.PluginId,
			ClusterId:   pStat.ClusterId,
			PluginPath:  pStat.PluginPath,
			State:       int(pStat.State),
			Error:       pStat.Error,
			Name:        pStat.Name,
			Description: pStat.Description,
			Version:     pStat.Version,
		})
	}
	return
}

func (g2s *GossipService) SendClusterMessageToNode(ctx context.Context, in *ClusterMessage) (reply *Void, err error) {
	msg := model.ClusterMessage{
		Event:            model.ClusterEvent(in.Event),
		SendType:         in.SendType,
		WaitForAllToSend: in.WaitForAllToSend,
		Data:             in.Data,
		Props:            in.Props,
	}

	if cb, found := g2s.cbMap[msg.Event]; found {
		go cb(&msg)
	}

	mlog.Debug("GRPC SendClusterMessageToNode")
	return
}

func (g2s *GossipService) CallSendClusterMessageToNode(addr string, msg *model.ClusterMessage) (err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	client := NewClusterClient(soc)
	_, err = client.SendClusterMessageToNode(context.Background(), &ClusterMessage{
		Event:            string(msg.Event),
		SendType:         msg.SendType,
		WaitForAllToSend: msg.WaitForAllToSend,
		Data:             msg.Data,
		Props:            msg.Props,
	})
	return
}

func (g2s *GossipService) GetClusterStats(ctx context.Context, in *Void) (reply *ClusterStatsReply, err error) {
	reply = &ClusterStatsReply{
		Id:                        g2s.cds.Id,
		TotalWebsocketConnections: int32(g2s.ps.TotalWebsocketConnections()),
		TotalReadDbConnections:    int32(g2s.ps.Store.TotalReadDbConnections()),
		TotalMasterDbConnections:  int32(g2s.ps.Store.TotalMasterDbConnections()),
	}
	return
}

func (g2s *GossipService) CallGetClusterStats(addr string) (result *model.ClusterStats, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *ClusterStatsReply
	client := NewClusterClient(soc)
	if reply, err = client.GetClusterStats(context.Background(), &Void{}); err != nil {
		return
	}

	result = &model.ClusterStats{
		Id:                        reply.Id,
		TotalWebsocketConnections: int(reply.TotalWebsocketConnections),
		TotalReadDbConnections:    int(reply.TotalReadDbConnections),
		TotalMasterDbConnections:  int(reply.TotalMasterDbConnections),
	}
	return
}

func (g2s *GossipService) GetLogs(ctx context.Context, in *LogRequest) (reply *LogReply, err error) {
	// TODO: figure out how to get logs
	reply = &LogReply{Entries: []string{}}
	return
}

func (g2s *GossipService) CallGetLogs(addr string, page, perPage int) (result *[]string, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *LogReply
	client := NewClusterClient(soc)
	if reply, err = client.GetLogs(context.Background(), &LogRequest{Page: int32(page), PerPage: int32(perPage)}); err != nil {
		return
	}

	result = &reply.Entries
	return
}

func (g2s *GossipService) QueryLogs(ctx context.Context, in *LogRequest) (reply *QueryLogReply, err error) {
	// TODO: figure out how to get logs
	reply = &QueryLogReply{Id: g2s.cds.Id, Entries: []string{}}
	return
}

func (g2s *GossipService) CallQueryLogs(addr string, page, perPage int) (result *map[string][]string, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *QueryLogReply
	client := NewClusterClient(soc)
	if reply, err = client.QueryLogs(context.Background(), &LogRequest{Page: int32(page), PerPage: int32(perPage)}); err != nil {
		return
	}

	result = &map[string][]string{
		reply.Id: reply.Entries,
	}
	return
}

func (g2s *GossipService) ConfigChanged(ctx context.Context, in *ConfigUpdateRequest) (reply *Void, err error) {
	var newConfig model.Config
	if err = json.NewDecoder(bytes.NewBuffer(in.NewConfigBuffer)).Decode(&newConfig); err != nil {
		return
	}
	_, _, err = g2s.ps.GetConfigStore().Set(&newConfig)
	return
}

func (g2s *GossipService) CallConfigChanged(addr string, previousConfig, newConfig *model.Config) (err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	ocb := bytes.NewBuffer([]byte{})
	json.NewEncoder(ocb).Encode(previousConfig)
	ncb := bytes.NewBuffer([]byte{})
	json.NewEncoder(ncb).Encode(newConfig)

	client := NewClusterClient(soc)
	_, err = client.ConfigChanged(context.Background(), &ConfigUpdateRequest{
		OldConfigBuffer: ocb.Bytes(),
		NewConfigBuffer: ncb.Bytes(),
	})
	return
}

func (g2s *GossipService) WebConnCountForUser(ctx context.Context, in *WebsocketCountForUserRequest) (reply *WebsocketCountForUserReply, err error) {
	reply = &WebsocketCountForUserReply{
		Count: int32(g2s.ps.WebConnCountForUser(in.UserId)),
	}
	return
}

func (g2s *GossipService) CallWebConnCountForUser(addr string, userID string) (result int, err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	var reply *WebsocketCountForUserReply
	client := NewClusterClient(soc)
	if reply, err = client.WebConnCountForUser(context.Background(), &WebsocketCountForUserRequest{UserId: userID}); err != nil {
		return
	}
	result = int(reply.Count)
	return
}
