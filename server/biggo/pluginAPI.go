package biggo

import (
	"context"
	"sync"
	"sync/atomic"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"

	"git.biggo.com/Funmula/mattermost-packages/grpc/grpcServer"
	"git.biggo.com/Funmula/mattermost-packages/pluginAPI"
)

type PluginAPIService struct {
	pluginAPI.UnsafePluginAPIServer
	ps *platform.PlatformService

	running atomic.Bool
	mx      sync.Mutex
}

var instance *PluginAPIService

func InitPluginAPI(ps *platform.PlatformService) {
	if instance == nil && ps != nil {
		instance = &PluginAPIService{ps: ps}
		instance.Start()
	}
}

func (s *PluginAPIService) Start() (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if !s.running.Load() {
		const (
			endpointTCP   string = "localhost:6222"
			endpointHTTP2 string = "localhost:6443"
		)

		var server *grpcServer.GRPCServer
		if server, err = grpcServer.NewGRPCServer[pluginAPI.PluginAPIServer](s, pluginAPI.RegisterPluginAPIServer, endpointTCP, endpointHTTP2); err != nil {
			return
		}

		go func() {
			defer s.running.Store(false)
			s.running.Store(true)
			server.ListenAndServe(true, true)
		}()
	}
	return
}

func (s *PluginAPIService) GetUserIdByAuthData(ctx context.Context, request *pluginAPI.UserIdByAuthDataRequest) (reply *pluginAPI.UserIdByAuthDataReply, err error) {
	var user *model.User
	if user, err = s.ps.Store.User().GetByAuth(&request.AuthData, request.AuthService); err == nil {
		reply = &pluginAPI.UserIdByAuthDataReply{
			UserId: user.Id,
		}
	}
	return
}
