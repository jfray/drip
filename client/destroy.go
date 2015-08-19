package drip_client

import (
	"fmt"
	"math"

	"github.com/digitalocean/godo"
)

func even(number int) bool {
	return number%2 == 0
}

func odd(number int) bool {
	return !even(number)
}

func (cc *DripClient) Destroy(ID int) (*godo.Response, error) {
	// get current count and determine whether or not deleting a node will
	// cause etcd to be unavailable. We'll use (N/2)-1 where N is number of
	// etcd nodes. Rounding up when needed. Disallowing destroy if it will
	// make etcd unusable.
	hostList, err := cc.List()
	if err != nil {
		return nil, fmt.Errorf("Can't count my chickens or my eggs: %q", err)
	}

	currentRobustitude := int(
		math.Ceil(
			(float64(len(hostList)) / float64(2)) - 1,
		),
	)

	potentialRobustitude := int(
		math.Ceil(
			(float64(len(hostList)-1) / float64(2)) - 1,
		),
	)

	// if there's only one host then there's no cluster to speak of. Bombs away!
	if len(hostList) > 1 {
		if even(len(hostList)) && potentialRobustitude < 1 {
			return nil, fmt.Errorf(
				"You have an even number of hosts (%d) in the cluster, but "+
					"removing one host will corrupt etcd values. Exiting. "+
					"Current Robustness: %d - Potential Robustness: %d",
				len(hostList),
				currentRobustitude,
				potentialRobustitude,
			)
		} else if odd(len(hostList)) && potentialRobustitude < 2 {
			return nil, fmt.Errorf(
				"You have an odd number of hosts (%d) in the cluster which "+
					"is hella rad and as it should be. However, removing "+
					"a host from the cluster will corrupt etcd values. "+
					"Exiting. Current Robustness: %d - Potential Robustness: %d",
				len(hostList),
				currentRobustitude,
				potentialRobustitude,
			)
		}
	}

	// If you've gotten this far, the host is safe to destroy
	resp, err := cc.Client.Droplets.Delete(ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
