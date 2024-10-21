package gossip

import (
	"fmt"
	"sync"

	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (g2s *GossipService) NewClient(addr string) (client ClusterClient, connection *grpc.ClientConn, err error) {
	if connection, err = grpc.Dial(fmt.Sprintf("%s:%d", addr, g2s.cds.GossipPort), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	client = NewClusterClient(connection)
	return
}

type ClusterCallFunc func(string)

func (g2s *GossipService) CallCluster(cb ClusterCallFunc, skipSelf bool) {
	if g2s.GetClusterDiscoveryService() == nil {
		return
	}

	cds := g2s.GetClusterDiscoveryService()
	if clusterDiscovery, err := g2s.ps.Store.ClusterDiscovery().GetAll(cds.Type, cds.ClusterName); err != nil {
		mlog.Error("Cluster Discovery Error", mlog.Err(err))
	} else {
		wg := sync.WaitGroup{}
		for _, cd := range clusterDiscovery {
			if cds.IsEqual(cd) && skipSelf {
				continue
			}

			wg.Add(1)
			go func(hostname string) {
				defer wg.Done()
				cb(hostname)
			}(cd.Hostname)
		}
		wg.Wait()
	}
}
