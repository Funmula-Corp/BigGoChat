// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package post_persistent_notifications

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

type AppIface interface {
	SendPersistentNotifications() error
	IsPersistentNotificationsEnabled() bool
}

func MakeWorker(jobServer *jobs.JobServer, app AppIface) *jobs.SimpleWorker {
	const workerName = "PostPersistentNotifications"

	isEnabled := func(_ *model.Config) bool {
		return app.IsPersistentNotificationsEnabled()
	}
	execute := func(logger mlog.LoggerIFace, job *model.Job) error {
		defer jobServer.HandleJobPanic(logger, job)
		return app.SendPersistentNotifications()
	}
	worker := jobs.NewSimpleWorker(workerName, jobServer, execute, isEnabled)
	return worker
}
