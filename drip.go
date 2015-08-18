package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jfray/drip/auth"
	"github.com/jfray/drip/client"
	"github.com/jfray/drip/config"
	"github.com/spf13/cobra"
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

func main() {
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS
	flag.Parse()
	fullPath, err := filepath.Abs(confDir)
	if err != nil {
		log.Fatalf("Error with conf directory: %q", err)
	}

	// import the main json config file
	var mainConfig config.MainConfig
	mainConfig.FilePath = fullPath
	mConfigFile, err := ioutil.ReadFile(
		fmt.Sprintf(
			"%s/conf.json",
			mainConfig.FilePath,
		),
	)
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
	clusterConfig.FilePath = fmt.Sprintf(
		"%s/clusters/%s/conf.json",
		fullPath,
		*cluster,
	)

	cConfigFile, err := ioutil.ReadFile(clusterConfig.FilePath)
	if err = json.Unmarshal(cConfigFile, &clusterConfig); err != nil {
		log.Fatalf("Could not unmarshal cluster config file: %q", err)
	}

	cc := drip_client.DripClient{
		Client:        client,
		MainConfig:    mainConfig,
		ClusterConfig: clusterConfig,
	}
	// COMMON PORTION BETWEEN ALL SUBCOMMANDS

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List running cluster instances",
		Long:  `List all running instances for the chosen cluster`,
		Run: func(cmd *cobra.Command, args []string) {
			ccList, err := cc.List()
			if err != nil {
				log.Fatalf("Could not retrieve droplet list: %q", err)
			}
			for _, l := range ccList {
				log.Printf("Name: %s, ID: %d", l.Name, l.ID)
			}
		},
	}

	var machineCount int
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create one or more new instances",
		Long: `this will build one or more new amachines based on the existin
		g cluster configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			droplets, err := cc.Create(machineCount)
			if err != nil {
				log.Fatalf("Could not create droplets: %q", err)
			}
			for _, droplet := range droplets {
				log.Printf(
					"Created. Name: %s, ID: %d, DATA: %+v",
					droplet.Name,
					droplet.ID,
					droplet,
				)
			}
		},
	}

	createCmd.Flags().IntVarP(
		&machineCount,
		"machine_count",
		"m",
		1,
		"How many machines to build at a time.",
	)

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show information about a specific instance.",
		Long: `get all information currently known about an image, keyed by the 
		image ID`,
		Run: func(cmd *cobra.Command, args []string) {
			ID, err := strconv.Atoi(strings.Join(args, " "))
			if err != nil {
				log.Fatalf(
					"Could not convert argument %s to valid image ID (%q)",
					strings.Join(args, " "),
					err,
				)
			}

			droplet, err := cc.Show(ID)
			if err != nil {
				log.Fatalf("Could not access image ID: %d (%q)", ID, err)
			}

			log.Printf(
				"Name: %s, ID: %d, DATA: %+v",
				droplet.Name,
				droplet.ID,
				droplet,
			)
		},
	}

	var destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy images by ID",
		Long:  `destroy the image referenced by their ID`,
		Run: func(cmd *cobra.Command, args []string) {
			ID, err := strconv.Atoi(strings.Join(args, " "))
			if err != nil {
				log.Fatalf(
					"Could not convert arguments %s to valid image ID (%q)",
					strings.Join(args, " "),
					err,
				)
			}
			response, err := cc.Destroy(ID)
			if err != nil {
				log.Fatalf("Could not delete image ID: %d (%q)", ID, err)
			}
			log.Printf("Response: %s", response)
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(listCmd, createCmd, showCmd, destroyCmd)
	rootCmd.Execute()
}
