package biggo

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"
)

func Cluster(ps *platform.PlatformService) einterfaces.ClusterInterface {
	InitPluginAPI(ps)
	return cluster.New(ps)
}
