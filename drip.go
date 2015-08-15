package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/jfray/drip/auth"
	"github.com/jfray/drip/client"
	"github.com/jfray/drip/config"
)

var (
	cluster = flag.String(
		"cluster",
		"sfo1-1",
		"which cluster are we working with?",
	)
	howMany = flag.Int(
		"how_many",
		1,
		"how many machines should we build/destroy/whatever?",
	)
	datacenter = strings.Split(*cluster, "-")[0]
	confDir    = "./.drip"
)

func main() {
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS
	flag.Parse()
	fullPath, err := filepath.Abs(confDir)
	if err != nil {
		log.Fatalf("Error with conf directory: %q", err)
	}

	// import the main json config file
	var mainConfig config.MainConfig
	mConfigFile, err := ioutil.ReadFile(fmt.Sprintf("%s/conf.json", fullPath))
	if err != nil {
		log.Fatalf("Could not read main config file: %q", err)
	}

	err = json.Unmarshal(mConfigFile, &mainConfig)
	if err != nil {
		log.Fatalf("Could not unmarshal main config file: %q", err)
	}
	client := auth.Authorize(mainConfig.Token)

	// get the cluster-specific json config
	var clusterConfig config.ClusterConfig

	// Add some basic info to the config
	clusterConfig.Name = *cluster
	clusterConfig.Datacenter = datacenter

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

	cc := drip_client.DripClient{
		client,
		cloudConfig,
		mainConfig,
		clusterConfig,
	}
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS

	cc.Create(*howMany)

}
