# Drip

basic go thingie for building a quick cluster of Digital Ocean droplets

## Config stuff

Relative to your checkout, there should be a .drip directory with the following contents:

```
.drip
.drip/conf.json
.drip/templates
.drip/templates/cloud-config.tmpl
.drip/clusters
.drip/clusters/[region]-[counter]
.drip/clusters/[region]-[counter]/conf.json
.drip/clusters/[region]-[counter]/ssh
.drip/clusters/[region]-[counter]/ssh/[region]-[counter]_rsa
.drip/clusters/[region]-[counter]/ssh/[region]-[counter]_rsa.pub
```

The .drip/conf.json should consist of the following:

```json
{
    "hostname_prefix": "[a short hostname prefix]",
    "max_hosts": [a reasonable number of hosts to build or destroy at once for safety reasons],
    "token": "[your main DO auth token"
}
```

The .drip/clusters/[region]-[counter]/conf.json should consist of:

```json
{
    "image": "[your core-os image, I use coreos-stable]",
    "size": "[DO uses memory to size their boxes, I use 2gb]",
    "ssh_key": [The SSH Key ID from the DO web console],
    "token": "[The discovery token from the etcd.io site]"
}
```

More info forthcoming. To be safe, I have .drip in the .gitconfig to avoid leaking any creds to your github account.

