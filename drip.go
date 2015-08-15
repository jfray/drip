package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var (
	cluster = flag.String(
		"cluster",
		"sfo1-1",
		"which cluster are we working with?",
	)
	datacenter = strings.Split(*cluster, "-")[0]
	confDir    = "./.drip"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type MainConfig struct {
	HostnamePrefix string `json:"hostname_prefix"`
	MaxHosts       int    `json:"max_hosts"`
	Token          string `json:"token"`
}

type ClusterConfig struct {
	Image  string `json:"image"`
	Size   string `json:"size"`
	SSHKey int    `json:"ssh_key"`
	Token  string `json:"token"`
}

type ConfiguredClient struct {
	Client      *godo.Client
	CloudConfig string
	MainConfig
	ClusterConfig
}

func authorize(authToken string) *godo.Client {
	tokenSource := &TokenSource{
		AccessToken: authToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	return godo.NewClient(oauthClient)
}

func (cc *ConfiguredClient) create(howMany int) (net.IP, error) {
	if howMany > cc.MainConfig.MaxHosts {
		log.Fatalf(
			"For safety measures, you cannot build more than %d hosts at a "+
				"time.",
			cc.MainConfig.MaxHosts,
		)
	}

	dropletName := fmt.Sprintf(
		"%s-%s-%s",
		cc.MainConfig.HostnamePrefix,
		*cluster,
		uuid.New()[:6],
	)
	sshKeyToUse := godo.DropletCreateSSHKey{ID: cc.ClusterConfig.SSHKey}

	createRequest := &godo.DropletCreateRequest{
		Name:              dropletName,
		Region:            datacenter,
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
		log.Fatalf("Something bad happened: %s\n\n", err)
	}
	log.Printf("My new dude is here: %d %s", newDroplet.ID, newDroplet.Name)

	return net.ParseIP("1.2.3.4"), nil
}

func main() {
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS
	flag.Parse()
	fullPath, err := filepath.Abs(confDir)
	if err != nil {
		log.Fatalf("Error with conf directory: %q", err)
	}

	// import the main json config file
	var mainConfig MainConfig
	mConfigFile, err := ioutil.ReadFile(fmt.Sprintf("%s/conf.json", fullPath))
	if err != nil {
		log.Fatalf("Could not read main config file: %q", err)
	}

	err = json.Unmarshal(mConfigFile, &mainConfig)
	if err != nil {
		log.Fatalf("Could not unmarshal main config file: %q", err)
	}
	client := authorize(mainConfig.Token)

	// get the cluster-specific json config
	var clusterConfig ClusterConfig
	cConfigFile, err := ioutil.ReadFile(
		fmt.Sprintf(
			"%s/clusters/%s/conf.json",
			fullPath,
			*cluster,
		),
	)
	if err = json.Unmarshal(cConfigFile, &clusterConfig); err != nil {
		log.Fatalf("Could not unmarshal cluster config file: %q", err)
	}

	var cloudConfigRendered bytes.Buffer
	t := template.New("cloud_config")
	cloudConfigTemplate, err := ioutil.ReadFile(
		fmt.Sprintf(
			"%s/templates/cloud-config.tmpl",
			fullPath,
		),
	)
	t, err = t.Parse(string(cloudConfigTemplate))
	if err != nil {
		log.Printf("Caught an error trying to load the template: %q", err)
	}

	if err = t.Execute(&cloudConfigRendered, clusterConfig); err != nil {
		log.Fatalf("Error caught executing template: %q", err)
	}
	cloudConfig := cloudConfigRendered.String()

	cc := ConfiguredClient{client, cloudConfig, mainConfig, clusterConfig}
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS

	cc.create(1)

}
