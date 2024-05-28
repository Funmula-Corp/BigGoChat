// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package indexer

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/storetest"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/utils/testutils"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/bleveengine"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

func TestBleveIndexer(t *testing.T) {
	mockStore := &storetest.Store{}
	defer mockStore.AssertExpectations(t)

	t.Run("Call GetOldestEntityCreationTime for the first indexing call", func(t *testing.T) {
		job := &model.Job{
			Id:       model.NewId(),
			CreateAt: model.GetMillis(),
			Status:   model.JobStatusPending,
			Type:     model.JobTypeBlevePostIndexing,
		}

		mockStore.JobStore.On("UpdateStatusOptimistically", job.Id, model.JobStatusPending, model.JobStatusInProgress).Return(true, nil)
		mockStore.JobStore.On("UpdateOptimistically", job, model.JobStatusInProgress).Return(true, nil)
		mockStore.PostStore.On("GetOldestEntityCreationTime").Return(int64(1), errors.New("")) // intentionally return error to return from function

		tempDir, err := os.MkdirTemp("", "setupConfigFile")
		require.NoError(t, err)

		t.Cleanup(func() {
			os.RemoveAll(tempDir)
		})

		cfg := &model.Config{
			BleveSettings: model.BleveSettings{
				EnableIndexing: model.NewBool(true),
				IndexDir:       model.NewString(tempDir),
			},
		}

		jobServer := &jobs.JobServer{
			Store: mockStore,
			ConfigService: &testutils.StaticConfigService{
				Cfg: cfg,
			},
		}

		bleveEngine := bleveengine.NewBleveEngine(cfg)
		aErr := bleveEngine.Start()
		require.Nil(t, aErr)

		worker := &BleveIndexerWorker{
			jobServer: jobServer,
			engine:    bleveEngine,
			logger:    mlog.CreateConsoleTestLogger(t),
		}

		worker.DoJob(job)
	})
}
