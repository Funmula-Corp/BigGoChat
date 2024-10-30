package cluster

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/proto"
)

// cluster information exchange
func (p *BiggoCluster) GetClusterInfos() []*model.ClusterInfo {
	result := []*model.ClusterInfo{p.GetMyClusterInfo()}
	mutex := sync.Mutex{}

	request := &proto.GetClusterInfosRequest{}
	p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.GetClusterInfos(context.Background(), request); err == nil {
			if len(resp.Buffer) == 0 {
				return
			}
			info := &model.ClusterInfo{}
			gob.NewDecoder(bytes.NewBuffer(resp.Buffer)).Decode(info)

			mutex.Lock()
			defer mutex.Unlock()
			result = append(result, info)
		} else {
			mlog.Error("cluster.client.GetClusterInfos.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true)

	if len(result) > 1 {
		slices.SortFunc(result, func(a, b *model.ClusterInfo) int {
			return strings.Compare(a.Hostname, b.Hostname)
		})
	}

	return result
}

// cluster message exchange
func (p *BiggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	request := &proto.SendClusterMessageToNodeRequest{Event: string(msg.Event), Data: msg.Data, Props: msg.Props}
	p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if _, err := client.SendClusterMessageToNode(context.Background(), request); err != nil {
			mlog.Error("cluster.client.SendClusterMessage.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true)
}

func (p *BiggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) error {
	request := &proto.SendClusterMessageToNodeRequest{Event: string(msg.Event), Data: msg.Data, Props: msg.Props}
	return p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if _, err := client.SendClusterMessageToNode(context.Background(), request); err != nil {
			mlog.Error("cluster.client.SendClusterMessageToNode.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true, nodeID)
}

// cluster status exchange
func (p *BiggoCluster) GetClusterStats() ([]*model.ClusterStats, *model.AppError) {
	result := []*model.ClusterStats{}
	mutex := sync.Mutex{}

	request := &proto.GetClusterStatsRequest{}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.GetClusterStats(context.Background(), request); err == nil {
			info := &model.ClusterStats{
				Id:                        node.Id,
				TotalWebsocketConnections: int(resp.GetTotalWebsocketConnections()),
				TotalReadDbConnections:    int(resp.GetTotalReadDbConnections()),
				TotalMasterDbConnections:  int(resp.GetTotalMasterDbConnections()),
			}

			mutex.Lock()
			defer mutex.Unlock()
			result = append(result, info)
		} else {
			mlog.Error("cluster.client.GetClusterStats.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return result, model.NewAppError("cluster.client.GetClusterStats", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (p *BiggoCluster) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	result := model.PluginStatuses{}
	mutex := sync.Mutex{}

	request := &proto.GetPluginStatusesRequest{}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.GetPluginStatuses(context.Background(), request); err == nil {

			mutex.Lock()
			defer mutex.Unlock()
			for idx := range resp.PluginStatus {
				result = append(result, &model.PluginStatus{
					ClusterId:   resp.PluginStatus[idx].ClusterId,
					PluginId:    resp.PluginStatus[idx].PluginId,
					PluginPath:  resp.PluginStatus[idx].PluginPath,
					State:       int(resp.PluginStatus[idx].State),
					Error:       resp.PluginStatus[idx].Error,
					Name:        resp.PluginStatus[idx].Name,
					Description: resp.PluginStatus[idx].Description,
					Version:     resp.PluginStatus[idx].Version,
				})
			}

		} else {
			mlog.Error("cluster.client.GetPluginStatuses.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return result, model.NewAppError("cluster.client.GetPluginStatuses", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (p *BiggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	if !sendToOtherServer || !p.IsLeader() {
		return nil
	}

	newConfBuffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(newConfBuffer).Encode(newConfig); err != nil {
		return model.NewAppError("cluster.client.ConfigChanged", "cluster.serialize.error", nil, err.Error(), http.StatusInternalServerError)
	}

	request := &proto.ConfigChangedRequest{ConfigBuffer: newConfBuffer.Bytes()}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if _, err := client.ConfigChanged(context.Background(), request); err != nil {
			mlog.Error("cluster.client.ConfigChanged.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return model.NewAppError("cluster.client.ConfigChanged", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (p *BiggoCluster) WebConnCountForUser(userID string) (int, *model.AppError) {
	result := atomic.Int64{}

	request := &proto.WebConnCountForUserRequest{UserID: userID}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.WebConnCountForUser(context.Background(), request); err == nil {
			result.Add(resp.Count)
		} else {
			mlog.Error("cluster.client.WebConnCountForUser.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return int(result.Load()), model.NewAppError("cluster.client.WebConnCountForUser", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return int(result.Load()), nil
}

// log exchange
func (p *BiggoCluster) GetLogs(page, perPage int) ([]string, *model.AppError) {
	result := []string{}
	mutex := sync.Mutex{}

	request := &proto.GetLogsRequest{Page: int64(page), PerPage: int64(perPage)}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.GetLogs(context.Background(), request); err == nil {

			if len(resp.LogRecord) > 0 {
				mutex.Lock()
				defer mutex.Unlock()
				result = append(result, resp.LogRecord...)
			}

		} else {
			mlog.Error("cluster.client.GetLogs.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return result, model.NewAppError("cluster.client.GetLogs", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (p *BiggoCluster) QueryLogs(page, perPage int, logFilter *model.LogFilter) (map[string][]string, *model.AppError) {
	result := map[string][]string{}

	request := &proto.QueryLogsRequest{Page: int64(page), PerPage: int64(perPage), LogFilter: &proto.QueryLogsLogFilterRequest{
		ServerNames: logFilter.ServerNames,
		LogLevels:   logFilter.LogLevels,
		DateFrom:    logFilter.DateFrom,
		DateTo:      logFilter.DateTo,
	}}
	if err := p.Call(func(client gossip.ClusterClient, node *model.ClusterDiscovery) {
		if resp, err := client.QueryLogs(context.Background(), request); err == nil {

			if len(resp.LogRecord) > 0 {
				result[node.Id] = resp.LogRecord
			}

		} else {
			mlog.Error("cluster.client.QueryLogs.error", mlog.String("node_id", node.Id), mlog.Err(err))
		}
	}, true); err != nil {
		return result, model.NewAppError("cluster.client.QueryLogs", "cluster.request.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return result, nil
}
