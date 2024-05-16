package biggo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

type G2Service struct {
	cluster *BiggoCluster

	mux *http.ServeMux
	svr *http.Server

	lock    sync.Mutex
	running atomic.Bool
}

func (g2s *G2Service) StartInterNodeCommunication() {
	g2s.lock.Lock()
	defer g2s.lock.Unlock()

	if !g2s.running.Load() {
		if *g2s.cluster.ps.Config().ClusterSettings.OverrideHostname != "" {
			g2s.cluster.cds.Hostname = *g2s.cluster.ps.Config().ClusterSettings.OverrideHostname
		} else if *g2s.cluster.ps.Config().ClusterSettings.UseIPAddress {
			g2s.cluster.cds.AutoFillIPAddress(
				*g2s.cluster.ps.Config().ClusterSettings.NetworkInterface,
				*g2s.cluster.ps.Config().ClusterSettings.AdvertiseAddress,
			)
		} else {
			g2s.cluster.cds.AutoFillHostname()
		}

		g2s.cluster.cds.ClusterName = *g2s.cluster.ps.Config().ClusterSettings.ClusterName
		g2s.cluster.cds.GossipPort = (int32)(*g2s.cluster.ps.Config().ClusterSettings.GossipPort)
		g2s.cluster.cds.Type = model.GetServiceEnvironment()

		g2s.mux = http.NewServeMux()
		g2s.svr = &http.Server{
			Addr:    fmt.Sprintf(":%d", g2s.cluster.cds.GossipPort),
			Handler: g2s.mux,
		}

		g2s.mux.HandleFunc("/gossip/cluster/info", g2s.clusterInfoHandler)
		g2s.mux.HandleFunc("/gossip/cluster/message", g2s.clusterMessageHandler)
		g2s.mux.HandleFunc("/gossip/cluster/stats", g2s.clusterStatsHandler)
		g2s.mux.HandleFunc("/gossip/cluster/plugin/statuses", g2s.clusterPluginStatusesHandler)
		g2s.mux.HandleFunc("/gossip/cluster/config", g2s.clusterConfigChangedHandler)

		go func() {
			defer g2s.running.Store(false)
			defer g2s.cluster.cds.Stop()
			g2s.running.Store(true)
			g2s.cluster.cds.Start()
			g2s.svr.ListenAndServe()
		}()
	}
}

func (g2s *G2Service) StopInterNodeCommunication() {
	g2s.lock.Lock()
	defer g2s.lock.Unlock()

	if g2s.running.Load() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		g2s.svr.Shutdown(ctx)
	}
}

func (g2s *G2Service) clusterInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g2s.cluster.GetMyClusterInfo())
}

func (g2s *G2Service) GetClusterInfo(host string) (info *model.ClusterInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("http://%s:%d/gossip/cluster/info", host, g2s.cluster.cds.GossipPort), nil,
	); err != nil {
		return
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	info = new(model.ClusterInfo)
	err = json.NewDecoder(res.Body).Decode(info)
	return
}

func (g2s *G2Service) clusterMessageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		msg := new(model.ClusterMessage)
		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g2s.cluster.cbMap[msg.Event](msg)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (g2s *G2Service) PostClusterMessage(host string, msg *model.ClusterMessage) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(msg); err != nil {
		return
	}

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://%s:%d/gossip/cluster/message", host, g2s.cluster.cds.GossipPort), buffer,
	); err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	return
}

func (g2s *G2Service) clusterStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&model.ClusterStats{Id: g2s.cluster.cds.Id,
		TotalWebsocketConnections: g2s.cluster.ps.TotalWebsocketConnections(),
		TotalReadDbConnections:    g2s.cluster.ps.Store.TotalReadDbConnections(),
		TotalMasterDbConnections:  g2s.cluster.ps.Store.TotalMasterDbConnections(),
	})
}

func (g2s *G2Service) GetClusterStats(host string) (stats *model.ClusterStats, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("http://%s:%d/gossip/cluster/stats", host, g2s.cluster.cds.GossipPort), nil,
	); err != nil {
		return
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	stats = new(model.ClusterStats)
	err = json.NewDecoder(res.Body).Decode(stats)
	return
}

func (g2s *G2Service) clusterPluginStatusesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	pStats, _ := g2s.cluster.ps.GetPluginStatuses()
	json.NewEncoder(w).Encode(&pStats)
}

func (g2s *G2Service) GetClusterPluginStatuses(host string) (pStats *model.PluginStatuses, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("http://%s:%d/gossip/cluster/plugin/statuses", host, g2s.cluster.cds.GossipPort), nil,
	); err != nil {
		return
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	pStats = new(model.PluginStatuses)
	err = json.NewDecoder(res.Body).Decode(pStats)
	return
}

func (g2s *G2Service) clusterConfigChangedHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg := new(model.Config)
		if err := json.NewDecoder(r.Body).Decode(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g2s.cluster.ps.GetConfigStore().Set(cfg)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (g2s *G2Service) PostClusterConfig(host string, cfg *model.Config) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(cfg); err != nil {
		return
	}

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://%s:%d/gossip/cluster/config", host, g2s.cluster.cds.GossipPort), buffer,
	); err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	return
}
