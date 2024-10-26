package cluster

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/gossip"
)

func (p *BiggoCluster) Call(callback func(client gossip.ClusterClient, node *model.ClusterDiscovery), skipSelf bool, nodes ...string) (err error) {
	var cluster []*model.ClusterDiscovery
	cluster, err = p.PlatformService.Store.ClusterDiscovery().GetAll(
		p.ClusterDiscoveryService.Type, p.ClusterDiscoveryService.ClusterName,
	)

	connectionErrors := []error{}
	connectionErrorsMutex := sync.Mutex{}

	waitGroup := sync.WaitGroup{}
	for idx := range cluster {
		if skipSelf && cluster[idx].Id == p.GetClusterId() {
			continue
		}
		if len(nodes) > 0 && !slices.Contains[[]string](nodes, cluster[idx].Id) {
			continue
		}

		waitGroup.Add(1)
		go func(hostname string, port int32) {
			defer waitGroup.Done()
			if connection, err := gossip.NewClusterConnection(hostname, port); err != nil {
				connectionErrorsMutex.Lock()
				defer connectionErrorsMutex.Unlock()
				connectionErrors = append(
					connectionErrors,
					fmt.Errorf("[cluster](%s:%d) connection error: %v", hostname, port, err))
			} else {
				defer connection.Close()
				callback(gossip.NewClusterClient(connection), cluster[idx])
			}
		}(cluster[idx].Hostname, cluster[idx].GossipPort)
		waitGroup.Wait()
	}

	if len(connectionErrors) > 0 {
		err = errors.Join(connectionErrors...)
	}
	return
}
