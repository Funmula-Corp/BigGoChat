package biggo

import (
	"context"
	"fmt"
	"net"
	"net/rpc"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/biggo_dbg"
)

type GossipService struct {
	gCtx    context.Context
	gCancel context.CancelFunc
}

type GossipServiceArgs struct {
	Message *model.ClusterMessage
}

type GossipServiceReply struct{}

func (gs *GossipService) Receive(args *GossipServiceArgs, _ *GossipServiceReply) (err error) {
	if args.Message != nil {
		if f, ok := cluster.(*biggoCluster).cbMap[args.Message.Event]; ok {
			biggo_dbg.Trace("====FOUND====")
			go f(args.Message)
		}
	}
	return
}

func (gs *GossipService) Publish(host string, port *int, msg *model.ClusterMessage) (err error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, *port))
	if err != nil {
		mlog.Error("Gossip IPC Connection Error", logr.Err(err))
		return
	}
	defer client.Close()

	args := GossipServiceArgs{Message: msg}
	if err = client.Call("GossipService.Receive", &args, new(GossipServiceReply)); err != nil {
		mlog.Error("Gossip IPC Call Error", logr.Err(err))
	}
	return
}

func (gs *GossipService) Start(port *int) {
	if gs.gCtx == nil && gs.gCancel == nil {
		gs.gCtx, gs.gCancel = context.WithCancel(context.Background())

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

		go func(ctx context.Context, port *int) {
			if srv, err := net.Listen("tcp", fmt.Sprintf(":%d", *port)); err != nil {
				mlog.Error("Gossip IPC Socket Error", logr.Err(err))
				return
			} else {
				defer srv.Close()
				hdlr(srv, ctx)
			}
		}(gs.gCtx, port)
	}
}

func (gs *GossipService) Stop() {
	if gs.gCtx != nil && gs.gCancel != nil {
		gs.gCancel()
		gs.gCtx, gs.gCancel = nil, nil
	}
}
