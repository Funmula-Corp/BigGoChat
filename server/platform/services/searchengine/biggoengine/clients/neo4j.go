package clients

import (
	"context"
	"fmt"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	neo4jClient   neo4j.DriverWithContext = nil
	neo4jDatabase string                  = "neo4j"
)

func InitNeo4jClient(config *model.Config) (err error) {
	if neo4jClient == nil {
		address := fmt.Sprintf("%s://%s:%0.f",
			cfg.Neo4jProtocol(config),
			cfg.Neo4jHost(config),
			cfg.Neo4jPort(config),
		)
		auth := neo4j.NoAuth()
		if cfg.Neo4jUseCredentials(config) {
			auth = neo4j.BasicAuth(
				cfg.Neo4jUsername(config),
				cfg.Neo4jPassword(config),
				"",
			)
		}
		neo4jClient, err = neo4j.NewDriverWithContext(address, auth)
	}
	return
}

func GraphQuery(query string, params map[string]interface{}) (result *neo4j.EagerResult, err error) {
	result, err = neo4j.ExecuteQuery(
		context.Background(), neo4jClient, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(neo4jDatabase),
	)
	return
}
