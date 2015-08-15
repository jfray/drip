package config

type TokenSource struct {
	AccessToken string
}

type MainConfig struct {
	HostnamePrefix string `json:"hostname_prefix"`
	MaxHosts       int    `json:"max_hosts"`
	Token          string `json:"token"`
}

type ClusterConfig struct {
	Name       string `"json:name"`
	Datacenter string `"json:datacenter"`
	Image      string `json:"image"`
	Size       string `json:"size"`
	SSHKey     int    `json:"ssh_key"`
	Token      string `json:"token"`
}
