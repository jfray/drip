package drip_client

import (
	"fmt"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/digitalocean/godo"
)

func (cc *DripClient) Create(howMany int) error {
	if howMany > cc.MainConfig.MaxHosts {
		return fmt.Errorf(
			"Enhance your calm and stop trying to build so many machines at "+
				"once, yo. Config is set to %d right now.",
			cc.MainConfig.MaxHosts,
		)
	}

	for i := 0; i < howMany; i++ {
		dropletName := fmt.Sprintf(
			"%s-%s-%s",
			cc.MainConfig.HostnamePrefix,
			cc.ClusterConfig.Name,
			uuid.New()[:6],
		)
		sshKeyToUse := godo.DropletCreateSSHKey{ID: cc.ClusterConfig.SSHKey}

		createRequest := &godo.DropletCreateRequest{
			Name:              dropletName,
			Region:            cc.ClusterConfig.Datacenter,
			Size:              cc.ClusterConfig.Size,
			SSHKeys:           []godo.DropletCreateSSHKey{sshKeyToUse},
			PrivateNetworking: true,
			UserData:          cc.CloudConfig,
			Image: godo.DropletCreateImage{
				Slug: cc.ClusterConfig.Image,
			},
		}
		newDroplet, _, err := cc.Client.Droplets.Create(createRequest)
		if err != nil {
			return fmt.Errorf("Something not so chill happened: %q", err)
		}
		log.Printf("My new dude is here: %d %s", newDroplet.ID, newDroplet.Name)

	}
	return nil
}
