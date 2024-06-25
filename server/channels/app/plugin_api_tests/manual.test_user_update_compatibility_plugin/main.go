// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type MyPlugin struct {
	plugin.MattermostPlugin
}


func (p *MyPlugin) MessageWillBePosted(_ *plugin.Context, post *model.Post) (*model.Post, string) {
	var user *model.User
	user, err := p.API.GetUser(post.UserId)
	if err != nil {
		return nil, fmt.Sprintf("Failed to getuser %v", err)
	}

	user.Nickname = "updated"
	_, err = p.API.UpdateUser(user)
	if err != nil {
		return nil, fmt.Sprintf("Failed to UpdateUser %v", err)
	}
	return nil, "OK"
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
