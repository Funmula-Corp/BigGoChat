package biggo_engine

import (
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
)

const schedFreq = time.Minute

func MakeScheduler(jobServer *jobs.JobServer) *jobs.PeriodicScheduler {
	isEnabled := func(cfg *model.Config) bool {
		return true
	}
	return jobs.NewPeriodicScheduler(jobServer, model.JobTypeBiggoIndexing, schedFreq, isEnabled)
}
