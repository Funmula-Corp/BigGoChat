
		package main

		import (
			"git.biggo.com/Funmula/BigGoChat/server/public/model"
			"git.biggo.com/Funmula/BigGoChat/server/public/plugin"
		)

		type MyPlugin struct {
			plugin.MattermostPlugin
		}

		func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
			panic("Uncaught error")
		}

		func main() {
			plugin.ClientMain(&MyPlugin{})
		}
	