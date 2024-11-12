package cluster

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	CancelVote              context.CancelFunc
	ClusterDiscoveryService *platform.ClusterDiscoveryService
	GossipServer            *gossip.ClusterServer
	KubeConfig              *rest.Config
	KubeClient              *kubernetes.Clientset
	Leader                  atomic.Bool
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
	return p.Leader.Load()
}

func (p *BiggoCluster) RegisterClusterMessageHandler(event model.ClusterEvent, callback einterfaces.ClusterMessageHandler) {
	if p.GossipServer != nil {
		p.GossipServer.ClusterMessageHandler[event] = callback
	}
}

func (p *BiggoCluster) StartInterNodeCommunication() {
	p.LoadConfigFromDB()
	if p.ClusterDiscoveryService != nil {
		if err := p.GossipServer.Start(uint16(p.ClusterDiscoveryService.GossipPort)); err == nil {
			p.ClusterDiscoveryService.Start()
			if p.KubeConfig, err = p.getKubeConfig(); err == nil {
				go p.JoinVote(p.GetClusterId(), "biggochat-leader", "default")
			} else {
				p.Leader.Swap(true)
			}
		}
	}
}

func (p *BiggoCluster) StopInterNodeCommunication() {
	if p.ClusterDiscoveryService != nil {
		p.ClusterDiscoveryService.Stop()
		p.GossipServer.Stop()
		if p.CancelVote != nil {
			p.CancelVote()
			p.CancelVote = nil
		}
	}
}

// cluster info exchange
func (p *BiggoCluster) GetMyClusterInfo() *model.ClusterInfo {
	configHash := func() string {
		buffer := bytes.NewBuffer([]byte{})
		json.NewEncoder(buffer).Encode(p.PlatformService.Config())

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
	info.Hostname, _ = os.Hostname()
	info.Version = dbVersion
	info.SchemaVersion = dbSchemaVersion
	return info
}

// unknown
func (p *BiggoCluster) NotifyMsg(buf []byte) {}

// attempt to load the config form teh database (called by the cluster member on join)
func (p *BiggoCluster) LoadConfigFromDB() {
	const settingName string = "SystemConfigFile"
	var (
		err     error
		setting *model.System
	)

	if setting, err = p.PlatformService.Store.System().GetByName(settingName); err != nil {
		mlog.Error("cluster.config.load.error", mlog.Err(err))
		return
	}

	buffer := bytes.NewBuffer([]byte(setting.Value))
	if err = json.NewDecoder(buffer).Decode(p.PlatformService.Config()); err != nil {
		mlog.Error("cluster.config.decoder.error", mlog.Err(err))
		return
	}

	if _, _, appErr := p.PlatformService.SaveConfig(p.PlatformService.Config(), false); appErr != nil {
		mlog.Error("cluster.config.save.error", mlog.Err(appErr))
		return
	}
}

// persist config to the database (should only be called by the cluster-leader)
func (p *BiggoCluster) SaveConfigToDB() {
	const settingName string = "SystemConfigFile"
	var (
		err     error
		setting *model.System
	)

	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(p.PlatformService.Config()); err != nil {
		mlog.Error("cluster.config.encoder.error", mlog.Err(err))
		return
	}

	if setting, err = p.PlatformService.Store.System().GetByName(settingName); err != nil {
		setting = &model.System{Name: settingName}
	}

	setting.Value = buffer.String()
	if err = p.PlatformService.Store.System().SaveOrUpdate(setting); err != nil {
		mlog.Error("cluster.config.save.error", mlog.Err(err))
		return
	}
}
