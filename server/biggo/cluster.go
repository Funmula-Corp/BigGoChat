package biggo

import (
	"errors"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/biggo_dbg"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

// distribute users across the MM cluster based on userID & consistent hashing via ISTIO proxy

type biggoCluster struct {
	ps  *platform.PlatformService
	cds *platform.ClusterDiscoveryService

	gService *GossipService

	cbMap map[model.ClusterEvent]einterfaces.ClusterMessageHandler
}

func (c *biggoCluster) StartInterNodeCommunication() {
	c.gService.Start(c.ps.Config().ClusterSettings.GossipPort)
	c.cds.Start()
}

func (c *biggoCluster) StopInterNodeCommunication() {
	c.cds.Stop()
	c.gService.Stop()
}

func (c *biggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	c.cbMap[event] = crm
}

// TODO: implement
func (c *biggoCluster) GetClusterId() string {
	biggo_dbg.Trace()
	return ""
}

// TODO: implement
func (c *biggoCluster) IsLeader() bool {
	biggo_dbg.Trace()
	return true
}

func (c *biggoCluster) HealthScore() int {
	// HealthScore returns a number which is indicative of how well an instance is meeting
	// the soft real-time requirements of the protocol. Lower numbers are better,
	// and zero means "totally healthy".
	biggo_dbg.Trace()
	return 0
}

// TODO: implement
func (c *biggoCluster) GetMyClusterInfo() *model.ClusterInfo {
	biggo_dbg.Trace()
	return nil
}

// TODO: implement
func (c *biggoCluster) GetClusterInfos() []*model.ClusterInfo {
	biggo_dbg.Trace()
	return nil
}

// TODO: implement gossipping (message distribution to neighbor)
func (c *biggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	biggo_dbg.Trace(msg)
	c.gService.Publish("127.0.0.1", c.ps.Config().ClusterSettings.GossipPort, msg)
}

// TODO: implement
func (c *biggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) error {
	biggo_dbg.Trace(nodeID, msg)
	return nil
}

func (c *biggoCluster) NotifyMsg(buf []byte) {
	mlog.Error("Cluster NotifyMsg Call Error", logr.Err(errors.New("NOT IMPLEMENTED")))
}

// TODO: implement
func (c *biggoCluster) GetClusterStats() ([]*model.ClusterStats, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

// TODO: implement
func (c *biggoCluster) GetLogs(page, perPage int) ([]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

// TODO: implement
func (c *biggoCluster) QueryLogs(page, perPage int) (map[string][]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

// TODO: implement
func (c *biggoCluster) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

// TODO: implement
func (c *biggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	biggo_dbg.Trace(previousConfig, newConfig, sendToOtherServer)
	return nil
}

// TODO: implement
func (c *biggoCluster) WebConnCountForUser(userID string) (int, *model.AppError) {
	biggo_dbg.Trace(userID)
	return 0, nil
}
