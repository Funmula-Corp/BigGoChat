package cluster

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"
)

func New(platformService *platform.PlatformService) *BiggoCluster {
	cluster := &BiggoCluster{
		ClusterDiscoveryService: platformService.NewClusterDiscoveryService(),
		GossipServer:            gossip.NewClusterServer(platformService),
		PlatformService:         platformService,
	}

	if *platformService.Config().ClusterSettings.OverrideHostname != "" {
		cluster.ClusterDiscoveryService.Hostname = *platformService.Config().ClusterSettings.OverrideHostname
	} else if *platformService.Config().ClusterSettings.UseIPAddress {
		cluster.ClusterDiscoveryService.AutoFillIPAddress(
			*platformService.Config().ClusterSettings.NetworkInterface,
			*platformService.Config().ClusterSettings.AdvertiseAddress,
		)
	} else {
		cluster.ClusterDiscoveryService.AutoFillHostname()
	}

	cluster.ClusterDiscoveryService.ClusterName = *cluster.PlatformService.Config().ClusterSettings.ClusterName
	cluster.ClusterDiscoveryService.GossipPort = (int32)(*cluster.PlatformService.Config().ClusterSettings.GossipPort)
	cluster.ClusterDiscoveryService.Type = model.GetServiceEnvironment()
	return cluster
}

type BiggoCluster struct {
	ClusterDiscoveryService *platform.ClusterDiscoveryService
	GossipServer            *gossip.ClusterServer
	PlatformService         *platform.PlatformService
}

// cluster meta functions
func (p *BiggoCluster) GetClusterId() string {
	if p.ClusterDiscoveryService != nil {
		return p.ClusterDiscoveryService.Id
	}
	return "UNKNOWN"
}

func (p *BiggoCluster) HealthScore() int {
	return 0
}

func (p *BiggoCluster) IsLeader() bool {
	return true
}

func (p *BiggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, callback einterfaces.ClusterMessageHandler) {
	if p.GossipServer != nil {
		p.GossipServer.ClusterMessageHandler[event] = callback
	}
}

func (p *BiggoCluster) StartInterNodeCommunication() {
	if p.ClusterDiscoveryService != nil {
		if err := p.GossipServer.Start(uint16(p.ClusterDiscoveryService.GossipPort)); err == nil {
			p.ClusterDiscoveryService.Start()
		}
	}
}

func (p *BiggoCluster) StopInterNodeCommunication() {
	if p.ClusterDiscoveryService != nil {
		p.ClusterDiscoveryService.Stop()
	}
}

// cluster info exchange
func (p *BiggoCluster) GetMyClusterInfo() *model.ClusterInfo {
	configHash := func() string {
		buffer := bytes.NewBuffer([]byte{})
		gob.NewEncoder(buffer).Encode(p.PlatformService.Config())

		hash := md5.New()
		hash.Write(buffer.Bytes())
		return hex.EncodeToString(hash.Sum(nil))
	}()

	dbVersion := func() string {
		if version, err := p.PlatformService.Store.GetDbVersion(true); err != nil {
			mlog.Error("Cluster Info Error", mlog.Err(err))
			return "ERROR"
		} else {
			return version
		}
	}()

	dbSchemaVersion := func() string {
		if version, err := p.PlatformService.Store.GetDBSchemaVersion(); err != nil {
			mlog.Error("Cluster Info Error", mlog.Err(err))
			return "ERROR"
		} else {
			return fmt.Sprintf("%d", version)
		}
	}()

	info := &model.ClusterInfo{}
	info.Id = p.ClusterDiscoveryService.Id
	info.ConfigHash = configHash
	info.IPAddress = p.ClusterDiscoveryService.Hostname
	info.Hostname = p.ClusterDiscoveryService.Hostname
	info.Version = dbVersion
	info.SchemaVersion = dbSchemaVersion
	return info
}

// unknown
func (p *BiggoCluster) NotifyMsg(buf []byte) {}
