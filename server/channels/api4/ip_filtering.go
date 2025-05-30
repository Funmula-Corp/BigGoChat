// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/audit"
	"git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func (api *API) InitIPFiltering() {
	api.BaseRoutes.IPFiltering.Handle("", api.APISessionRequired(getIPFilters)).Methods("GET")
	api.BaseRoutes.IPFiltering.Handle("", api.APISessionRequired(applyIPFilters)).Methods("POST")
	api.BaseRoutes.IPFiltering.Handle("/my_ip", api.APISessionRequired(myIP)).Methods("GET")
}

func ensureIPFilteringInterface(c *Context, where string) (einterfaces.IPFilteringInterface, bool) {
	if c.App.IPFiltering() == nil || !c.App.Config().FeatureFlags.CloudIPFiltering || c.App.License() == nil || c.App.License().SkuShortName != model.LicenseShortSkuEnterprise {
		c.Err = model.NewAppError(where, "api.context.ip_filtering.not_available.app_error", nil, "", http.StatusNotImplemented)
		return nil, false
	}
	return c.App.IPFiltering(), true
}

func getIPFilters(c *Context, w http.ResponseWriter, r *http.Request) {
	ipFiltering, ok := ensureIPFilteringInterface(c, "getIPFilters")
	if !ok {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadIPFilters) {
		c.SetPermissionError(model.PermissionSysconsoleReadIPFilters)
		return
	}

	allowedRanges, err := ipFiltering.GetIPFilters()
	if err != nil {
		c.Err = model.NewAppError("getIPFilters", "api.context.ip_filtering.get_ip_filters.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	if err := json.NewEncoder(w).Encode(allowedRanges); err != nil {
		c.Err = model.NewAppError("getIPFilters", "api.context.ip_filtering.get_ip_filters.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}
}

func applyIPFilters(c *Context, w http.ResponseWriter, r *http.Request) {
	ipFiltering, ok := ensureIPFilteringInterface(c, "applyIPFilters")
	if !ok {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleWriteIPFilters) {
		c.SetPermissionError(model.PermissionSysconsoleWriteIPFilters)
		return
	}

	auditRec := c.MakeAuditRecord("applyIPFilters", audit.Fail)
	defer c.LogAuditRecWithLevel(auditRec, app.LevelContent)

	allowedRanges := &model.AllowedIPRanges{} // Initialize the allowedRanges variable
	if err := json.NewDecoder(r.Body).Decode(allowedRanges); err != nil {
		c.Err = model.NewAppError("applyIPFilters", "api.context.ip_filtering.apply_ip_filters.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	audit.AddEventParameterAuditable(auditRec, "IPFilter", allowedRanges)

	updatedAllowedRanges, err := ipFiltering.ApplyIPFilters(allowedRanges)

	if err != nil {
		c.Err = model.NewAppError("applyIPFilters", "api.context.ip_filtering.apply_ip_filters.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	auditRec.Success()

	go c.App.SendIPFiltersChangedEmail(c.AppContext, c.AppContext.Session().UserId)

	if err := json.NewEncoder(w).Encode(updatedAllowedRanges); err != nil {
		c.Err = model.NewAppError("getIPFilters", "api.context.ip_filtering.get_ip_filters.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}
}

func myIP(c *Context, w http.ResponseWriter, r *http.Request) {
	_, ok := ensureIPFilteringInterface(c, "myIP")

	if !ok {
		return
	}

	response := &model.GetIPAddressResponse{
		IP: c.AppContext.IPAddress(),
	}

	json, err := json.Marshal(response)
	if err != nil {
		c.Err = model.NewAppError("myIP", "api.context.ip_filtering.get_my_ip.failed", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	w.Write(json)
}
