package biggoengine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
	"git.biggo.com/Funmula/BigGoChat/server/v8/platform/services/searchengine/biggoengine/cfg"
)

func (be *BiggoEngine) DeleteChannel(channel *model.Channel) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexChannel(rctx request.CTX, channel *model.Channel, userIDs, teamMemberIDs []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexChannelsBulk(channels []*model.Channel) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchChannels(teamId, userID, term string, isGuest bool, page, perPage int) (result []string, aErr *model.AppError) {
	mlog.Debug("SEARCH-ENGINE", mlog.String("type", "channel"))
	const keyErrWhere string = "biggo.search_engine.search_channels"
	const endpoint string = "/api/v1/search/channel"
	var (
		err error
		req *http.Request
		res *http.Response

		searchUrl *url.URL
	)

	var buffer []byte
	parameters := map[string]any{
		"term":     term,
		"team_id":  teamId,
		"user_id":  userID,
		"is_guest": isGuest,
		"page":     page,
		"size":     perPage,
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

	result = []string{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		aErr = model.NewAppError(keyErrWhere, "parse_response", nil, aErr.Error(), http.StatusInternalServerError)
		mlog.Error("failed to parse search result", mlog.Err(err))
		return
	}
	mlog.Info("SEARCH-RESULT", mlog.Any("result", result))
	return
}
