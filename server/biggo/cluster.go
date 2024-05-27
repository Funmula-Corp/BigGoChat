package biggo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/biggo/gossip"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

// distribute users across the MM cluster based on userID & consistent hashing via ISTIO proxy

type BiggoCluster struct {
	ps  *platform.PlatformService
	g2s *gossip.GossipService
}

func (c *BiggoCluster) StartInterNodeCommunication() {
	c.g2s.StartInterNodeCommunication()
}

func (c *BiggoCluster) StopInterNodeCommunication() {
	c.g2s.StopInterNodeCommunication()
}

func (c *BiggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	c.g2s.RegisterClusterMessageHandler(event, crm)
}

func (c *BiggoCluster) GetClusterId() string {
	if c.g2s.GetClusterDiscoveryService() == nil {
		return ""
	}
	return c.g2s.GetClusterDiscoveryService().Id
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
	info.Id = c.g2s.GetClusterDiscoveryService().Id
	info.ConfigHash = hex.EncodeToString(hash.Sum(nil))
	info.IPAddress = c.g2s.GetClusterDiscoveryService().Hostname
	info.Hostname, _ = os.Hostname()
	info.Version = dbVersion
	info.SchemaVersion = dbSchemaVersion
	return
}

func (c *BiggoCluster) GetClusterInfos() (result []*model.ClusterInfo) {
	mx := sync.Mutex{}
	result = []*model.ClusterInfo{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallGetMyClusterInfo(hostname); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result = append(result, res)
		}
	}, false)
	return
}

/*
func (c *BiggoCluster) GetClusterInfos() (result []*model.ClusterInfo) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	result = make([]*model.ClusterInfo, 0)
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				mx.Lock()
				result = append(result, c.GetMyClusterInfo())
				mx.Unlock()
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()

			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	c.g2s.CallCluster(func(hostname string) {
		c.g2s.CallSendClusterMessageToNode(hostname, msg)
	}, true)
}

/*
func (c *BiggoCluster) SendClusterMessage(msg *model.ClusterMessage) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if !cds.IsEqual(cd) {
				wg.Add(1)
				go func(nodeID, hostname string) {
					if err := c.SendClusterMessageToNode(nodeID, msg); err != nil {
						mlog.Error("Send Cluster Message Error", logr.String("hostname", hostname))
					}
				}(cd.Id, cd.Hostname)
			}
		}
		wg.Wait()
	}
}
*/

func (c *BiggoCluster) SendClusterMessageToNode(nodeID string, msg *model.ClusterMessage) (err error) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) || cd.Id != nodeID {
				continue
			}

			wg.Add(1)
			go func(hostname string, msg *model.ClusterMessage) {
				defer wg.Done()
				c.g2s.CallSendClusterMessageToNode(hostname, msg)
			}(cd.Hostname, msg)
		}
		wg.Wait()
	}
	return
}

func (c *BiggoCluster) NotifyMsg(buf []byte) {
	mlog.Error("Cluster NotifyMsg Call Error", logr.Err(errors.New("NOT IMPLEMENTED")))
}

func (c *BiggoCluster) GetClusterStats() (result []*model.ClusterStats, aErr *model.AppError) {
	mx := sync.Mutex{}
	result = []*model.ClusterStats{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallGetClusterStats(hostname); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result = append(result, res)
		}
	}, true)
	return
}

/*
func (c *BiggoCluster) GetClusterStats() (result []*model.ClusterStats, aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	result = []*model.ClusterStats{}
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if res, err := c.g2s.CallGetClusterStats(hostname); err == nil {
					mx.Lock()
					defer mx.Unlock()
					result = append(result, res)
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) GetLogs(page, perPage int) (result []string, aErr *model.AppError) {
	mx := sync.Mutex{}
	result = []string{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallGetLogs(hostname, page, perPage); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result = append(result, res...)
		}
	}, true)
	return
}

/*
func (c *BiggoCluster) GetLogs(page, perPage int) (result []string, aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	result = []string{}
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if res, err := c.g2s.CallGetLogs(hostname, page, perPage); err == nil {
					mx.Lock()
					defer mx.Unlock()
					result = append(result, res...)
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) QueryLogs(page, perPage int) (result map[string][]string, aErr *model.AppError) {
	mx := sync.Mutex{}
	result = map[string][]string{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallGetLogs(hostname, page, perPage); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result[hostname] = res
		}
	}, true)
	return
}

/*
func (c *BiggoCluster) QueryLogs(page, perPage int) (result map[string][]string, aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	result = map[string][]string{}
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if res, err := c.g2s.CallQueryLogs(hostname, page, perPage); err == nil {
					mx.Lock()
					defer mx.Unlock()
					result[hostname] = res
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) GetPluginStatuses() (result model.PluginStatuses, aErr *model.AppError) {
	mx := sync.Mutex{}
	result = model.PluginStatuses{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallGetPluginStatuses(hostname); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result = append(result, res...)
		}
	}, true)
	return
}

/*
func (c *BiggoCluster) GetPluginStatuses() (result model.PluginStatuses, aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	result = model.PluginStatuses{}
	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if res, err := c.g2s.CallGetPluginStatuses(hostname); err == nil {
					mx.Lock()
					defer mx.Unlock()
					result = append(result, res...)
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) (aErr *model.AppError) {
	c.g2s.CallCluster(func(hostname string) {
		c.g2s.CallConfigChanged(hostname, previousConfig, newConfig)
	}, true)
	return
}

/*
func (c *BiggoCluster) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) (aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if err := c.g2s.CallConfigChanged(hostname, previousConfig, newConfig); err != nil {
					mlog.Error("Get Cluster Config Update Error", mlog.Err(err))
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/

func (c *BiggoCluster) WebConnCountForUser(userID string) (result int, aErr *model.AppError) {
	mx := sync.Mutex{}
	c.g2s.CallCluster(func(hostname string) {
		if res, err := c.g2s.CallWebConnCountForUser(hostname, userID); err == nil {
			mx.Lock()
			defer mx.Unlock()
			result += res
		}
	}, true)
	return
}

/*
func (c *BiggoCluster) WebConnCountForUser(userID string) (count int, aErr *model.AppError) {
	cds := c.g2s.GetClusterDiscoveryService()
	if cds == nil {
		return
	}

	if clusterDiscovery, err := c.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", logr.Err(err))
	} else {
		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) {
				count += c.ps.WebConnCountForUser(userID)
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				if res, err := c.g2s.CallWebConnCountForUser(hostname, userID); err == nil {
					mx.Lock()
					defer mx.Unlock()
					count += res
				}
			}(cd.Hostname)
		}
		wg.Wait()
	}
	return
}
*/
