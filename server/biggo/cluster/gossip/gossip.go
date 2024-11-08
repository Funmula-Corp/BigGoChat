package gossip

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip/client"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip/server"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/proto"
)

type ClusterClient = proto.ClusterClient
type ClusterServer = server.GossipServer

var (
	NewClusterClient     = client.NewClusterClient
	NewClusterConnection = client.NewClusterConnection
	NewClusterServer     = server.NewClusterServer
)
