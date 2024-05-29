package biggo

import (
	"net/http"
	"sync"
	"sync/atomic"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app/platform"
)

type PluginAPIService struct {
	ps *platform.PlatformService

	running atomic.Bool

	mx  sync.Mutex
	mux *http.ServeMux
}

var instance *PluginAPIService

func InitPluginAPI(ps *platform.PlatformService) {
	if instance == nil {
		instance = &PluginAPIService{ps: ps}
		instance.Start()
	}
}

func (s *PluginAPIService) Start() (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if !s.running.Load() {
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/api/v1/user/id", s.GetUserId)
		go func() {
			defer s.running.Store(false)
			s.running.Store(true)
			http.ListenAndServe(":9999", s.mux)
		}()
	}
	return
}

func (s *PluginAPIService) GetUserId(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	authData := r.FormValue("auth_data")
	if authData == "" {
		http.Error(w, "query auth_data missing", http.StatusBadRequest)
		return
	}

	authService := r.FormValue("auth_service")
	if authData == "" {
		http.Error(w, "query auth_service missing", http.StatusBadRequest)
		return
	}

	if user, err := s.ps.Store.User().GetByAuth(&authData, authService); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write([]byte(user.Id))
	}
}
