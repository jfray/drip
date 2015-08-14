package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/digitalocean/godo"
	"github.com/gocql/gocql"
	"golang.org/x/oauth2"
)

var (
	baseDir       = "/Users/jfray/Dropbox/Hosting/DigitalOcean"
	hostnameBase  = "a8"
	tokenLocation = flag.String(
		"token_location",
		fmt.Sprintf("%s/secret/auth_token", baseDir),
		"have your token in a file and put it here.",
	)
	configLocation = flag.String(
		"config_location",
		fmt.Sprintf("%s/conf/cloud-config.tmpl", baseDir),
		"Your cloud-config template should be here.",
	)

	datacenter = flag.String(
		"datacenter",
		"sfo1",
		"which datacenter do you want to put this thing in?",
	)
	size = flag.String(
		"size",
		"2gb",
		"what size image do you want to build?",
	)
	image = flag.String(
		"image",
		"coreos-stable",
		"which os image do you want to install?",
	)
	sshKeys = flag.Int(
		"ssh_keys",
		1201425,
		"which ssh key pair to use? (ID)",
	)
	clusterName = flag.String(
		"cluster_name",
		"sfo1-1",
		"name of unique cluster",
	)
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

type Config struct {
	Token string
}

func main() {
	flag.Parse()

	authToken, err := ioutil.ReadFile(*tokenLocation)
	if err != nil {
		log.Fatalf("Could not read token: %q", err)
	}
	log.Printf("TOKEN: %s", authToken)
	cc, err := ioutil.ReadFile(*configLocation)
	if err != nil {
		log.Fatalf("Could not read config template: %q", err)
	}
	log.Printf("CONFIG: %s", cc)

	tokenSource := &TokenSource{
		AccessToken: string(authToken),
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	client := godo.NewClient(oauthClient)
	dropletName := fmt.Sprintf("%s-%s", hostnameBase, gocql.TimeUUID().String()[:6])

	sshKeyToUse := godo.DropletCreateSSHKey{ID: *sshKeys}

	var cfg bytes.Buffer
	t := template.New("config_template") //create a new template
	t, err = t.Parse(string(cc))         //open and parse a template text file
	if err != nil {
		log.Fatalf("Caught an error trying to load the template: %q", err)
	}

	clusterTokenLocation := fmt.Sprintf("%s/conf/cluster_token-%s", baseDir, *clusterName)
	clusterToken, err := ioutil.ReadFile(clusterTokenLocation)
	if err != nil {
		log.Fatalf("Could not get cluster config token: %q", err)
	}
	log.Printf("Cluster token for %s cluster: %s", clusterName, string(clusterToken))
	config := Config{Token: string(clusterToken)}

	if err = t.Execute(&cfg, config); err != nil {
		log.Fatalf("Error caught executing template: %q", err)
	}
	cloudConfig := cfg.String()

	createRequest := &godo.DropletCreateRequest{
		Name:              dropletName,
		Region:            *datacenter,
		Size:              *size,
		SSHKeys:           []godo.DropletCreateSSHKey{sshKeyToUse},
		PrivateNetworking: true,
		UserData:          cloudConfig, // read this in as a template
		Image: godo.DropletCreateImage{
			Slug: *image,
		},
	}
	newDroplet, _, err := client.Droplets.Create(createRequest)
	if err != nil {
		log.Fatalf("Something bad happened: %s\n\n", err)
	}
	log.Printf("My new dude is here: %+v", newDroplet)
}
