// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"os"

	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mattermost/commands"
	// Import and register app layer slash commands
	_ "git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/slashcommands"
	// Plugins
	_ "git.biggo.com/Funmula/BigGoChat/server/v8/biggo/oauth"
	_ "git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/oauthproviders/gitlab"

)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
