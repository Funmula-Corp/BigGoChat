package biggo

import (
	context "context"
	"net"
	"sync"
	"sync/atomic"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/mattermost-packages/pluginAPI"

	"google.golang.org/grpc"
)

type PluginAPIService struct {
	pluginAPI.UnimplementedPluginAPIServer
	ps *platform.PlatformService

	running atomic.Bool

	grpcSoc net.Listener
	grpcSvr *grpc.Server

	mx sync.Mutex
}

var instance *PluginAPIService

func InitPluginAPI(ps *platform.PlatformService) {
	if instance == nil {
		instance = &PluginAPIService{ps: ps}
		instance.Start()
	}
}

func (s *PluginAPIService) Start() (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if !s.running.Load() {
		if s.grpcSoc, err = net.Listen("tcp", "localhost:9999"); err != nil {
			return
		}
		s.grpcSvr = grpc.NewServer()

		go func() {
			defer s.running.Store(false)
			defer s.grpcSoc.Close()
			s.running.Store(true)
			s.grpcSvr.Serve(s.grpcSoc)
		}()
	}
	return
}

func (s *PluginAPIService) GetUserIdByAuthData(ctx context.Context, in *pluginAPI.UserIdByAuthDataRequest) (reply *pluginAPI.UserIdByAuthDataReply, err error) {
	var user *model.User
	if user, err = s.ps.Store.User().GetByAuth(&in.AuthData, "service"); err != nil {
		return
	}
	reply = &pluginAPI.UserIdByAuthDataReply{
		UserId: user.Id,
	}
	return
}
