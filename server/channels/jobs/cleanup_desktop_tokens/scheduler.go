// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package cleanup_desktop_tokens

import (
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

const schedFreq = 1 * time.Hour

func MakeScheduler(jobServer *jobs.JobServer) *jobs.PeriodicScheduler {
	isEnabled := func(cfg *model.Config) bool {
		return true
	}
	return jobs.NewPeriodicScheduler(jobServer, model.JobTypeCleanupDesktopTokens, schedFreq, isEnabled)
}
