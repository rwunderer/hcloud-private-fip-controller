package fipcontroller

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "k8s.io/apimachinery/pkg/util/wait"

    "github.com/hetznercloud/hcloud-go/hcloud"

    "github.com/rwunderer/hcloud-private-fip-controller/internal/pkg/config"
)

type Controller struct {
    HetznerClient   *hcloud.Client
    Config          *config.Config
    Backoff         wait.Backoff
}

// NewController creates a new Controller
func NewController(config *config.Config) (*Controller, error) {

    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("Controller config invalid: %v", err)
    }

    hetznerClient, err := newHetznerClient(config.ApiToken)
    if err != nil {
        return nil, fmt.Errorf("Could not initialise hetzner client: %v", err)
    }

    backoff := wait.Backoff{
        Duration: time.Second,
        Factor:   1.2,
        Steps:    5,
    }

    return &Controller{
        HetznerClient:  hetznerClient,
        Config:         config,
        Backoff:        backoff,
    }, nil
}

func (controller *Controller) Run(ctx context.Context) error {
    log.Print("Initialization complete. Starting reconciliation")

    if err := controller.watchIp(ctx); err != nil {
        return err
    }

    for {
        select {
        case <-ctx.Done():
            log.Printf("Context Done. Shutting down")
            return nil
        case <-time.After(30 * time.Second):
            if err := controller.watchIp(ctx); err != nil {
                return err
            }
        }
    }
}

func (controller *Controller) watchIp(ctx context.Context) error {
    ip := net.ParseIP(controller.Config.IpAddress)

    if controller.isLocalAddress(ip) {
        netName := controller.Config.NetworkName
        gw      := net.ParseIP(controller.Config.HostIp)
        log.Printf("Found address %v on local interface for gateway %v.", ip, gw)

        net, err := controller.getNetwork(ctx, netName)
        if err != nil {
            log.Print(fmt.Errorf("Error getting network: %v", err))
            return nil
        }

        routeOk, err := controller.verifyRoute(ctx, net, ip, gw)
        if err != nil {
            log.Print(fmt.Errorf("Error verifying route: %v", err))
            return nil
        }

        if !routeOk {
            log.Printf("Adding route to network %v", netName)

            err := controller.setRoute(ctx, net, ip, gw)
            if err != nil {
                log.Print(fmt.Errorf("Error setting route: %v", err))
                return nil
            }
        }

    } else {
        log.Printf("Address %v is not local to %v.", ip, controller.Config.HostIp)
    }

    return nil
}

func alwaysRetry(_ error) bool {
    return true
}
