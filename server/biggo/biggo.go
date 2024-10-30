package biggo

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/gossip"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"
)

func Cluster(ps *platform.PlatformService) (cluster einterfaces.ClusterInterface) {
	cluster = &BiggoCluster{ps: ps}
	cluster.(*BiggoCluster).g2s = gossip.NewGossipService(cluster, ps)
	InitPluginAPI(ps)
	return
}
