package common

import (
	"crypto/rand"
	"errors"
	"net"
	"strings"
)

func BuildInfo(buildVersion string, buildDate string, buildCommit string) string {
	var sb strings.Builder

	writePart := func(msg string, val string) {
		sb.WriteString(msg)

		if len(val) > 0 {
			sb.WriteString(val)
		} else {
			sb.WriteString("N/A")
		}
		sb.WriteString("\n")
	}

	writePart("Build version: ", buildVersion)
	writePart("Build date: ", buildDate)
	writePart("Build commit: ", buildCommit)

	return sb.String()
}

func GenerateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Получаем текущий локальный IP адрес
// https://gist.github.com/miguelmota/8544989558d8723b42068aec5bc72ebf
func LocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
