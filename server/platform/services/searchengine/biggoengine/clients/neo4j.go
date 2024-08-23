package clients

import (
	"context"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver     neo4j.DriverWithContext = nil
	driverLock sync.Mutex              = sync.Mutex{}

	neo4jAddress  string = "bolt://172.17.0.2:7687"
	neo4jDatabase string = "neo4j"
)

func getDriver() (driver neo4j.DriverWithContext, err error) {
	if driver == nil {
		driverLock.Lock()
		defer driverLock.Unlock()
		if driver == nil {
			driver, err = neo4j.NewDriverWithContext(neo4jAddress, neo4j.NoAuth())
		}
	}
	return
}

func GraphQuery(query string, params map[string]interface{}) (result *neo4j.EagerResult, err error) {
	if driver, err = getDriver(); err != nil {
		return
	}
	result, err = neo4j.ExecuteQuery(
		context.Background(), driver, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(neo4jDatabase),
	)
	return
}
