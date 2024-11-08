package server

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/proto"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"

	"google.golang.org/grpc"
)

var (
	ErrMissingRequest error = errors.New("request is nil")
)

func NewClusterServer(platformService *platform.PlatformService) *GossipServer {
	server := &GossipServer{
		ClusterMessageHandler: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
		grpcServer:            grpc.NewServer(),
		platformService:       platformService,
	}
	proto.RegisterClusterServer(server.grpcServer, server)
	return server
}

type GossipServer struct {
	proto.ClusterServer

	ClusterMessageHandler map[model.ClusterEvent]einterfaces.ClusterMessageHandler

	grpcServer        *grpc.Server
	grpcServerRunning atomic.Bool
	platformService   *platform.PlatformService
}

func (p *GossipServer) Start(port uint16) (err error) {
	if !p.grpcServerRunning.Swap(true) {
		var listener net.Listener
		if listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
			return
		}
		go p.grpcServer.Serve(listener)
	}
	return
}

func (p *GossipServer) Stop() {
	if p.grpcServerRunning.Swap(false) {
		p.grpcServer.Stop()
	}
}
