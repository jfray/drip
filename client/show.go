package drip_client

import "github.com/digitalocean/godo"

func (cc *DripClient) Show(ID int) (*godo.Droplet, error) {
	droplet, _, err := cc.Client.Droplets.Get(ID)
	if err != nil {
		return nil, err
	}

	return droplet, nil
}
