package drip_client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"code.google.com/p/go-uuid/uuid"
	"github.com/digitalocean/godo"
)

func (cc *DripClient) Create(howMany int) ([]*godo.Droplet, error) {
	allDroplets := make([]*godo.Droplet, 0)

	if howMany > cc.MainConfig.MaxHosts {
		return nil, fmt.Errorf(
			"Enhance your calm and stop trying to build so many machines at "+
				"once, yo. Config is set to %d right now.",
			cc.MainConfig.MaxHosts,
		)
	}

	// check to see how many machines already exist. If none then generate a
	// discovery token and write it to the json
	if err := cc.UpdateDiscoveryToken(); err != nil {
		return nil, fmt.Errorf("No chill here: %q", err)
	}

	// process the cloud config template and add to the config
	var cloudConfigRendered bytes.Buffer
	t := template.New("cloud_config")
	cloudConfigTemplate, err := ioutil.ReadFile(
		fmt.Sprintf(
			"%s/templates/cloud-config.tmpl",
			cc.MainConfig.FilePath,
		),
	)
	t, err = t.Parse(string(cloudConfigTemplate))
	if err != nil {
		return nil, fmt.Errorf("Caught an error trying to load the template: %q", err)
	}

	if err = t.Execute(&cloudConfigRendered, cc.ClusterConfig); err != nil {
		return nil, fmt.Errorf("Error caught executing template: %q", err)
	}

	cloudConfig := cloudConfigRendered.String()
	cc.CloudConfig = cloudConfig

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
			return nil, fmt.Errorf("Something not so chill happened: %q", err)
		}
		allDroplets = append(allDroplets, newDroplet)

	}
	return allDroplets, nil
}
