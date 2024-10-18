// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import (
	"fmt"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/utils/fileutils"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/i18n"
)

// TranslationsPreInit loads translations from filesystem if they are not
// loaded already and assigns english while loading server config.
func TranslationsPreInit() error {
	translationsDir := "i18n"

	i18nDirectory, found := fileutils.FindDirRelBinary(translationsDir)
	if !found {
		return fmt.Errorf("unable to find i18n directory at %q", translationsDir)
	}

	return i18n.TranslationsPreInit(i18nDirectory)
}
