package biggoindexer

import (
	"sync"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/jobs"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine"
)

type BiggoIndexerWorker struct {
	name      string
	stateMut  sync.Mutex
	stopCh    chan struct{}
	stoppedCh chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	logger    mlog.LoggerIFace
	engine    *biggoengine.BiggoEngine
	stopped   bool
}

func MakeWorker(jobServer *jobs.JobServer, engine *biggoengine.BiggoEngine) *BiggoIndexerWorker {
	if engine == nil {
		return nil
	}
	const workerName = "BiggoIndexer"
	return &BiggoIndexerWorker{
		name:      workerName,
		stoppedCh: make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: jobServer,
		logger:    jobServer.Logger().With(mlog.String("worker_name", workerName)),
		engine:    engine,
		stopped:   true,
	}
}

func (worker *BiggoIndexerWorker) IsEnabled(cfg *model.Config) bool {
	return true
}

func (worker *BiggoIndexerWorker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *BiggoIndexerWorker) Run() {
	worker.stateMut.Lock()
	// We have to re-assign the stop channel again, because
	// it might happen that the job was restarted due to a config change.
	if worker.stopped {
		worker.stopped = false
		worker.stopCh = make(chan struct{})
	} else {
		worker.stateMut.Unlock()
		return
	}
	// Run is called from a separate goroutine and doesn't return.
	// So we cannot Unlock in a defer clause.
	worker.stateMut.Unlock()

	worker.logger.Debug("Worker Started")

	defer func() {
		worker.logger.Debug("Worker: Finished")
		worker.stoppedCh <- true
	}()

	for {
		select {
		case <-worker.stopCh:
			worker.logger.Debug("Worker: Received stop signal")
			return
		case job := <-worker.jobs:
			worker.DoJob(&job)
		}
	}
}

func (worker *BiggoIndexerWorker) Stop() {
	worker.stateMut.Lock()
	defer worker.stateMut.Unlock()

	// Set to close, and if already closed before, then return.
	if worker.stopped {
		return
	}
	worker.stopped = true
	worker.logger.Debug("Worker Stopping")
	close(worker.stopCh)
	<-worker.stoppedCh
}

func (worker *BiggoIndexerWorker) DoJob(job *model.Job) {
}
