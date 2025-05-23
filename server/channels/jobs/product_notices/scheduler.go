// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package product_notices

import (
	"time"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/jobs"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

type Scheduler struct {
	*jobs.PeriodicScheduler
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, _ time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	nextTime := time.Now().Add(time.Duration(*cfg.AnnouncementSettings.NoticesFetchFrequency) * time.Second)
	return &nextTime
}

func MakeScheduler(jobServer *jobs.JobServer) *Scheduler {
	isEnabled := func(cfg *model.Config) bool {
		return *cfg.AnnouncementSettings.AdminNoticesEnabled || *cfg.AnnouncementSettings.UserNoticesEnabled
	}
	return &Scheduler{PeriodicScheduler: jobs.NewPeriodicScheduler(jobServer, model.JobTypeProductNotices, 0, isEnabled)}
}
