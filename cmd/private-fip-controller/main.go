package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/rwunderer/hcloud-private-fip-controller/internal/app/fipcontroller"
    "github.com/rwunderer/hcloud-private-fip-controller/internal/pkg/config"
)

func main() {
    configFile := os.Getenv("CONFIG_FILE")

    config := &config.Config{}

    // define command-line flags
    flag.StringVar(&config.ApiToken, "hcloud-token", os.Getenv("HCLOUD_TOKEN"), "Hetzner Cloud API Token")
    flag.StringVar(&config.HostIp, "host-ip", os.Getenv("HOST_IP"), "IP address of local node")
    flag.StringVar(&config.IpAddress, "ip-address", os.Getenv("IP_ADDRESS"), "Internal ip address to use as 'floating' ip")
    flag.StringVar(&config.NetworkName, "hcloud-network", os.Getenv("HCLOUD_NETWORK"), "ID or name of the Hetzner Cloud private network")

    // optionally read flags from config file
    if _, err := os.Stat(configFile); err == nil {
        if err := config.ReadFile(configFile); err != nil {
            log.Print(fmt.Errorf("could not parse controller config file: %v", err))
            os.Exit(1)
        }
    }

    // parse command-line flags
    flag.Parse()

    // run actual controller
    controller, err := fipcontroller.NewController(config)
    if err != nil {
        log.Print(fmt.Errorf("could not initialize controller: %v", err))
        os.Exit(1)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    controller.Run(ctx)
}
