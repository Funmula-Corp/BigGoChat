// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package post_persistent_notifications

import (
	"time"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/jobs"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

type Scheduler struct {
	*jobs.PeriodicScheduler
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, _ time.Time, _ bool, _ *model.Job) *time.Time {
	nextTime := time.Now().Add((time.Duration(*cfg.ServiceSettings.PersistentNotificationIntervalMinutes) * time.Minute) / 2)
	return &nextTime
}

func MakeScheduler(jobServer *jobs.JobServer, licenseFunc func() *model.License) *Scheduler {
	enabledFunc := func(_ *model.Config) bool {
		l := licenseFunc()
		return l != nil && (l.SkuShortName == model.LicenseShortSkuProfessional || l.SkuShortName == model.LicenseShortSkuEnterprise)
	}
	return &Scheduler{jobs.NewPeriodicScheduler(jobServer, model.JobTypePostPersistentNotifications, 0, enabledFunc)}
}
