package cluster_test

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/plugin"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/pluginapi/cluster"
)

//nolint:staticcheck
func ExampleMutex() {
	// Use p.API from your plugin instead.
	pluginAPI := plugin.API(nil)

	m, err := cluster.NewMutex(pluginAPI, "key")
	if err != nil {
		panic(err)
	}
	m.Lock()
	// critical section
	m.Unlock()
}
