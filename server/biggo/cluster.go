package biggo

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/biggo_dbg"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

// distribute users across the MM cluster based on userID & consistent hashing

type biggoCluster struct{}

func (c *biggoCluster) StartInterNodeCommunication() {
	biggo_dbg.Trace()
}

func (c *biggoCluster) StopInterNodeCommunication() {
	biggo_dbg.Trace()
}

func (c *biggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	biggo_dbg.Trace(event, crm)
}

func (c *biggoCluster) GetClusterId() string {
	biggo_dbg.Trace()
	return ""
}

func (c *biggoCluster) IsLeader() bool {
	biggo_dbg.Trace()
	return true
}

func (c *biggoCluster) HealthScore() int {
	biggo_dbg.Trace()
	return 0 // 0 is healthy
}

func (c *biggoCluster) GetMyClusterInfo() *model.ClusterInfo {
	biggo_dbg.Trace()
	return nil
}

func (c *biggoCluster) GetClusterInfos() []*model.ClusterInfo {
	biggo_dbg.Trace()
	return nil
}

func (c *biggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	biggo_dbg.Trace(msg)
}

func (c *biggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) error {
	biggo_dbg.Trace(nodeID, msg)
	return nil
}

func (c *biggoCluster) NotifyMsg(buf []byte) {
	biggo_dbg.Trace(string(buf))
}

func (c *biggoCluster) GetClusterStats() ([]*model.ClusterStats, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

func (c *biggoCluster) GetLogs(page, perPage int) ([]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

func (c *biggoCluster) QueryLogs(page, perPage int) (map[string][]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

func (c *biggoCluster) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

func (c *biggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	biggo_dbg.Trace(previousConfig, newConfig, sendToOtherServer)
	return nil
}

func (c *biggoCluster) WebConnCountForUser(userID string) (int, *model.AppError) {
	biggo_dbg.Trace(userID)
	return 0, nil
}
