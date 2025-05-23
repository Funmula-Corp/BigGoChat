// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/audit"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func (api *API) InitSystemLocal() {
	api.BaseRoutes.System.Handle("/ping", api.APILocal(getSystemPing)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/logs", api.APILocal(getLogs)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APILocal(setServerBusy)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APILocal(getServerBusyExpires)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APILocal(clearServerBusy)).Methods("DELETE")
	api.BaseRoutes.System.Handle("/support_packet", api.APILocal(generateSupportPacket)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/integrity", api.APILocal(localCheckIntegrity)).Methods("POST")
	api.BaseRoutes.System.Handle("/schema/version", api.APILocal(getAppliedSchemaMigrations)).Methods("GET")
}

func localCheckIntegrity(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("localCheckIntegrity", audit.Fail)
	defer c.LogAuditRec(auditRec)

	var results []model.IntegrityCheckResult
	resultsChan := c.App.CheckIntegrity()
	for result := range resultsChan {
		results = append(results, result)
	}

	data, err := json.Marshal(results)
	if err != nil {
		c.Err = model.NewAppError("Api4.localCheckIntegrity", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	auditRec.Success()
	w.Write(data)
}
