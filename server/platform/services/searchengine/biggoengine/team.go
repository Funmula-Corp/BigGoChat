package biggoengine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
)

func (be *BiggoEngine) SearchTeams(userId string, searchParams []*model.SearchParams, page, perPage int) (result []string, total int64, aErr *model.AppError) {
	const keyErrWhere string = "biggo.search_engine.search_teams"
	const endpoint string = "/api/v1/search/team"
	var (
		err error
		req *http.Request
		res *http.Response

		searchUrl *url.URL
	)

	var buffer []byte
	parameters := map[string]any{
		"term": searchParams[0].Terms,
		"page": page,
		"size": perPage,
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

	results := struct {
		TeamIds []string `json:"ids"`
		Total   int64    `json:"total"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&results); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_response", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to parse search result", mlog.Err(err))
		return
	}

	result = results.TeamIds
	total = results.Total
	return
}
