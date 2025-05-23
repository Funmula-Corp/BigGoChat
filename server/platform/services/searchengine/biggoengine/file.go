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

func (be *BiggoEngine) DeleteFile(fileID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeleteFilesBatch(rctx request.CTX, endTime, limit int64) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeletePostFiles(rctx request.CTX, postID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeleteUserFiles(rctx request.CTX, userID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexFile(file *model.FileInfo, channelId string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexFilesBulk(files []*model.FileForIndexing) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchFiles(userId string, searchParams []*model.SearchParams, page, perPage int) (result []string, aErr *model.AppError) {
	mlog.Debug("SEARCH-ENGINE", mlog.String("type", "file"))
	const keyErrWhere string = "biggo.search_engine.search_files"
	const endpoint string = "/api/v1/search/file"
	var (
		err error
		req *http.Request
		res *http.Response

		searchUrl *url.URL
	)

	if len(searchParams) <= 0 {
		result = []string{}
		return
	}
	searchParam := searchParams[0]

	var buffer []byte
	parameters := map[string]any{
		"term":    searchParam.Terms,
		"user_id": userId,
		"page":    page,
		"size":    perPage,
		"params":  *searchParam,
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
