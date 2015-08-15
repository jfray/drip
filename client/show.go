package drip_client

import "log"

func (cc *DripClient) Show(ID int) error {
	// create options. initially, these will be blank
	droplet, _, err := cc.Client.Droplets.Get(ID)
	if err != nil {
		return err
	}

	log.Printf("ID: %d, Name: %s, IP: %s", droplet.ID, droplet.Name, droplet.Networks)

	return nil
}
