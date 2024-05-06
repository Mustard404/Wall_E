package plugins

import (
	"Agent/common"
	"context"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

func GoNmap(ip string, portID string) common.Port {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// with a 5-minute timeout.
	args := []string{"-Pn"}
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(ip),
		nmap.WithPorts(portID),
		nmap.WithCustomArguments(args...),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	var port common.Port

	// Use the results to print an example output
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, portInfo := range host.Ports {
			port = common.Port{
				// ID:       0,
				Port:     int(portInfo.ID),
				Protocol: portInfo.Protocol,
				State:    portInfo.State.State,
				Service:  portInfo.Service.Name,
				White:    false,
				AssetID:  0,
			}
		}
	}
	return port
}
