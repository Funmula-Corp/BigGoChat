package biggo

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

// distribute users across the MM cluster based on userID & consistent hashing

type biggoCluster struct{}

func (c *biggoCluster) StartInterNodeCommunication() {
	Trace()
}

func (c *biggoCluster) StopInterNodeCommunication() {
	Trace()
}

func (c *biggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	Trace(event, crm)
}

func (c *biggoCluster) GetClusterId() string {
	Trace()
	return ""
}

func (c *biggoCluster) IsLeader() bool {
	Trace()
	return true
}

func (c *biggoCluster) HealthScore() int {
	Trace()
	return 0 // 0 is healthy
}

func (c *biggoCluster) GetMyClusterInfo() *model.ClusterInfo {
	Trace()
	return nil
}

func (c *biggoCluster) GetClusterInfos() []*model.ClusterInfo {
	Trace()
	return nil
}

func (c *biggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	Trace(msg)
}

func (c *biggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) error {
	Trace(nodeID, msg)
	return nil
}

func (c *biggoCluster) NotifyMsg(buf []byte) {
	Trace(string(buf))
}

func (c *biggoCluster) GetClusterStats() ([]*model.ClusterStats, *model.AppError) {
	Trace()
	return nil, nil
}

func (c *biggoCluster) GetLogs(page, perPage int) ([]string, *model.AppError) {
	Trace(page, perPage)
	return nil, nil
}

func (c *biggoCluster) QueryLogs(page, perPage int) (map[string][]string, *model.AppError) {
	Trace(page, perPage)
	return nil, nil
}

func (c *biggoCluster) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	Trace()
	return nil, nil
}

func (c *biggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	Trace(previousConfig, newConfig, sendToOtherServer)
	return nil
}

func WebConnCountForUser(userID string) (int, *model.AppError) {
	Trace(userID)
	return 0, nil
}
