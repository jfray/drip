package drip_client

import (
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
)

func (cc *DripClient) Destroy(IDs ...int) ([]*godo.Response, error) {
	allResponses := make([]*godo.Response, 0)
	allErrors := make([]string, 0)

	// This is an actual fatal error in this case, so let's bomb out
	if len(IDs) > cc.MainConfig.MaxHosts {
		return nil, fmt.Errorf(
			"Some people just want to see the world burn. Others just " +
				"mistype how many machines they want to delete concurrently." +
				"Bill from Steel Mountain just wants to take instagram pics " +
				"of him and his cats. Config is set to %d right now.",
		)

	}

	// We don't necessarily want to error out the entire function if not all
	// hosts are destroyed. Let's keep track of both indepedently and maybe
	// requeue the failures?
	//
	// TODO: figure out a good queue strategy
	for _, ID := range IDs {
		resp, err := cc.Client.Droplets.Delete(ID)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		allResponses = append(allResponses, resp)
	}

	return allResponses, fmt.Errorf(strings.Join(allErrors, ","))
}
