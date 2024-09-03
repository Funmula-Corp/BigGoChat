package helper

import (
	"fmt"
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
)

const (
	// "_source" filtering to possibly reduce the communicaiton overhead between ES -> Neo4j
	cypherESQueryStart string = `
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			_source: [''],
			query: {
				bool: {
					should: [`
	cypherESQueryEnd string = `
					]
				}
			}, from: $from, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit`
)

// this function composes the elasticsearch section of cypher queries
func ComposeSearchParamsQuery(config *model.Config, index string, page int, perSize int, termsField string, searchParam *model.SearchParams) (query string, queryParams map[string]any) {
	builder := strings.Builder{}
	builder.WriteString(cypherESQueryStart)
	builder.WriteString(composeSearchQuery(&termsField, searchParam))
	builder.WriteString(cypherESQueryEnd)
	query = builder.String()

	queryParams = map[string]any{
		"es_address": fmt.Sprintf("%s://%s:%.0f",
			cfg.ElasticsearchProtocol(config),
			cfg.ElasticsearchHost(config),
			cfg.ElasticsearchPort(config),
		),
		"es_index": index,
		"size":     perSize,
		"from":     perSize * page,
	}

	if searchParam.Terms != "" || searchParam.ExcludedTerms != "" {
		queryParams["term"] = searchParam.Terms
	}

	if len(searchParam.InChannels) > 0 {
		queryParams["in_channels"] = searchParam.InChannels
	}

	if len(searchParam.ExcludedChannels) > 0 {
		queryParams["exclude_channels"] = searchParam.ExcludedChannels
	}

	if len(searchParam.FromUsers) > 0 {
		queryParams["from_users"] = searchParam.FromUsers
	}

	if len(searchParam.ExcludedUsers) > 0 {
		queryParams["exclude_users"] = searchParam.ExcludedUsers
	}
	return
}

func composeSearchQuery(fieldName *string, searchParam *model.SearchParams) (param string) {
	mlog.Warn("BIGGO-INDEXER", mlog.Any("SEARCH_PARAM", searchParam))
	should := []string{}
	must := []string{}
	must_not := []string{}

	if searchParam.Terms != "" {
		field := *fieldName
		if searchParam.IsHashtag {
			field = "hashtags"
		}
		if searchParam.OrTerms {
			should = append(should, composeTermsQuery(field))
		} else {
			must = append(must, composeTermsQuery(field))
		}
	}

	if searchParam.ExcludedTerms != "" {
		field := *fieldName
		if searchParam.IsHashtag {
			field = "hashtags"
		}
		must_not = append(must_not, composeTermsQuery(field))
	}

	// filter out system messages for: joining team, joining channel, archiving channel
	//must_not = append(must_not, "{terms: {type: ['system_join_team', 'system_join_channel', 'system_channel_deleted']}}")

	if len(searchParam.InChannels) > 0 {
		must = append(must, "{terms: {channel_id: $in_channels}}")
	}

	if len(searchParam.ExcludedChannels) > 0 {
		must_not = append(must_not, "{terms: {channel_id: $exclude_channels}}")
	}

	if len(searchParam.FromUsers) > 0 {
		must = append(must, "{terms: {user_id: $from_users}}")
	}

	if len(searchParam.ExcludedUsers) > 0 {
		must_not = append(must_not, "{terms: {user_id: $exclude_users}}")
	}

	builder := strings.Builder{}
	builder.WriteString("{ bool: {")
	if len(should) > 0 {
		builder.WriteString("should: [")
		builder.WriteString(strings.Join(should, ","))
		builder.WriteString("]")
	}
	if len(must) > 0 {
		builder.WriteString("must: [")
		builder.WriteString(strings.Join(must, ","))
		builder.WriteString("]")
	}
	if len(must_not) > 0 {
		if len(must) > 0 {
			builder.WriteString(",")
		}
		builder.WriteString("must_not: [")
		builder.WriteString(strings.Join(must_not, ","))
		builder.WriteString("]")
	}
	builder.WriteString("}}")

	param = builder.String()
	return
}

func composeTermsQuery(fieldName string) string {
	// TODO: weigh the results based on terms, match, prefix, keyword
	query := strings.Builder{}
	query.WriteString("{bool:{should:[")
	query.WriteString(fmt.Sprintf("{match: {%s: {query: $term, boost: 4}}},", fieldName))
	query.WriteString(fmt.Sprintf("{prefix: {%s: {value: $term, boost: 2}}},", fieldName))
	query.WriteString(fmt.Sprintf("{term: {%s: {value: $term, boost: 0}}},", fieldName))
	query.WriteString(fmt.Sprintf("{term: {`%s.keyword`: {value: $term, boost: 0}}}", fieldName))
	query.WriteString("]}}")
	return query.String()
}
