// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

func InitTest(logger mlog.LoggerIFace) {
	initStores(logger)
}

func TearDownTest() {
	tearDownStores()
}
