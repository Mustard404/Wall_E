package plugins

import (
	"Agent/common"
	"github.com/zan8in/masscan"
	"log"
)

func GoMasscan(ip string) []int {
	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(ip),
		masscan.SetParamPorts("1-65535"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(common.MSConfig.Rate),
	)
	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	scanResult, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("masscan encountered an error: %v", err)
	}

	if scanResult != nil {
		var ports []int
		for i := range scanResult.Hosts {
			ports = append(ports, scanResult.Ports[i].Port)
		}
		return ports
	}

	return nil
}
