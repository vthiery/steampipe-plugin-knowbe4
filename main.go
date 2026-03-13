package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/vthiery/steampipe-plugin-knowbe4/knowbe4"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: knowbe4.Plugin})
}
