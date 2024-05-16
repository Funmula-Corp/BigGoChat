package biggo

import (
	"net/rpc"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

func Cluster(ps *platform.PlatformService) (cluster einterfaces.ClusterInterface) {
	cluster = &BiggoCluster{
		ps: ps, cds: ps.NewClusterDiscoveryService(),
		cbMap: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
	}
	cluster.(*BiggoCluster).gService = &GossipService{c: cluster.(*BiggoCluster)}
	rpc.RegisterName("GossipService", cluster.(*BiggoCluster).gService)
	return
}
