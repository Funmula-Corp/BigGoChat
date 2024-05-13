package biggo

import (
	"net/rpc"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

var cluster einterfaces.ClusterInterface

func Cluster(ps *platform.PlatformService) einterfaces.ClusterInterface {
	if cluster == nil {
		cluster = &biggoCluster{
			ps: ps, cds: ps.NewClusterDiscoveryService(),
			cbMap: map[model.ClusterEvent]einterfaces.ClusterMessageHandler{},
		}
		cluster.(*biggoCluster).cds.ClusterDiscovery.AutoFillHostname()
		cluster.(*biggoCluster).cds.ClusterDiscovery.AutoFillIPAddress(
			*cluster.(*biggoCluster).ps.Config().ClusterSettings.NetworkInterface,
			*cluster.(*biggoCluster).ps.Config().ClusterSettings.AdvertiseAddress,
		)
		cluster.(*biggoCluster).cds.ClusterDiscovery.ClusterName = *cluster.(*biggoCluster).ps.Config().ClusterSettings.ClusterName
		cluster.(*biggoCluster).cds.ClusterDiscovery.GossipPort = (int32)(*cluster.(*biggoCluster).ps.Config().ClusterSettings.GossipPort)
		cluster.(*biggoCluster).cds.ClusterDiscovery.Type = model.GetServiceEnvironment()

		cluster.(*biggoCluster).gService = new(GossipService)
		rpc.Register(cluster.(*biggoCluster).gService)
	}
	return cluster
}
