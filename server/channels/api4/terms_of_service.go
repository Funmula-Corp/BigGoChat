// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/audit"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
)

func (api *API) InitTermsOfService() {
	api.BaseRoutes.TermsOfService.Handle("", api.APISessionRequired(getLatestTermsOfService)).Methods("GET")
	api.BaseRoutes.TermsOfService.Handle("", api.APISessionRequired(createTermsOfService)).Methods("POST")
}

func getLatestTermsOfService(c *Context, w http.ResponseWriter, r *http.Request) {
	termsOfService, err := c.App.GetLatestTermsOfService()
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(termsOfService); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func createTermsOfService(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	if license := c.App.Channels().License(); license == nil || !*license.Features.CustomTermsOfService {
		c.Err = model.NewAppError("createTermsOfService", "api.create_terms_of_service.custom_terms_of_service_disabled.app_error", nil, "", http.StatusBadRequest)
		return
	}

	auditRec := c.MakeAuditRecord("createTermsOfService", audit.Fail)
	defer c.LogAuditRec(auditRec)

	props := model.MapFromJSON(r.Body)
	text := props["text"]
	userId := c.AppContext.Session().UserId

	if text == "" {
		c.Err = model.NewAppError("Config.IsValid", "api.create_terms_of_service.empty_text.app_error", nil, "", http.StatusBadRequest)
		return
	}

	oldTermsOfService, err := c.App.GetLatestTermsOfService()
	if err != nil && err.Id != app.ErrorTermsOfServiceNoRowsFound {
		c.Err = err
		return
	}

	if oldTermsOfService == nil || oldTermsOfService.Text != text {
		termsOfService, err := c.App.CreateTermsOfService(text, userId)
		if err != nil {
			c.Err = err
			return
		}

		if err := json.NewEncoder(w).Encode(termsOfService); err != nil {
			c.Logger.Warn("Error while writing response", mlog.Err(err))
		}
	} else {
		if err := json.NewEncoder(w).Encode(oldTermsOfService); err != nil {
			c.Logger.Warn("Error while writing response", mlog.Err(err))
		}
	}
	auditRec.Success()
}
