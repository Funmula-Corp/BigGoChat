package cluster

import (
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/plugin"
)

func ExampleSchedule() {
	// Use p.API from your plugin instead.
	pluginAPI := plugin.API(nil)

	callback := func() {
		// periodic work to do
	}

	job, err := Schedule(pluginAPI, "key", MakeWaitForInterval(5*time.Minute), callback)
	if err != nil {
		panic("failed to schedule job")
	}

	// main thread

	defer job.Close()
}
