// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package migrations

import (
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

const (
	MigrationStateUnscheduled = "unscheduled"
	MigrationStateInProgress  = "in_progress"
	MigrationStateCompleted   = "completed"

	JobDataKeyMigration         = "migration_key"
	JobDataKeyMigrationLastDone = "last_done"
)

func MakeMigrationsList() []string {
	return []string{
		model.MigrationKeyAdvancedPermissionsPhase2,
	}
}

func GetMigrationState(c request.CTX, migration string, store store.Store) (string, *model.Job, *model.AppError) {
	if _, err := store.System().GetByName(migration); err == nil {
		return MigrationStateCompleted, nil, nil
	}

	jobs, err := store.Job().GetAllByType(c, model.JobTypeMigrations)
	if err != nil {
		return "", nil, model.NewAppError("GetMigrationState", "app.job.get_all.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, job := range jobs {
		if key, ok := job.Data[JobDataKeyMigration]; ok {
			if key != migration {
				continue
			}

			switch job.Status {
			case model.JobStatusInProgress, model.JobStatusPending:
				return MigrationStateInProgress, job, nil
			default:
				return MigrationStateUnscheduled, job, nil
			}
		}
	}

	return MigrationStateUnscheduled, nil, nil
}
