package drip_client

import (
	"fmt"
	"log"
)

func (cc *DripClient) Destroy(IDs ...int) error {
	if len(IDs) > cc.MainConfig.MaxHosts {
		return fmt.Errorf(
			"Some people just want to see the world burn. Others just " +
				"mistype how many machines they want to delete concurrently." +
				"Bill from Steel Mountain just wants to take instagram pics " +
				"of him and his cats. Config is set to %d right now.",
		)

	}

	for _, ID := range IDs {
		resp, err := cc.Client.Droplets.Delete(ID)
		if err != nil {
			return err
		}
		log.Printf("Response: %+v", resp)
	}

	return nil
}
