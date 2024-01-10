package fipcontroller

import (
	"context"
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"k8s.io/client-go/util/retry"
)

func newHetznerClient(token string) (*hcloud.Client, error) {
	hetznerClient := hcloud.NewClient(hcloud.WithToken(token))
	return hetznerClient, nil
}

func (controller *Controller) getNetwork(ctx context.Context, netName string) (netw *hcloud.Network, err error) {
	err = retry.OnError(controller.Backoff, alwaysRetry, func() error {
		netw, _, err = controller.HetznerClient.Network.Get(ctx, netName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("could not get network %v: %v", netName, err)
	}

	return netw, nil
}

func (controller *Controller) verifyRoute(ctx context.Context, netw *hcloud.Network, ip net.IP, gw net.IP) (bool, error) {
	for _, r := range netw.Routes {
		if r.Destination.Contains(ip) && r.Gateway.Equal(gw) {
			return true, nil
		}
	}

	return false, nil
}

func (controller *Controller) setRoute(ctx context.Context, netw *hcloud.Network, ip net.IP, gw net.IP) error {
	for _, r := range netw.Routes {
		if r.Destination.Contains(ip) {
			route := hcloud.NetworkDeleteRouteOpts{
				Route: r,
			}

			err := retry.OnError(controller.Backoff, alwaysRetry, func() error {
				_, _, err := controller.HetznerClient.Network.DeleteRoute(ctx, netw, route)
				return err
			})

			if err != nil {
				return fmt.Errorf("Could not delete old route to %v: %v", ip, err)
			}
		}
	}

	route := hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: &net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(32, 32),
			},
			Gateway: gw,
		},
	}

	err := retry.OnError(controller.Backoff, alwaysRetry, func() error {
		_, _, err := controller.HetznerClient.Network.AddRoute(ctx, netw, route)
		return err
	})

	if err != nil {
		return fmt.Errorf("Could not add route to %v: %v", ip, err)
	}

	return nil
}
