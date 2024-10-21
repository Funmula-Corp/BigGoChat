// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import "git.biggo.com/Funmula/BigGoChat/server/public/model"

type IPFilteringInterface interface {
	ApplyIPFilters(allowedIPRanges *model.AllowedIPRanges) (*model.AllowedIPRanges, error)
	GetIPFilters() (*model.AllowedIPRanges, error)
}
