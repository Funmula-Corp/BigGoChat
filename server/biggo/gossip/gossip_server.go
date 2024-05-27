package gossip

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (g2s *GossipService) GetMyClusterInfo(ctx context.Context, in *Void) (reply *ClusterInfoReply, err error) {
	mlog.Debug("GetMyClusterInfo-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	cInfo := g2s.ci.GetMyClusterInfo()
	reply = &ClusterInfoReply{
		Id:            cInfo.Id,
		Version:       cInfo.Version,
		SchemaVersion: cInfo.SchemaVersion,
		ConfigHash:    cInfo.ConfigHash,
		IPAddress:     cInfo.IPAddress,
		Hostname:      cInfo.Hostname,
	}
	mlog.Debug("GetMyClusterInfo-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) SendClusterMessageToNode(ctx context.Context, in *ClusterMessage) (reply *Void, err error) {
	mlog.Debug("SendClusterMessageToNode-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
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
	mlog.Debug("SendClusterMessageToNode-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Err(err))
	return
}

func (g2s *GossipService) GetClusterStats(ctx context.Context, in *Void) (reply *ClusterStatsReply, err error) {
	mlog.Debug("GetClusterStats-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	reply = &ClusterStatsReply{
		Id:                        g2s.cds.Id,
		TotalWebsocketConnections: int32(g2s.ps.TotalWebsocketConnections()),
		TotalReadDbConnections:    int32(g2s.ps.Store.TotalReadDbConnections()),
		TotalMasterDbConnections:  int32(g2s.ps.Store.TotalMasterDbConnections()),
	}
	mlog.Debug("GetClusterStats-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) GetLogs(ctx context.Context, in *LogRequest) (reply *LogReply, err error) {
	mlog.Debug("GetLogs-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	logs, appErr := g2s.ps.GetLogsSkipSend(int(in.PerPage), int(in.Page), &model.LogFilter{})
	reply = &LogReply{Entries: logs}
	if appErr != nil {
		err = appErr.Unwrap()
	}
	mlog.Debug("GetLogs-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) QueryLogs(ctx context.Context, in *LogRequest) (reply *QueryLogReply, err error) {
	mlog.Debug("QueryLogs-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	logs, appErr := g2s.ps.GetLogsSkipSend(int(in.PerPage), int(in.Page), &model.LogFilter{ServerNames: []string{g2s.cds.Hostname}})
	reply = &QueryLogReply{Entries: logs}
	if appErr != nil {
		err = appErr.Unwrap()
	}
	mlog.Debug("QueryLogs-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) GetPluginStatuses(ctx context.Context, in *Void) (reply *PluginStatusReply, err error) {
	mlog.Debug("GetPluginStatuses-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
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
	}
	mlog.Debug("GetPluginStatuses-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) ConfigChanged(ctx context.Context, in *ConfigUpdateRequest) (reply *Void, err error) {
	mlog.Debug("ConfigChanged-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	var newConfig model.Config
	if err = json.NewDecoder(bytes.NewBuffer(in.NewConfigBuffer)).Decode(&newConfig); err != nil {
		return
	}
	_, _, err = g2s.ps.GetConfigStore().Set(&newConfig)
	mlog.Debug("ConfigChanged-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}

func (g2s *GossipService) WebConnCountForUser(ctx context.Context, in *WebsocketCountForUserRequest) (reply *WebsocketCountForUserReply, err error) {
	mlog.Debug("WebConnCountForUser-Begin", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("request", in))
	reply = &WebsocketCountForUserReply{
		Count: int32(g2s.ps.WebConnCountForUser(in.UserId)),
	}
	mlog.Debug("WebConnCountForUser-Complete", mlog.String("hostname", g2s.cds.Hostname), mlog.Any("reply", reply), mlog.Err(err))
	return
}
