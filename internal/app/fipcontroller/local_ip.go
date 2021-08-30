package fipcontroller

import (
	"fmt"
	"log"
	"net"
)

func (controller *Controller) isLocalAddress(needle net.IP) bool {
	var foundIface string

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return false
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
			continue
		}

		for _, a := range addrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil {
				log.Print(fmt.Errorf("localAddresses: %v\n", err))
				continue
			}

			if ip.Equal(needle) {
				foundIface = i.Name
				break
			}
		}
	}

	if foundIface != "" {
		log.Printf("Found address %v on interface %v", needle, foundIface)
		return true
	}

	return false
}
