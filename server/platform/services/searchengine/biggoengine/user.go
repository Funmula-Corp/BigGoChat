package biggoengine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
)

func (be *BiggoEngine) DeleteUser(user *model.User) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexUser(rctx request.CTX, user *model.User, teamsIds, channelsIds []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexUsersBulk(users []*model.UserForIndexing) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchUsers(userId string, term string, page, perPage int) (result []string, aErr *model.AppError) {
	_, result, aErr = be.SearchUsersInChannel(userId, "", term, page, perPage)
	return
}

func (be *BiggoEngine) SearchUsersInChannel(userId string, channelId string, term string, page, perPage int) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	const keyErrWhere string = "biggo.search_engine.search_users"
	const endpoint string = "/api/v1/search/user"
	var (
		err error
		req *http.Request
		res *http.Response

		searchUrl *url.URL
	)

	var buffer []byte
	parameters := map[string]any{
		"term":       term,
		"page":       page,
		"size":       perPage,
		"channel_id": channelId,
		"user_id":    userId,
	}
	if buffer, err = json.Marshal(&parameters); err != nil {
		aErr = model.NewAppError(keyErrWhere, "marshal_search_buffer", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to marshal search buffer", mlog.Err(err))
		return
	}

	if searchUrl, err = url.Parse(cfg.SearchServiceHost(be.config)); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_service_url", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to parse service url", mlog.Err(err))
		return
	}
	searchUrl = searchUrl.JoinPath(endpoint)

	if req, err = http.NewRequest(http.MethodPost, searchUrl.String(), bytes.NewReader(buffer)); err != nil {
		aErr = model.NewAppError(keyErrWhere, "create_request", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to create request", mlog.Err(err))
		return
	}

	if res, err = http.DefaultClient.Do(req); err != nil {
		aErr = model.NewAppError(keyErrWhere, "get_response", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to get response", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	results := map[string][]string{}
	if err = json.NewDecoder(res.Body).Decode(&results); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_response", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to parse search result", mlog.Err(err))
		return
	}

	found := false
	if userInChannel, found = results["in_channel"]; !found {
		userInChannel = []string{}
	}
	if userNotInChannel, found = results["not_in_channel"]; !found {
		userNotInChannel = []string{}
	}
	return
}
