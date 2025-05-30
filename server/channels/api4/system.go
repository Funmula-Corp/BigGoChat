// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/audit"
	"git.biggo.com/Funmula/BigGoChat/server/v8/platform/services/cache"
	"git.biggo.com/Funmula/BigGoChat/server/v8/platform/shared/web"
)

const (
	RedirectLocationCacheSize     = 10000
	RedirectLocationMaximumLength = 2100
	RedirectLocationCacheExpiry   = 1 * time.Hour
	DefaultServerBusySeconds      = 3600
	MaxServerBusySeconds          = 86400
)

var redirectLocationDataCache = cache.NewLRU(cache.LRUOptions{
	Size: RedirectLocationCacheSize,
})

func (api *API) InitSystem() {
	api.BaseRoutes.System.Handle("/ping", api.APIHandler(getSystemPing)).Methods("GET")

	api.BaseRoutes.System.Handle("/timezones", api.APISessionRequired(getSupportedTimezones)).Methods("GET")

	api.BaseRoutes.APIRoot.Handle("/audits", api.APISessionRequired(getAudits)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/email/test", api.APISessionRequired(testEmail)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/site_url/test", api.APISessionRequired(testSiteURL)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/file/s3_test", api.APISessionRequired(testS3)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/database/recycle", api.APISessionRequired(databaseRecycle)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/caches/invalidate", api.APISessionRequired(invalidateCaches)).Methods("POST")

	api.BaseRoutes.APIRoot.Handle("/logs", api.APISessionRequired(getLogs)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/logs/download", api.APISessionRequired(downloadLogs)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/logs/query", api.APISessionRequired(queryLogs)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/logs", api.APIHandler(postLog)).Methods("POST")

	api.BaseRoutes.APIRoot.Handle("/analytics/old", api.APISessionRequired(getAnalytics)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/latest_version", api.APISessionRequired(getLatestVersion)).Methods("GET")

	api.BaseRoutes.APIRoot.Handle("/redirect_location", api.APISessionRequiredTrustRequester(getRedirectLocation)).Methods("GET")

	api.BaseRoutes.APIRoot.Handle("/notifications/ack", api.APISessionRequired(pushNotificationAck)).Methods("POST")

	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APISessionRequired(setServerBusy)).Methods("POST")
	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APISessionRequired(getServerBusyExpires)).Methods("GET")
	api.BaseRoutes.APIRoot.Handle("/server_busy", api.APISessionRequired(clearServerBusy)).Methods("DELETE")
	api.BaseRoutes.APIRoot.Handle("/restart", api.APISessionRequired(restart)).Methods("POST")
	api.BaseRoutes.System.Handle("/notices/{team_id:[A-Za-z0-9]+}", api.APISessionRequired(getProductNotices)).Methods("GET")
	api.BaseRoutes.System.Handle("/notices/view", api.APISessionRequired(updateViewedProductNotices)).Methods("PUT")
	api.BaseRoutes.System.Handle("/support_packet", api.APISessionRequired(generateSupportPacket)).Methods("GET")
	api.BaseRoutes.System.Handle("/onboarding/complete", api.APISessionRequired(getOnboarding)).Methods("GET")
	api.BaseRoutes.System.Handle("/onboarding/complete", api.APISessionRequired(completeOnboarding)).Methods("POST")
	api.BaseRoutes.System.Handle("/schema/version", api.APISessionRequired(getAppliedSchemaMigrations)).Methods("GET")
}

func generateSupportPacket(c *Context, w http.ResponseWriter, r *http.Request) {
	const FileMime = "application/zip"
	const OutputDirectory = "support_packet"

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("generateSupportPacket", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	// Support packet generation is limited to system admins (MM-42271).
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	// We support the existing API hence the logs are always included
	// if nothing specified.
	includeLogs := true
	if r.FormValue("basic_server_logs") == "false" {
		includeLogs = false
	}
	supportPacketOptions := &model.SupportPacketOptions{
		IncludeLogs:   includeLogs,
		PluginPackets: r.Form["plugin_packets"],
	}

	// Checking to see if the server has a e10 or e20 license (this feature is only permitted for servers with licenses)
	if c.App.Channels().License() == nil {
		c.Err = model.NewAppError("Api4.generateSupportPacket", "api.no_license", nil, "", http.StatusForbidden)
		return
	}

	fileDatas := c.App.GenerateSupportPacket(c.AppContext, supportPacketOptions)

	// Constructing the ZIP file name as per spec (mattermost_support_packet_YYYY-MM-DD-HH-MM.zip)
	// Note that this filename is also being checked at the webapp, please update the
	// regex within the commercial_support_modal.tsx file if the naming convention ever changes.
	now := time.Now()
	outputZipFilename := fmt.Sprintf("mattermost_support_packet_%s.zip", now.Format("2006-01-02-03-04"))

	fileStorageBackend := c.App.FileBackend()

	// We do this incase we get concurrent requests, we will always have a unique directory.
	// This is to avoid the situation where we try to write to the same directory while we are trying to delete it (further down)
	outputDirectoryToUse := OutputDirectory + "_" + model.NewId()
	err := c.App.CreateZipFileAndAddFiles(fileStorageBackend, fileDatas, outputZipFilename, outputDirectoryToUse)
	if err != nil {
		c.Err = model.NewAppError("Api4.generateSupportPacket", "api.unable_to_create_zip_file", nil, "", http.StatusForbidden).Wrap(err)
		return
	}

	fileBytes, err := fileStorageBackend.ReadFile(path.Join(outputDirectoryToUse, outputZipFilename))
	defer fileStorageBackend.RemoveDirectory(outputDirectoryToUse)
	if err != nil {
		c.Err = model.NewAppError("Api4.generateSupportPacket", "api.unable_to_read_file_from_backend", nil, "", http.StatusForbidden).Wrap(err)
		return
	}
	fileBytesReader := bytes.NewReader(fileBytes)

	// Send the zip file back to client
	// We are able to pass 0 for content size due to the fact that Golang's serveContent (https://golang.org/src/net/http/fs.go)
	// already sets that for us
	web.WriteFileResponse(outputZipFilename, FileMime, 0, now, *c.App.Config().ServiceSettings.WebserverMode, fileBytesReader, true, w, r)
}

func getSystemPing(c *Context, w http.ResponseWriter, r *http.Request) {
	reqs := c.App.Config().ClientRequirements

	s := make(map[string]string)
	s[model.STATUS] = model.StatusOk
	s["AndroidLatestVersion"] = reqs.AndroidLatestVersion
	s["AndroidMinVersion"] = reqs.AndroidMinVersion
	s["IosLatestVersion"] = reqs.IosLatestVersion
	s["IosMinVersion"] = reqs.IosMinVersion

	testflag := c.App.Config().FeatureFlags.TestFeature
	if testflag != "off" {
		s["TestFeatureFlag"] = testflag
	}

	actualGoroutines := runtime.NumGoroutine()
	if *c.App.Config().ServiceSettings.GoroutineHealthThreshold > 0 && actualGoroutines >= *c.App.Config().ServiceSettings.GoroutineHealthThreshold {
		c.Logger.Warn("The number of running goroutines is over the health threshold", mlog.Int("goroutines", actualGoroutines), mlog.Int("health_threshold", *c.App.Config().ServiceSettings.GoroutineHealthThreshold))
		s[model.STATUS] = model.StatusUnhealthy
	}

	// Enhanced ping health check:
	// If an extra form value is provided then perform extra health checks for
	// database and file storage backends.
	if r.FormValue("get_server_status") == "true" {
		dbStatusKey := "database_status"
		s[dbStatusKey] = model.StatusOk

		writeErr := c.App.DBHealthCheckWrite()
		if writeErr != nil {
			c.Logger.Warn("Unable to write to database.", mlog.Err(writeErr))
			s[dbStatusKey] = model.StatusUnhealthy
			s[model.STATUS] = model.StatusUnhealthy
		}

		writeErr = c.App.DBHealthCheckDelete()
		if writeErr != nil {
			c.Logger.Warn("Unable to remove ping health check value from database.", mlog.Err(writeErr))
			s[dbStatusKey] = model.StatusUnhealthy
			s[model.STATUS] = model.StatusUnhealthy
		}

		if s[dbStatusKey] == model.StatusOk {
			c.Logger.Debug("Able to write to database.")
		}

		filestoreStatusKey := "filestore_status"
		s[filestoreStatusKey] = model.StatusOk
		appErr := c.App.TestFileStoreConnection()
		if appErr != nil {
			c.Logger.Warn("Unable to test filestore connection.", mlog.Err(appErr))
			s[filestoreStatusKey] = model.StatusUnhealthy
			s[model.STATUS] = model.StatusUnhealthy
		}

		w.Header().Set(model.STATUS, s[model.STATUS])
		w.Header().Set(dbStatusKey, s[dbStatusKey])
		w.Header().Set(filestoreStatusKey, s[filestoreStatusKey])
	}

	if deviceID := r.FormValue("device_id"); deviceID != "" {
		s["CanReceiveNotifications"] = c.App.SendTestPushNotification(deviceID)
	}

	s["ActiveSearchBackend"] = c.App.ActiveSearchBackend()

	if s[model.STATUS] != model.StatusOk && r.FormValue("use_rest_semantics") != "true" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(model.MapToJSON(s)))
}

func testEmail(c *Context, w http.ResponseWriter, r *http.Request) {
	var cfg *model.Config
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		c.Logger.Warn("Error decoding the config", mlog.Err(err))
	}
	if cfg == nil {
		cfg = c.App.Config()
	}

	if checkHasNilFields(&cfg.EmailSettings) {
		c.Err = model.NewAppError("testEmail", "api.file.test_connection_email_settings_nil.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionTestEmail) {
		c.SetPermissionError(model.PermissionTestEmail)
		return
	}

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("testEmail", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	appErr := c.App.TestEmail(c.AppContext, c.AppContext.Session().UserId, cfg)
	if appErr != nil {
		c.Err = appErr
		return
	}

	ReturnStatusOK(w)
}

func testSiteURL(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionTestSiteURL) {
		c.SetPermissionError(model.PermissionTestSiteURL)
		return
	}

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("testSiteURL", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	props := model.MapFromJSON(r.Body)
	siteURL := props["site_url"]
	if siteURL == "" {
		c.SetInvalidParam("site_url")
		return
	}

	appErr := c.App.TestSiteURL(c.AppContext, siteURL)
	if appErr != nil {
		c.Err = appErr
		return
	}

	ReturnStatusOK(w)
}

func getAudits(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("getAudits", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionReadAudits) {
		c.SetPermissionError(model.PermissionReadAudits)
		return
	}

	audits, appErr := c.App.GetAuditsPage(c.AppContext, "", c.Params.Page, c.Params.PerPage)
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()
	audit.AddEventParameter(auditRec, "page", c.Params.Page)
	audit.AddEventParameter(auditRec, "audits_per_page", c.Params.LogsPerPage)

	if err := json.NewEncoder(w).Encode(audits); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func databaseRecycle(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionRecycleDatabaseConnections) {
		c.SetPermissionError(model.PermissionRecycleDatabaseConnections)
		return
	}

	auditRec := c.MakeAuditRecord("databaseRecycle", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("databaseRecycle", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	c.App.RecycleDatabaseConnection(c.AppContext)

	auditRec.Success()
	ReturnStatusOK(w)
}

func invalidateCaches(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionInvalidateCaches) {
		c.SetPermissionError(model.PermissionInvalidateCaches)
		return
	}

	auditRec := c.MakeAuditRecord("invalidateCaches", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("invalidateCaches", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	appErr := c.App.Srv().InvalidateAllCaches()
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ReturnStatusOK(w)
}

func queryLogs(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("queryLogs", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("queryLogs", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionGetLogs) {
		c.SetPermissionError(model.PermissionGetLogs)
		return
	}

	var logFilter *model.LogFilter
	err := json.NewDecoder(r.Body).Decode(&logFilter)
	if err != nil || logFilter == nil {
		c.Err = model.NewAppError("queryLogs", "api.system.logs.invalidFilter", nil, "", http.StatusInternalServerError)
		return
	}

	logs, logerr := c.App.QueryLogs(c.AppContext, c.Params.Page, c.Params.LogsPerPage, logFilter)
	if logerr != nil {
		c.Err = logerr
		return
	}

	logsJSON := make(map[string][]interface{})
	var result interface{}
	for node, logLines := range logs {
		for _, log := range logLines {
			err2 := json.Unmarshal([]byte(log), &result)
			if err2 == nil {
				logsJSON[node] = append(logsJSON[node], result)
			} else {
				c.Logger.Warn("Error parsing log line in Server Logs")
			}
		}
	}

	auditRec.AddMeta("page", c.Params.Page)
	auditRec.AddMeta("logs_per_page", c.Params.LogsPerPage)

	w.Write(model.ToJSON(logsJSON))
}

func getLogs(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("getLogs", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("getLogs", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionGetLogs) {
		c.SetPermissionError(model.PermissionGetLogs)
		return
	}

	lines, appErr := c.App.GetLogs(c.AppContext, c.Params.Page, c.Params.LogsPerPage)
	if appErr != nil {
		c.Err = appErr
		return
	}

	audit.AddEventParameter(auditRec, "page", c.Params.Page)
	audit.AddEventParameter(auditRec, "logs_per_page", c.Params.LogsPerPage)

	w.Write([]byte(model.ArrayToJSON(lines)))
}

func downloadLogs(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("downloadLogs", audit.Fail)
	defer c.LogAuditRec(auditRec)
	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("downloadLogs", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionGetLogs) {
		c.SetPermissionError(model.PermissionGetLogs)
		return
	}

	fileData, err := c.App.GetMattermostLog(c.AppContext)
	if err != nil {
		c.Err = model.NewAppError("downloadLogs", "api.system.logs.download_bytes_buffer.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(fileData.Body)
	web.WriteFileResponse("mattermost.log",
		"text/plain",
		int64(len(fileData.Body)),
		time.Now(),
		*c.App.Config().ServiceSettings.WebserverMode,
		reader,
		true,
		w,
		r)
}

func postLog(c *Context, w http.ResponseWriter, r *http.Request) {
	forceToDebug := false

	if !*c.App.Config().ServiceSettings.EnableDeveloper {
		if c.AppContext.Session().UserId == "" {
			c.Err = model.NewAppError("postLog", "api.context.permissions.app_error", nil, "", http.StatusForbidden)
			return
		}

		if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
			forceToDebug = true
		}
	}

	var m map[string]string
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		c.Logger.Warn("Error decoding request.", mlog.Err(err))
	}
	if m == nil {
		m = map[string]string{}
	}

	lvl := m["level"]
	msg := m["message"]

	if len(msg) > 400 {
		msg = msg[0:399]
	}

	msg = "Client Logs API Endpoint Message: " + msg
	fields := []mlog.Field{
		mlog.String("type", "client_message"),
		mlog.String("user_agent", c.AppContext.UserAgent()),
	}

	if !forceToDebug && lvl == "ERROR" {
		mlog.Error(msg, fields...)
	} else {
		mlog.Debug(msg, fields...)
	}

	m["message"] = msg
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		c.Logger.Warn("Error while writing response.", mlog.Err(err))
	}
}

func getAnalytics(c *Context, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	teamId := r.URL.Query().Get("team_id")

	if name == "" {
		name = "standard"
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionGetAnalytics) {
		c.SetPermissionError(model.PermissionGetAnalytics)
		return
	}

	rows, appErr := c.App.GetAnalytics(c.AppContext, name, teamId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if rows == nil {
		c.SetInvalidParam("name")
		return
	}

	if err := json.NewEncoder(w).Encode(rows); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func getLatestVersion(c *Context, w http.ResponseWriter, r *http.Request) {
	if *c.App.Config().ExperimentalSettings.RestrictSystemAdmin {
		c.Err = model.NewAppError("latestVersion", "api.restricted_system_admin", nil, "", http.StatusForbidden)
		return
	}

	resp, appErr := c.App.GetLatestVersion(c.AppContext, "https://api.github.com/repos/mattermost/mattermost-server/releases/latest")
	if appErr != nil {
		c.Err = appErr
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		c.Logger.Warn("Unable to marshal JSON for latest version.", mlog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}

func getSupportedTimezones(c *Context, w http.ResponseWriter, r *http.Request) {
	supportedTimezones := c.App.Timezones().GetSupported()
	if supportedTimezones == nil {
		supportedTimezones = make([]string, 0)
	}

	b, err := json.Marshal(supportedTimezones)
	if err != nil {
		c.Logger.Warn("Unable to marshal JSON in timezones.", mlog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}

func testS3(c *Context, w http.ResponseWriter, r *http.Request) {
	var cfg *model.Config
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		c.Logger.Warn("Error decoding the config", mlog.Err(err))
	}
	if cfg == nil {
		cfg = c.App.Config()
	}

	if checkHasNilFields(&cfg.FileSettings) {
		c.Err = model.NewAppError("testS3", "api.file.test_connection_s3_settings_nil.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionTestS3) {
		c.SetPermissionError(model.PermissionTestS3)
		return
	}

	appErr := c.App.CheckMandatoryS3Fields(&cfg.FileSettings)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if *cfg.FileSettings.AmazonS3SecretAccessKey == model.FakeSetting {
		cfg.FileSettings.AmazonS3SecretAccessKey = c.App.Config().FileSettings.AmazonS3SecretAccessKey
	}

	appErr = c.App.TestFileStoreConnectionWithConfig(&cfg.FileSettings)
	if appErr != nil {
		c.Err = appErr
		return
	}

	ReturnStatusOK(w)
}

func getRedirectLocation(c *Context, w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["location"] = ""

	if !*c.App.Config().ServiceSettings.EnableLinkPreviews {
		w.Write([]byte(model.MapToJSON(m)))
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" {
		c.SetInvalidParam("url")
		return
	}

	var location string
	if err := redirectLocationDataCache.Get(url, &location); err == nil {
		m["location"] = location
		w.Write([]byte(model.MapToJSON(m)))
		return
	}

	client := c.App.HTTPService().MakeClient(false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Head(url)
	if err != nil {
		// Cache failures to prevent retries.
		redirectLocationDataCache.SetWithExpiry(url, "", RedirectLocationCacheExpiry)
		// Always return a success status and a JSON string to limit information returned to client.
		w.Write([]byte(model.MapToJSON(m)))
		return
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	location = res.Header.Get("Location")

	// If the location length is > 2100, we can probably ignore. Fixes https://mattermost.atlassian.net/browse/MM-54219
	if len(location) > RedirectLocationMaximumLength {
		// Treating as a "failure". Cache failures to prevent retries.
		redirectLocationDataCache.SetWithExpiry(url, "", RedirectLocationCacheExpiry)
		// Always return a success status and a JSON string to limit information returned to client.
		w.Write([]byte(model.MapToJSON(m)))
		return
	}

	redirectLocationDataCache.SetWithExpiry(url, location, RedirectLocationCacheExpiry)
	m["location"] = location

	w.Write([]byte(model.MapToJSON(m)))
}

func pushNotificationAck(c *Context, w http.ResponseWriter, r *http.Request) {
	var ack model.PushNotificationAck
	if jsonErr := json.NewDecoder(r.Body).Decode(&ack); jsonErr != nil {
		c.Err = model.NewAppError("pushNotificationAck",
			"api.push_notifications_ack.message.parse.app_error",
			nil,
			"",
			http.StatusBadRequest,
		).Wrap(jsonErr)
		return
	}

	if !*c.App.Config().EmailSettings.SendPushNotifications {
		c.Err = model.NewAppError("pushNotificationAck", "api.push_notification.disabled.app_error", nil, "", http.StatusNotImplemented)
		return
	}

	if ack.NotificationType == model.PushTypeMessage {
		c.App.CountNotificationAck(model.NotificationTypePush, ack.ClientPlatform)
	}

	err := c.App.SendAckToPushProxy(&ack)
	if ack.IsIdLoaded {
		if err != nil {
			c.App.CountNotificationReason(model.NotificationStatusError, model.NotificationTypePush, model.NotificationReasonPushProxySendError, ack.ClientPlatform)
			c.App.NotificationsLog().Error("Notification ack not sent to push proxy",
				mlog.String("type", model.NotificationTypePush),
				mlog.String("status", model.NotificationStatusError),
				mlog.String("reason", model.NotificationReasonPushProxySendError),
				mlog.String("ack_id", ack.Id),
				mlog.String("push_type", ack.NotificationType),
				mlog.String("post_id", ack.PostId),
				mlog.String("ack_type", ack.NotificationType),
				mlog.String("device_type", ack.ClientPlatform),
				mlog.Int("received_at", ack.ClientReceivedAt),
				mlog.Err(err),
			)
		} else {
			c.App.CountNotificationReason(model.NotificationStatusSuccess, model.NotificationTypePush, model.NotificationReason(""), ack.ClientPlatform)
		}
		// Return post data only when PostId is passed.
		if ack.PostId != "" && ack.NotificationType == model.PushTypeMessage {
			if _, appErr := c.App.GetPostIfAuthorized(c.AppContext, ack.PostId, c.AppContext.Session(), false); appErr != nil {
				c.Err = appErr
				return
			}

			notificationInterface := c.App.Notification()

			if notificationInterface == nil {
				c.Err = model.NewAppError("pushNotificationAck", "api.system.id_loaded.not_available.app_error", nil, "", http.StatusFound)
				return
			}

			msg, appError := notificationInterface.GetNotificationMessage(c.AppContext, &ack, c.AppContext.Session().UserId)
			if appError != nil {
				c.Err = model.NewAppError("pushNotificationAck", "api.push_notification.id_loaded.fetch.app_error", nil, "", http.StatusInternalServerError).Wrap(appError)
				return
			}
			if err2 := json.NewEncoder(w).Encode(msg); err2 != nil {
				c.Logger.Warn("Error while writing response", mlog.Err(err2))
			}
		}

		return
	} else if err != nil {
		c.Err = model.NewAppError("pushNotificationAck", "api.push_notifications_ack.forward.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	c.App.CountNotificationReason(model.NotificationStatusSuccess, model.NotificationTypePush, model.NotificationReason(""), ack.ClientPlatform)
	ReturnStatusOK(w)
}

func setServerBusy(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	// number of seconds to keep server marked busy
	secs := r.URL.Query().Get("seconds")
	if secs == "" {
		secs = strconv.FormatInt(DefaultServerBusySeconds, 10)
	}

	i, err := strconv.ParseInt(secs, 10, 64)
	if err != nil || i <= 0 || i > MaxServerBusySeconds {
		c.SetInvalidURLParam(fmt.Sprintf("seconds must be 1 - %d", MaxServerBusySeconds))
		return
	}

	auditRec := c.MakeAuditRecord("setServerBusy", audit.Fail)
	defer c.LogAuditRec(auditRec)
	audit.AddEventParameter(auditRec, "seconds", i)

	c.App.Srv().Platform().Busy.Set(time.Second * time.Duration(i))
	c.Logger.Warn("server busy state activated - non-critical services disabled", mlog.Int("seconds", i))

	auditRec.Success()
	ReturnStatusOK(w)
}

func clearServerBusy(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	auditRec := c.MakeAuditRecord("clearServerBusy", audit.Fail)
	defer c.LogAuditRec(auditRec)

	c.App.Srv().Platform().Busy.Clear()
	c.Logger.Info("server busy state cleared - non-critical services enabled")

	auditRec.Success()
	ReturnStatusOK(w)
}

func getServerBusyExpires(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	// We call to ToJSON because it actually returns a different struct
	// along with doing some computations.
	sbsJSON, jsonErr := c.App.Srv().Platform().Busy.ToJSON()
	if jsonErr != nil {
		c.Logger.Warn(jsonErr.Error())
	}

	if _, err := w.Write(sbsJSON); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func restart(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("restartServer", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	auditRec.Success()
	ReturnStatusOK(w)
	time.Sleep(1 * time.Second)

	go func() {
		c.App.Srv().Restart()
	}()
}

func getProductNotices(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	client, parseError := model.NoticeClientTypeFromString(r.URL.Query().Get("client"))
	if parseError != nil {
		c.SetInvalidParamWithErr("client", parseError)
		return
	}
	clientVersion := r.URL.Query().Get("clientVersion")
	locale := r.URL.Query().Get("locale")

	notices, appErr := c.App.GetProductNotices(c.AppContext, c.AppContext.Session().UserId, c.Params.TeamId, client, clientVersion, locale)
	if appErr != nil {
		c.Err = appErr
		return
	}
	result, _ := notices.Marshal()
	_, _ = w.Write(result)
}

func updateViewedProductNotices(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("updateViewedProductNotices", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempt")

	ids, err := model.SortedArrayFromJSON(r.Body)
	if err != nil {
		c.Err = model.NewAppError("updateViewedProductNotices", model.PayloadParseError, nil, "", http.StatusBadRequest).Wrap(err)
		return
	}
	appErr := c.App.UpdateViewedProductNotices(c.AppContext.Session().UserId, ids)
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()
	ReturnStatusOK(w)
}

func getOnboarding(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("getOnboarding", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempt")

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	firstAdminCompleteSetupObj, err := c.App.GetOnboarding()

	if err != nil {
		c.Err = model.NewAppError("getOnboarding", "app.system.get_onboarding_request.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	auditRec.Success()
	if err := json.NewEncoder(w).Encode(firstAdminCompleteSetupObj); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func completeOnboarding(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.Err = model.NewAppError("completeOnboarding", "app.system.complete_onboarding_request.no_first_user", nil, "", http.StatusForbidden)
		return
	}

	auditRec := c.MakeAuditRecord("completeOnboarding", audit.Fail)
	defer c.LogAuditRec(auditRec)

	onboardingRequest, err := model.CompleteOnboardingRequestFromReader(r.Body)
	if err != nil {
		c.Err = model.NewAppError("completeOnboarding", "app.system.complete_onboarding_request.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		return
	}
	audit.AddEventParameter(auditRec, "install_plugin", onboardingRequest.InstallPlugins)
	audit.AddEventParameterAuditable(auditRec, "onboarding_request", onboardingRequest)

	appErr := c.App.CompleteOnboarding(c.AppContext, onboardingRequest)
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()
	ReturnStatusOK(w)
}

func getAppliedSchemaMigrations(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionToAny(*c.AppContext.Session(), model.SysconsoleReadPermissions) {
		c.SetPermissionError(model.SysconsoleReadPermissions...)
		return
	}

	auditRec := c.MakeAuditRecord("getAppliedSchemaMigrations", audit.Fail)
	defer c.LogAuditRec(auditRec)

	migrations, appErr := c.App.GetAppliedSchemaMigrations()
	if appErr != nil {
		c.Err = appErr
		return
	}

	js, err := json.Marshal(migrations)
	if err != nil {
		c.Err = model.NewAppError("getAppliedMigrations", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	w.Write(js)
	auditRec.Success()
}

// returns true if the data has nil fields
// this is being used for testS3 and testEmail methods
func checkHasNilFields(value any) bool {
	v := reflect.Indirect(reflect.ValueOf(value))
	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			return true
		}
	}

	return false
}
