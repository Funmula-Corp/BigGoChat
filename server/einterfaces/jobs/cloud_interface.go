// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package jobs

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

type CloudJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() Scheduler
}
