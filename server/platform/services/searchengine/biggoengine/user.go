package biggoengine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
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

func (be *BiggoEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (result []string, aErr *model.AppError) {
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
		"term": term,
		"page": 0,
		"size": options.Limit,
	}
	if buffer, err = json.Marshal(&parameters); err != nil {
		aErr = model.NewAppError(keyErrWhere, "marshal_search_buffer", nil, aErr.Error(), http.StatusInternalServerError)
		return
	}

	if searchUrl, err = url.Parse(cfg.SearchServiceHost(be.config)); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_service_url", nil, aErr.Error(), http.StatusInternalServerError)
		return
	}
	searchUrl = searchUrl.JoinPath(endpoint)

	if req, err = http.NewRequest(http.MethodPost, searchUrl.String(), bytes.NewReader(buffer)); err != nil {
		aErr = model.NewAppError(keyErrWhere, "create_request", nil, aErr.Error(), http.StatusInternalServerError)
		return
	}

	if res, err = http.DefaultClient.Do(req); err != nil {
		aErr = model.NewAppError(keyErrWhere, "get_response", nil, aErr.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	result = []string{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_response", nil, aErr.Error(), http.StatusInternalServerError)
		return
	}
	return
}
