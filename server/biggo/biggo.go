package biggo

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

func Cluster(ps *platform.PlatformService) (cluster einterfaces.ClusterInterface) {
	cluster = &BiggoCluster{
		ps: ps, cds: ps.NewClusterDiscoveryService(),
		cbMap: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
	}
	cluster.(*BiggoCluster).g2Service = &G2Service{cluster: cluster.(*BiggoCluster)}
	return
}
