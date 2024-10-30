package pluginapi_test

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/pluginapi"

	"git.biggo.com/Funmula/BigGoChat/server/public/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	client *pluginapi.Client
}

func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)

	return nil
}

func Example() {
}
