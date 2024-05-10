package biggo

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

var cluster einterfaces.ClusterInterface

type msgRegistry map[model.ClusterEvent]einterfaces.ClusterMessageHandler

func Cluster(ps *platform.PlatformService) einterfaces.ClusterInterface {
	if cluster == nil {
		cluster = &biggoCluster{
			ps: ps, msgHdlr: msgRegistry{},
		}
	}
	return cluster
}
