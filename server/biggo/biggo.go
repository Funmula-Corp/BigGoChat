package biggo

import (
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

var cluster einterfaces.ClusterInterface

func Cluster(ps *platform.PlatformService) einterfaces.ClusterInterface {
	if cluster == nil {
		cluster = &biggoCluster{}
	}
	return cluster
}
