// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package expirynotify

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/jobs"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
)

func MakeWorker(jobServer *jobs.JobServer, notifySessionsExpired func() error) *jobs.SimpleWorker {
	const workerName = "ExpiryNotify"

	isEnabled := func(cfg *model.Config) bool {
		return *cfg.ServiceSettings.ExtendSessionLengthWithActivity
	}
	execute := func(logger mlog.LoggerIFace, job *model.Job) error {
		defer jobServer.HandleJobPanic(logger, job)

		return notifySessionsExpired()
	}
	return jobs.NewSimpleWorker(workerName, jobServer, execute, isEnabled)
}
