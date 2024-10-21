// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type AccountMigrationInterface interface {
	MigrateToLdap(c request.CTX, fromAuthService string, foreignUserFieldNameToMatch string, force bool, dryRun bool) *model.AppError
	MigrateToSaml(c request.CTX, fromAuthService string, usersMap map[string]string, auto bool, dryRun bool) *model.AppError
}
