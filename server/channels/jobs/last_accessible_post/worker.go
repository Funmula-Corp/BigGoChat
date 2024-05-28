// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package last_accessible_post

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

type AppIface interface {
	ComputeLastAccessiblePostTime() error
}

func MakeWorker(jobServer *jobs.JobServer, license *model.License, app AppIface) *jobs.SimpleWorker {
	const workerName = "LastAccessiblePost"

	isEnabled := func(_ *model.Config) bool {
		return license != nil && license.Features != nil && *license.Features.Cloud
	}
	execute := func(logger mlog.LoggerIFace, job *model.Job) error {
		defer jobServer.HandleJobPanic(logger, job)

		return app.ComputeLastAccessiblePostTime()
	}
	worker := jobs.NewSimpleWorker(workerName, jobServer, execute, isEnabled)
	return worker
}
