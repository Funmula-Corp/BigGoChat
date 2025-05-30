// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchengine

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func NewBroker(cfg *model.Config) *Broker {
	return &Broker{
		cfg: cfg,
	}
}

func (seb *Broker) RegisterElasticsearchEngine(es SearchEngineInterface) {
	seb.ElasticsearchEngine = es
}

func (seb *Broker) RegisterBiggoEngine(be SearchEngineInterface) {
	seb.BiggoEngine = be
}

type Broker struct {
	cfg                 *model.Config
	ElasticsearchEngine SearchEngineInterface
	BiggoEngine         SearchEngineInterface
}

func (seb *Broker) UpdateConfig(cfg *model.Config) *model.AppError {
	seb.cfg = cfg
	if seb.ElasticsearchEngine != nil {
		seb.ElasticsearchEngine.UpdateConfig(cfg)
	}
	if seb.BiggoEngine != nil {
		seb.BiggoEngine.UpdateConfig(cfg)
	}

	return nil
}

func (seb *Broker) GetActiveEngines() []SearchEngineInterface {
	engines := []SearchEngineInterface{}
	if seb.ElasticsearchEngine != nil && seb.ElasticsearchEngine.IsActive() {
		engines = append(engines, seb.ElasticsearchEngine)
	}
	if seb.BiggoEngine != nil && seb.BiggoEngine.IsActive() && seb.BiggoEngine.IsIndexingEnabled() {
		engines = append(engines, seb.BiggoEngine)
	}
	return engines
}

func (seb *Broker) ActiveEngine() string {
	activeEngines := seb.GetActiveEngines()
	if len(activeEngines) > 0 {
		return activeEngines[0].GetName()
	}
	if *seb.cfg.SqlSettings.DisableDatabaseSearch {
		return "none"
	}
	return "database"
}
