package biggoengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

const (
	EsUserIndex string = "mm_biggoengine_user"
)

func (be *BiggoEngine) DeleteUser(user *model.User) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexUser(rctx request.CTX, user *model.User, teamsIds, channelsIds []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexUsersBulk(users []*model.UserForIndexing) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	if indexer, err = clients.EsBulkIndex(EsUserIndex); err != nil {
		aErr = model.NewAppError("BiggoIndexer.IndexUsersBulk", "engine.biggo.bulk_index.create.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	for _, value := range users {
		var buffer []byte
		if buffer, err = json.Marshal(value); err != nil {
			//aErr = model.NewAppError("BiggoIndexer.IndexUsersBulk", "engine.biggo.unmarshal.user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: value.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      EsUserIndex,
		}); err != nil {
			//aErr = model.NewAppError("BiggoIndexer.IndexUsersBulk", "engine.biggo.index.user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}
	}
	return
}

func (be *BiggoEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	var (
		cli *elasticsearch.Client
		err error
		res *esapi.Response
	)

	if cli, err = clients.EsClient(); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInChannel", "engine.biggo.client.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	query := `{
		"query": {
			"prefix": {
				"username": {
					"value": "%s"
				}
			}
		}
	}`

	if res, err = cli.Search(
		cli.Search.WithIndex(EsUserIndex),
		cli.Search.WithBody(strings.NewReader(fmt.Sprintf(query, term))),
	); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInChannel", "engine.biggo.search.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			aErr = model.NewAppError("BiggoIndexer.SearchUsersInChannel", "engine.biggo.response.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(errors.New(string(buffer)))
		}
		return
	}

	var response clients.EnvelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInChannel", "engine.biggo.parse.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	userNotInChannel = []string{}
	for _, hit := range response.Hits.Hits {
		userNotInChannel = append(userNotInChannel, hit.ID)
	}
	return
}

func (be *BiggoEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (result []string, aErr *model.AppError) {
	var (
		cli *elasticsearch.Client
		err error
		res *esapi.Response
	)

	if cli, err = clients.EsClient(); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInTeam", "engine.biggo.client.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	if term == "" {
		term = "*"
	}

	query := `{
		"query": {
			"prefix": {
				"username": {
					"value": "%s"
				}
			}
		}
	}`

	if res, err = cli.Search(
		cli.Search.WithIndex(EsUserIndex),
		cli.Search.WithBody(strings.NewReader(fmt.Sprintf(query, term))),
	); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInTeam", "engine.biggo.search.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			aErr = model.NewAppError("BiggoIndexer.SearchUsersInTeam", "engine.biggo.response.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(errors.New(string(buffer)))
		}
		return
	}

	var response clients.EnvelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		aErr = model.NewAppError("BiggoIndexer.SearchUsersInTeam", "engine.biggo.parse.elasticsearch.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		return
	}

	result = []string{}
	for _, hit := range response.Hits.Hits {
		result = append(result, hit.ID)
	}
	return
}
