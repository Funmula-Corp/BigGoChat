package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/proto"
)

func (p *GossipServer) GetClusterInfos(ctx context.Context, request *proto.GetClusterInfosRequest) (*proto.GetClusterInfosReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	info := p.platformService.Cluster().GetMyClusterInfo()
	buffer := bytes.NewBuffer([]byte{})
	if err := gob.NewEncoder(buffer).Encode(info); err != nil {
		return nil, err
	}

	return &proto.GetClusterInfosReply{
		Buffer: buffer.Bytes(),
	}, nil
}

func (p *GossipServer) SendClusterMessageToNode(ctx context.Context, request *proto.SendClusterMessageToNodeRequest) (*proto.SendClusterMessageToNodeReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	p.ClusterMessageHandler[model.ClusterEvent(request.GetEvent())](&model.ClusterMessage{
		Event: model.ClusterEvent(request.GetEvent()),
		Data:  request.GetData(),
		Props: request.GetProps(),
	})

	return &proto.SendClusterMessageToNodeReply{}, nil
}

func (p *GossipServer) GetPluginStatuses(ctx context.Context, request *proto.GetPluginStatusesRequest) (*proto.GetPluginStatusesReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	pluginStatus := []*proto.PluginStatuses{}

	pluginStatuses, _ := p.platformService.GetPluginStatuses()
	for _, status := range pluginStatuses {
		pluginStatus = append(pluginStatus, &proto.PluginStatuses{
			PluginId:    status.PluginId,
			ClusterId:   status.ClusterId,
			PluginPath:  status.PluginPath,
			State:       int64(status.State),
			Error:       status.Error,
			Name:        status.Name,
			Description: status.Description,
			Version:     status.Version,
		})
	}

	return &proto.GetPluginStatusesReply{
		PluginStatus: pluginStatus,
	}, nil
}

func (p *GossipServer) GetClusterStats(ctx context.Context, request *proto.GetClusterStatsRequest) (*proto.GetClusterStatsReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	return &proto.GetClusterStatsReply{
		TotalWebsocketConnections: int64(p.platformService.TotalWebsocketConnections()),
		TotalReadDbConnections:    int64(p.platformService.Store.TotalReadDbConnections()),
		TotalMasterDbConnections:  int64(p.platformService.Store.TotalMasterDbConnections()),
	}, nil
}

func (p *GossipServer) ConfigChanged(ctx context.Context, request *proto.ConfigChangedRequest) (*proto.ConfigChangedReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := gob.NewEncoder(buffer).Encode(p.platformService.Config()); err != nil {
		return nil, fmt.Errorf("cluster.serialize.config.current.error: %v", err)
	}

	hash := md5.New()
	hash.Write(buffer.Bytes())
	if hex.EncodeToString(hash.Sum(nil)) != request.Hash {
		return &proto.ConfigChangedReply{}, nil
	}

	newConfig := &model.Config{}
	if err := gob.NewDecoder(bytes.NewBuffer(request.GetConfigBuffer())).Decode(newConfig); err != nil {
		return nil, fmt.Errorf("cluster.deserialize.config.new.error: %v", err)
	}
	p.platformService.SaveConfig(newConfig, false)
	return &proto.ConfigChangedReply{}, nil
}

func (p *GossipServer) WebConnCountForUser(ctx context.Context, request *proto.WebConnCountForUserRequest) (*proto.WebConnCountForUserReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}
	return &proto.WebConnCountForUserReply{
		Count: int64(p.platformService.WebConnCountForUser(request.UserID)),
	}, nil
}

func (p *GossipServer) GetLogs(ctx context.Context, request *proto.GetLogsRequest) (*proto.GetLogsReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	reply := &proto.GetLogsReply{}
	if logs, err := p.platformService.GetLogsSkipSend(int(request.Page), int(request.PerPage), &model.LogFilter{}); err == nil {
		reply.LogRecord = logs
	}
	return reply, nil
}

func (p *GossipServer) QueryLogs(ctx context.Context, request *proto.QueryLogsRequest) (*proto.QueryLogsReply, error) {
	if request == nil {
		return nil, ErrMissingRequest
	}

	filter := &model.LogFilter{}
	if request.LogFilter != nil {
		filter.DateFrom = request.LogFilter.DateFrom
		filter.DateTo = request.LogFilter.DateTo
		filter.LogLevels = request.LogFilter.LogLevels
		filter.ServerNames = request.LogFilter.ServerNames
	}

	reply := &proto.QueryLogsReply{}
	if logs, err := p.platformService.GetLogsSkipSend(int(request.Page), int(request.PerPage), filter); err == nil {
		reply.LogRecord = logs
	}
	return reply, nil
}
