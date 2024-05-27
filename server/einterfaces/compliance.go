// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

type ComplianceInterface interface {
	StartComplianceDailyJob()
	RunComplianceJob(rctx request.CTX, job *model.Compliance) *model.AppError
}
