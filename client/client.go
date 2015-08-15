package drip_client

import (
	"github.com/digitalocean/godo"
	"github.com/jfray/drip/config"
)

type DripClient struct {
	Client      *godo.Client
	CloudConfig string
	config.MainConfig
	config.ClusterConfig
}
