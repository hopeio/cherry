package neti

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

func GetIP() string {
	ip, _ := ExternalIP()
	return ip.String()
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
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
	return nil, errors.New("network error")
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
		return nil // not an ipv4 address
	}

	return ip
}

func GetCommonIP() (string, error) {
	res, err := http.Get("http://txt.go.sohu.com/ip/soip")
	if err != nil {
		return "", errors.New("network error")
	}
	body, _ := ioutil.ReadAll(res.Body)
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	return string(reg.Find(body)), nil
}

// 获取本机ip地址
func GetLocalIPv4Address() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no ipv4 address found")
}
