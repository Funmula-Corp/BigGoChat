package biggo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/biggo_dbg"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

// distribute users across the MM cluster based on userID & consistent hashing via ISTIO proxy

type BiggoCluster struct {
	ps  *platform.PlatformService
	cds *platform.ClusterDiscoveryService

	g2Service *G2Service

	cbMap map[model.ClusterEvent]einterfaces.ClusterMessageHandler
}

func (c *BiggoCluster) StartInterNodeCommunication() {
	c.g2Service.StartInterNodeCommunication()
}

func (c *BiggoCluster) StopInterNodeCommunication() {
	c.g2Service.StopInterNodeCommunication()
}

func (c *BiggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	c.cbMap[event] = crm
}

func (c *BiggoCluster) GetClusterId() string {
	return c.cds.Id
}

func (c *BiggoCluster) IsLeader() bool {
	// TODO: implement
	hostname, _ := os.Hostname()
	return hostname == "biggo-chat-0"
}

func (c *BiggoCluster) HealthScore() int {
	// HealthScore returns a number which is indicative of how well an instance is meeting
	// the soft real-time requirements of the protocol. Lower numbers are better,
	// and zero means "totally healthy".
	return 0
}

func (c *BiggoCluster) GetMyClusterInfo() (info *model.ClusterInfo) {
	buffer := bytes.NewBuffer([]byte{})
	json.NewEncoder(buffer).Encode(c.ps.Config())

	hash := md5.New()
	hash.Write(buffer.Bytes())

	dbVersion := func() string {
		if version, err := c.ps.Store.GetDbVersion(true); err != nil {
			mlog.Error("Cluster Info Error", logr.Err(err))
			return "ERROR"
		} else {
			return version
		}
	}()

	dbSchemaVersion := func() string {
		if version, err := c.ps.Store.GetDBSchemaVersion(); err != nil {
			mlog.Error("Cluster Info Error", logr.Err(err))
			return "ERROR"
		} else {
			return fmt.Sprintf("%d", version)
		}
	}()

	info = new(model.ClusterInfo)
	info.Id = c.cds.Id
	info.ConfigHash = hex.EncodeToString(hash.Sum(nil))
	info.IPAddress = c.cds.Hostname
	info.Hostname, _ = os.Hostname()
	info.Version = dbVersion
	info.SchemaVersion = dbSchemaVersion
	mlog.Info("=====DEBUG=====", logr.Any("cluster_info", info))
	return
}

func (c *BiggoCluster) GetClusterInfos() (result []*model.ClusterInfo) {
	result = make([]*model.ClusterInfo, 0)
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(c.cds.Type, c.cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		for _, cd := range clusterDiscovery {
			if c.cds.IsEqual(cd) {
				result = append(result, c.GetMyClusterInfo())
				continue
			}

			if res, err := c.g2Service.GetClusterInfo(cd.Hostname); res != nil && err == nil {
				result = append(result, res)
			}
		}
	}
	return
}

func (c *BiggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(c.cds.Type, c.cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		for _, cd := range clusterDiscovery {
			if !c.cds.IsEqual(cd) {
				if err := c.g2Service.PostClusterMessage(cd.Hostname, msg); err != nil {
					mlog.Error("SendClusterMessage Error", logr.String("hostname", cd.Hostname))
				}
			}
		}
	}
}

func (c *BiggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) (err error) {
	var clusterDiscovery []*model.ClusterDiscovery
	if clusterDiscovery, err = c.ps.Store.ClusterDiscovery().GetAll(c.cds.Type, c.cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		for _, cd := range clusterDiscovery {
			if cd.Id == nodeID {
				err = c.g2Service.PostClusterMessage(cd.Hostname, msg)
			}
		}
	}
	return
}

func (c *BiggoCluster) NotifyMsg(buf []byte) {
	mlog.Error("Cluster NotifyMsg Call Error", logr.Err(errors.New("NOT IMPLEMENTED")))
}

// TODO: implement
func (c *BiggoCluster) GetClusterStats() ([]*model.ClusterStats, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

// TODO: implement
func (c *BiggoCluster) GetLogs(page, perPage int) ([]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

// TODO: implement
func (c *BiggoCluster) QueryLogs(page, perPage int) (map[string][]string, *model.AppError) {
	biggo_dbg.Trace(page, perPage)
	return nil, nil
}

// TODO: implement
func (c *BiggoCluster) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	biggo_dbg.Trace()
	return nil, nil
}

func (c *BiggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) (aerr *model.AppError) {
	if sendToOtherServer {
		if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(c.cds.Type, c.cds.ClusterName); err != nil {
			mlog.Error("Cluster Discovery Error", logr.Err(err))
		} else {
			for _, cd := range clusterDiscovery {
				if !c.cds.IsEqual(cd) {
					c.g2Service.PostClusterConfig(cd.Hostname, newConfig)
				}
			}
		}
	}
	return
}

// TODO: implement
func (c *BiggoCluster) WebConnCountForUser(userID string) (int, *model.AppError) {
	biggo_dbg.Trace(userID)
	return 0, nil
}
