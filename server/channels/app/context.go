// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/sqlstore"
	"git.biggo.com/Funmula/BigGoChat/server/public/plugin"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

// RequestContextWithMaster adds the context value that master DB should be selected for this request.
func RequestContextWithMaster(c request.CTX) request.CTX {
	return sqlstore.RequestContextWithMaster(c)
}

func pluginContext(c request.CTX) *plugin.Context {
	context := &plugin.Context{
		RequestId:      c.RequestId(),
		SessionId:      c.Session().Id,
		IPAddress:      c.IPAddress(),
		AcceptLanguage: c.AcceptLanguage(),
		UserAgent:      c.UserAgent(),
	}
	return context
}
