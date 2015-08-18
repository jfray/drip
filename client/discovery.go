package drip_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func (cc *DripClient) UpdateDiscoveryToken() error {
	var err error
	dropletList, err := cc.List()
	if err != nil {
		fmt.Errorf("Listing to the left: %q", err)
	}
	if len(dropletList) == 0 {
		log.Println("No hosts were found, creating a new token.")

		discoveryURL := fmt.Sprintf("%s/new", cc.ClusterConfig.DiscoveryURL)
		resp, err := http.Get(discoveryURL)
		if err != nil {
			return fmt.Errorf("Couldn't discover nothin: %q", err)
		}
		defer resp.Body.Close()
		newToken, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Discovery ain't got no body: %q", err)
		}
		tokenParts := strings.Split(string(newToken), "/")
		token := tokenParts[len(tokenParts)-1]

		match, _ := regexp.MatchString("^[a-f0-9]{32}$", token)
		if !match {
			return fmt.Errorf("The token received was bogus. Try again.")
		}
		// update the json file in place I guess?
		cc.ClusterConfig.Token = token
		writeContent, err := json.Marshal(cc.ClusterConfig)

		if err != nil {
			return fmt.Errorf("Jason can't write the JSON: %q", err)
		}

		if err = ioutil.WriteFile(
			cc.ClusterConfig.FilePath,
			writeContent,
			0644,
		); err != nil {
			return fmt.Errorf("Can't file this one away: %q", err)
		}

	}
	return nil
}
