[![GitHub license](https://img.shields.io/github/license/rwunderer/hcloud-private-fip-controller.svg)](https://github.com/rwunderer/hcloud-private-fip-controller/blob/main/LICENSE)

# hcloud-private-fip-controller
small k8s controller to simulate a floating ip on hetzner cloud's private network.

The configuration of the ip on the node's interface is expected to be done by
some external means. The controller runs on each possible target node and waits
for the "floating" ip to be assigned locally. As soon as it detects the ip on
a local interface it configures a corresponding route via hcloud API.

My personal use case for this feature is [Talos' virtual IP](https://www.talos.dev/docs/v1.8/guides/vip/)
on the control-plane of my cluster.

## Usage

* Define a private network with a large IP range (eg. `10.0.0.0/8`) and a smaller subnet
to place your nodes in (eg. `10.0.0.0/16`). Choose an ip outside the subnet but from
within the network range as the "floating" ip (eg. `10.255.255.1`).

* Configure some kind of vip handling for the ip (eg. Talos' virtual IP).

* Configure and launch private-fip-controller on each node.

## Configuration

The controller needs the following environment variables set:

* `HOST_IP` - the IP address to use as *destination* for the route
* `HCLOUD_TOKEN` - the [API token](https://docs.hetzner.cloud/#authentication) for Hetzner Cloud.
* `HCLOUD_NETWORK` - name or ID of the hetzner cloud private network
* `IP_ADDRESS` - the actual "floating" ip

See [deploy/daemonset.yaml](deploy/daemonset.yaml) for details.

Alternatively a yaml configuration file can be mounted into the container. Then
only the environment variable `CONFIG_FILE` needs to be set to the full path of
the configuration file:

```
hcloudToken: ABCDD....
hostIP: 1.2.3.4
ipAddress: 10.255.255.1
networkName: MyNet
```

## Acknowledgements

As this is my first attempt at both go programming and a kubernetes controller.
So I copied lots of code from various places in order to try and follow best
practices.

Among others:

* [Alex Ellis' Everyday Go Book](https://openfaas.gumroad.com/l/everyday-golang)
* [hcloud-fip-controller](https://github.com/cbeneke/hcloud-fip-controller)
* [hcloud-cloud-controller-manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager)
