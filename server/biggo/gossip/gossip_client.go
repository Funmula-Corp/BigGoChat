package gossip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (g2s *GossipService) CallGetMyClusterInfo(addr string) (result *model.ClusterInfo, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *ClusterInfoReply
	if reply, err = client.GetMyClusterInfo(context.Background(), &Void{}); err != nil {
		mlog.Error("CallGetMyClusterInfo", mlog.Err(err))
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

func (g2s *GossipService) CallSendClusterMessageToNode(addr string, msg *model.ClusterMessage) (err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	if _, err = client.SendClusterMessageToNode(context.Background(), &ClusterMessage{
		Event:            string(msg.Event),
		SendType:         msg.SendType,
		WaitForAllToSend: msg.WaitForAllToSend,
		Data:             msg.Data,
		Props:            msg.Props,
	}); err != nil {
		mlog.Error("CallSendClusterMessageToNode", mlog.Err(err))
		return
	}
	return
}

func (g2s *GossipService) CallGetClusterStats(addr string) (result *model.ClusterStats, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *ClusterStatsReply
	if reply, err = client.GetClusterStats(context.Background(), &Void{}); err != nil {
		mlog.Error("CallGetClusterStats", mlog.Err(err))
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

func (g2s *GossipService) CallGetLogs(addr string, page, perPage int) (result []string, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *LogReply
	if reply, err = client.GetLogs(context.Background(), &LogRequest{Page: int32(page), PerPage: int32(perPage)}); err != nil {
		mlog.Error("CallGetLogs", mlog.Err(err))
		return
	}

	result = reply.Entries
	return
}

func (g2s *GossipService) CallQueryLogs(addr string, page, perPage int) (result []string, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *QueryLogReply
	if reply, err = client.QueryLogs(context.Background(), &LogRequest{Page: int32(page), PerPage: int32(perPage)}); err != nil {
		mlog.Error("CallQueryLogs", mlog.Err(err))
		return
	}

	result = reply.Entries
	return
}

func (g2s *GossipService) CallGetPluginStatuses(addr string) (result model.PluginStatuses, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *PluginStatusReply
	if reply, err = client.GetPluginStatuses(context.Background(), &Void{}); err != nil {
		mlog.Error("CallGetPluginStatuses", mlog.Err(err))
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

func (g2s *GossipService) CallConfigChanged(addr string, previousConfig, newConfig *model.Config) (err error) {
	var soc *grpc.ClientConn
	if soc, err = grpc.NewClient(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	defer soc.Close()

	mlog.Debug("CallConfigChanged", mlog.Any("newConfig", newConfig), mlog.Err(err))
	ocb := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(ocb).Encode(previousConfig); err != nil {
		return
	}
	ncb := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(ncb).Encode(newConfig); err != nil {
		return
	}

	client := NewClusterClient(soc)
	_, err = client.ConfigChanged(context.Background(), &ConfigUpdateRequest{
		OldConfigBuffer: ocb.Bytes(),
		NewConfigBuffer: ncb.Bytes(),
	})
	return
}

func (g2s *GossipService) CallWebConnCountForUser(addr string, userID string) (result int, err error) {
	var (
		client     ClusterClient
		connection *grpc.ClientConn
	)

	if client, connection, err = g2s.NewClient(addr); err != nil {
		return
	}
	defer connection.Close()

	var reply *WebsocketCountForUserReply
	if reply, err = client.WebConnCountForUser(context.Background(), &WebsocketCountForUserRequest{UserId: userID}); err != nil {
		mlog.Error("CallWebConnCountForUser", mlog.Err(err))
		return
	}
	result = int(reply.Count)
	return
}
