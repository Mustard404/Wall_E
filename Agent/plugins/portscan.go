package plugins

import (
	"Agent/common"
	"encoding/json"
	"strconv"
)

func PortScan(asset common.Asset) bool {
	// Masscan 扫描端口
	portsID := GoMasscan(asset.IP)
	// Nmap 识别服务
	for _, portID := range portsID {
		portString := strconv.Itoa(portID)
		var port common.Port
		port = GoNmap(asset.IP, portString)
		port.AssetID = asset.ID
		jsonPort, _ := json.Marshal(port)
		Producer("port_scan_return", jsonPort)
	}

	return false
}
