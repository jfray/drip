package drip_client

import (
	"log"

	"github.com/digitalocean/godo"
)

func (cc *DripClient) List() error {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := cc.Client.Droplets.List(opt)
		if err != nil {
			return err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	for _, dl := range list {
		log.Printf("ID: %d, Name: %s, IP: %s", dl.ID, dl.Name, dl.Networks)
	}

	return nil
}
