package biggo

import (
	"context"
	"fmt"
	"net"
	"net/rpc"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
)

type GossipService struct{}
type GossipServiceArgs struct {
	Message *model.ClusterMessage
}
type GossipServiceReply struct{}

var gsInstance *GossipService

func (gs *GossipService) RecvMessage(args *GossipServiceArgs, _ *GossipServiceReply) (err error) {
	if f, ok := cluster.(*biggoCluster).msgHdlr[args.Message.Event]; ok {
		go f(args.Message)
	}
	return
}

func (c *biggoCluster) sendMessage(rcvr string, msg *model.ClusterMessage) (err error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", rcvr, *c.ps.Config().ClusterSettings.GossipPort))
	if err != nil {
		mlog.Error("Gossip IPC Connection Error", logr.Err(err))
		return
	}
	defer client.Close()

	args := GossipServiceArgs{Message: msg}
	if err = client.Call("GossipService.RecvMessage", &args, new(GossipServiceReply)); err != nil {
		mlog.Error("Gossip IPC Call Error", logr.Err(err))
	}
	return
}

func (c *biggoCluster) startInterNodeListener() {
	if c.gCtx == nil && c.gCancel == nil {
		if gsInstance == nil {
			gsInstance = new(GossipService)
			rpc.Register(gsInstance)
		}

		c.gCtx, c.gCancel = context.WithCancel(context.Background())
		port := c.ps.Config().ClusterSettings.GossipPort

		hdlr := func(srv net.Listener, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					conn, err := srv.Accept()
					if err != nil {
						mlog.Error("Gossip IPC Accept Error", logr.Err(err))
						continue
					}
					go rpc.ServeConn(conn)
				}
			}
		}

		go func(ps *platform.PlatformService, ctx context.Context, port *int) {
			if srv, err := net.Listen("tcp", fmt.Sprintf(":%d", *port)); err != nil {
				mlog.Error("Gossip IPC Socket Error", logr.Err(err))
				return
			} else {
				defer srv.Close()
				hdlr(srv, ctx)
			}
		}(c.ps, c.gCtx, port)
	}
}

func (c *biggoCluster) stopInterNodeListener() {
	if c.gCtx != nil && c.gCancel != nil {
		c.gCancel()
		c.gCtx, c.gCancel = nil, nil
	}
}
