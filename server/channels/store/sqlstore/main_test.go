// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore_test

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/sqlstore"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/testlib"
)

var mainHelper *testlib.MainHelper

func TestMain(m *testing.M) {
	mainHelper = testlib.NewMainHelperWithOptions(nil)
	defer mainHelper.Close()

	sqlstore.InitTest(mainHelper.Logger)

	mainHelper.Main(m)
	sqlstore.TearDownTest()
}
