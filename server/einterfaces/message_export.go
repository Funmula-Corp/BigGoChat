// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type MessageExportInterface interface {
	StartSynchronizeJob(c request.CTX, exportFromTimestamp int64) (*model.Job, *model.AppError)
	RunExport(c request.CTX, format string, since int64, limit int) (int64, *model.AppError)
}
