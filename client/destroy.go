package drip_client

import "github.com/digitalocean/godo"

func (cc *DripClient) Destroy(ID int) (*godo.Response, error) {
	resp, err := cc.Client.Droplets.Delete(ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
