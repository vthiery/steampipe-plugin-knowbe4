package knowbe4

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// knowbe4Config stores the connection configuration for the plugin.
type knowbe4Config struct {
	APIKey    *string `hcl:"api_key"`
	APIRegion *string `hcl:"api_region"`
}

// ConfigInstance returns a new instance of knowbe4Config (used by the SDK).
func ConfigInstance() interface{} {
	return &knowbe4Config{}
}

// GetConfig retrieves and casts the connection config from the plugin query data.
func GetConfig(connection *plugin.Connection) knowbe4Config {
	if connection == nil || connection.Config == nil {
		return knowbe4Config{}
	}
	config, _ := connection.Config.(knowbe4Config)
	return config
}
