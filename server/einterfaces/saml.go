// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type SamlInterface interface {
	ConfigureSP(c request.CTX) error
	BuildRequest(c request.CTX, relayState string) (*model.SamlAuthRequest, *model.AppError)
	DoLogin(c request.CTX, encodedXML string, relayState map[string]string) (*model.User, *model.AppError)
	GetMetadata(c request.CTX) (string, *model.AppError)
	CheckProviderAttributes(c request.CTX, SS *model.SamlSettings, ouser *model.User, patch *model.UserPatch) string
}
