// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/audit"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
)

func (api *API) InitLicenseLocal() {
	api.BaseRoutes.APIRoot.Handle("/license", api.APILocal(localAddLicense, handlerParamFileAPI)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/license", api.APILocal(localRemoveLicense)).Methods("DELETE")
}

func localAddLicense(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("localAddLicense", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempt")

	err := r.ParseMultipartForm(*c.App.Config().FileSettings.MaxFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := r.MultipartForm

	fileArray, ok := m.File["license"]
	if !ok {
		c.Err = model.NewAppError("addLicense", "api.license.add_license.no_file.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if len(fileArray) <= 0 {
		c.Err = model.NewAppError("addLicense", "api.license.add_license.array.app_error", nil, "", http.StatusBadRequest)
		return
	}

	fileData := fileArray[0]
	audit.AddEventParameter(auditRec, "filename", fileData.Filename)

	file, err := fileData.Open()
	if err != nil {
		c.Err = model.NewAppError("addLicense", "api.license.add_license.open.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)

	license, appErr := c.App.Srv().SaveLicense(buf.Bytes())
	if appErr != nil {
		if appErr.Id == model.ExpiredLicenseError {
			c.LogAudit("failed - expired or non-started license")
		} else if appErr.Id == model.InvalidLicenseError {
			c.LogAudit("failed - invalid license")
		} else {
			c.LogAudit("failed - unable to save license")
		}
		c.Err = appErr
		return
	}

	auditRec.Success()
	c.LogAudit("success")

	if err := json.NewEncoder(w).Encode(license); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func localRemoveLicense(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("localRemoveLicense", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempt")

	if err := c.App.Srv().RemoveLicense(); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("success")

	ReturnStatusOK(w)
}
