// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package jobs

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

type IndexerJobInterface interface {
	MakeWorker() model.Worker
}
