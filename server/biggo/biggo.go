package biggo

import (
	"github.com/mattermost/mattermost/server/v8/biggo/gossip"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

func Cluster(ps *platform.PlatformService) (cluster einterfaces.ClusterInterface) {
	cluster = &BiggoCluster{ps: ps}
	cluster.(*BiggoCluster).g2s = gossip.NewGossipService(cluster, ps)
	return
}
