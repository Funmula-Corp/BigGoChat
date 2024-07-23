package gossip

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/einterfaces"

	"google.golang.org/grpc"
)

func NewGossipService(ci einterfaces.ClusterInterface, ps *platform.PlatformService) *GossipService {
	return &GossipService{
		ci: ci, ps: ps, cbMap: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
	}
}

type GossipService struct {
	UnimplementedClusterServer

	ci  einterfaces.ClusterInterface
	cds *platform.ClusterDiscoveryService
	ps  *platform.PlatformService

	running atomic.Bool

	cbMap map[model.ClusterEvent]einterfaces.ClusterMessageHandler

	grpcSoc net.Listener
	grpcSvr *grpc.Server

	mx sync.Mutex
}

func (g2s *GossipService) GetClusterDiscoveryService() *platform.ClusterDiscoveryService {
	return g2s.cds
}

func (g2s *GossipService) NewClusterDiscoveryService() {
	if g2s.cds != nil {
		g2s.cds.Stop()
	}

	g2s.cds = g2s.ps.NewClusterDiscoveryService()
	if *g2s.ps.Config().ClusterSettings.OverrideHostname != "" {
		g2s.cds.Hostname = *g2s.ps.Config().ClusterSettings.OverrideHostname
	} else if *g2s.ps.Config().ClusterSettings.UseIPAddress {
		g2s.cds.AutoFillIPAddress(
			*g2s.ps.Config().ClusterSettings.NetworkInterface,
			*g2s.ps.Config().ClusterSettings.AdvertiseAddress,
		)
	} else {
		g2s.cds.AutoFillHostname()
	}

	g2s.cds.ClusterName = *g2s.ps.Config().ClusterSettings.ClusterName
	g2s.cds.GossipPort = (int32)(*g2s.ps.Config().ClusterSettings.GossipPort)
	g2s.cds.Type = model.GetServiceEnvironment()
}

func (g2s *GossipService) RegisterClusterMessageHandler(event model.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	g2s.cbMap[event] = crm
}

func (g2s *GossipService) StartInterNodeCommunication() (err error) {
	g2s.mx.Lock()
	defer g2s.mx.Unlock()

	if !g2s.running.Load() {
		g2s.NewClusterDiscoveryService()
		if g2s.grpcSoc, err = net.Listen("tcp", fmt.Sprintf(":%d", g2s.cds.GossipPort)); err != nil {
			return
		}

		g2s.grpcSvr = grpc.NewServer()
		RegisterClusterServer(g2s.grpcSvr, g2s)

		go func() {
			defer g2s.running.Store(false)
			defer g2s.grpcSoc.Close()
			defer g2s.cds.Stop()
			g2s.running.Store(true)
			g2s.cds.Start()
			g2s.grpcSvr.Serve(g2s.grpcSoc)
		}()
	}
	return
}

func (g2s *GossipService) StopInterNodeCommunication() {
	g2s.mx.Lock()
	defer g2s.mx.Unlock()

	if g2s.running.Load() {
		g2s.grpcSvr.Stop()
	}
}
