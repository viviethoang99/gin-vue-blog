package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"xojoc.pw/useragent"
)

var IP = new(ipUtil)

type ipUtil struct{}

// Get the IP address from which the user sent the request
// If the server is not behind a proxy, you can get the IP directly via c.Request.RemoteAddr
// In common architectures, requests usually pass through a proxy (most commonly Nginx) before reaching the server; directly obtaining the IP yields the proxy server's IP
func (*ipUtil) GetIpAddress(c *gin.Context) (ipAddress string) {
	// c.ClientIP() gets the proxy server's IP (Nginx)

	// X-Real-IP: Nginx proxy header; since this project explicitly uses Nginx, prefer this first
	ipAddress = c.Request.Header.Get("X-Real-IP")

	// X-Forwarded-For is added when passing through HTTP proxies or load balancers
	// Format: client1,proxy1,proxy2
	// Typically, the first IP is the real client IP; the rest are proxy servers
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ips := c.Request.Header.Get("X-Forwarded-For") // "ip1,ip2,ip3"
		splitIps := strings.Split(ips, ",")            // ["ip1", "ip2", "ip3"]
		if len(splitIps) > 0 {
			ipAddress = splitIps[0]
		}
	}

	// Proxy-Client-IP: Apache proxy header
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("Proxy-Client-IP")
	}

	// WL-Proxy-Client-IP: Weblogic proxy header
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
	}

	// RemoteAddr: the remote host IP of the request (will be the proxy IP if behind a proxy)
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.RemoteAddr
	}

	// If local IP is detected, fetch the LAN IP address
	if strings.HasPrefix(ipAddress, "127.0.0.1") || strings.HasPrefix(ipAddress, "[::1]") {
		ip, err := externalIP()
		if err != nil {
			slog.Error("GetIpAddress, externalIP, err: ", err)
		}
		ipAddress = ip.String()
	}

	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Index(ipAddress, ",") > 0 {
			ipAddress = ipAddress[:strings.Index(ipAddress, ",")]
		}
	}
	return ipAddress
}

// Get IP source
// https://github.com/lionsoul2014/ip2region
var vIndex []byte // Cache VectorIndex index to reduce a fixed IO operation

// Get region info: China|0|Jiangsu Province|Suzhou City|Telecom
func (*ipUtil) GetIpSource(ipAddress string) string {
	var dbPath = "../assets/ip2region.xdb" // IP database file
	// File-only query: read from file each time
	// searcher, err := xdb.NewWithFileOnly(dbPath)

	// Cache VectorIndex to reduce a fixed IO operation
	if vIndex == nil {
		var err error
		vIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to load vector index from `%s`: %s\n", dbPath, err))
			return ""
		}
	}
	searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex)

	if err != nil {
		slog.Error("failed to create searcher with vector index: ", err)
		return ""
	}
	defer searcher.Close()

	// Format: Country|Region|Province|City|ISP
	// Only China's data is mostly accurate to city; for other countries, often only to country; the rest are 0
	region, err := searcher.SearchByStr(ipAddress)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to search ip(%s): %s\n", ipAddress, err))
		return ""
	}
	return region
}

// Get simplified IP source, e.g., "Jiangsu Suzhou Telecom"
func (i *ipUtil) GetIpSourceSimpleIdle(ipAddress string) string {
	region := i.GetIpSource(ipAddress) // Country|Region|Province|City|ISP

	// Detected as intranet, return "Intranet IP"
	// Example: 0|0|0|Intranet IP|Intranet IP
	if strings.Contains(region, "内网IP") {
		return "内网IP"
	}

	// Often unable to get region
	// Example: China|0|Jiangsu Province|Suzhou City|Telecom
	ipSource := strings.Split(region, "|")
	if ipSource[0] != "中国" && ipSource[0] != "0" {
		return ipSource[0]
	}
	if ipSource[2] == "0" {
		ipSource[2] = ""
	}
	if ipSource[3] == "0" {
		ipSource[3] = ""
	}
	if ipSource[4] == "0" {
		ipSource[4] = ""
	}
	if ipSource[2] == "" && ipSource[3] == "" && ipSource[4] == "" {
		return ipSource[0]
	}
	return ipSource[2] + ipSource[3] + " " + ipSource[4]
}

func (*ipUtil) GetUserAgent(c *gin.Context) *useragent.UserAgent {
	return useragent.Parse(c.Request.UserAgent())
}

// Get LAN IP that is not 127.0.0.1
func externalIP() (net.IP, error) {
	// Get the server's network interface list
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		// Not up
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// Unicast interface address list
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
