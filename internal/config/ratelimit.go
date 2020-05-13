package config

import (
	"context"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"strings"
)

var RateLimits rateLimits

type rateLimits struct {
	Enabled            bool        `json:"enabled" yaml:"enabled"`
	WhiteList          []string    `json:"whiteList" yaml:"whiteList"`
	BlackList          []string    `json:"blackList" yaml:"blackList"`
	WhiteIPList        []net.IP    `json:"white_ip_list" yaml:"-"`
	BlackIPList        []net.IP    `json:"black_ip_list" yaml:"-"`
	WhiteCIDRList      []net.IPNet `json:"white_cidr_list" yaml:"-"`
	BlackCIDRList      []net.IPNet `json:"black_cidr_list" yaml:"-"`
	AuthCreate         int         `json:"authCreate" yaml:"authCreate"`
	SearchAsync        int         `json:"searchAsync" yaml:"searchAsync"`
	SaveSelectedRoutes int         `json:"saveSelectedRoutes" yaml:"saveSelectedRoutes"`
	SendSms            int         `json:"sendSms" yaml:"sendSms"`
}

func loadRateLimits() {
	resp, err := etcdClient.Get(context.Background(), C.ETCD.Paths.RateLimits)
	if err != nil {
		log.Fatalf("cannot read key from etcd - %s", err)
	}

	for _, v := range resp.Kvs {
		if string(v.Key) == C.ETCD.Paths.RateLimits {
			err = parseRateLimits(v.Value)
			if err != nil {
				log.Fatalf("cannot read key from etcd - %s", err)
			}
		}
	}
}

func parseRateLimits(val []byte) (err error) {
	var tmpRatelimits rateLimits
	err = yaml.Unmarshal(val, &tmpRatelimits)
	if err != nil {
		return
	}

	tmpRatelimits.WhiteIPList = make([]net.IP, 0)
	tmpRatelimits.WhiteCIDRList = make([]net.IPNet, 0)

	for _, v := range tmpRatelimits.WhiteList {
		if strings.Contains(v, "/") { // если подсеть - 192.168.0.0/16
			if _, ipNet, err := net.ParseCIDR(v); err == nil {
				tmpRatelimits.WhiteCIDRList = append(tmpRatelimits.WhiteCIDRList, *ipNet)
			}
		} else { // если ip - 192.168.32.110
			if ip := net.ParseIP(v); ip != nil {
				tmpRatelimits.WhiteIPList = append(tmpRatelimits.WhiteIPList, ip)
			}
		}
	}

	tmpRatelimits.BlackIPList = make([]net.IP, 0)
	tmpRatelimits.BlackCIDRList = make([]net.IPNet, 0)

	for _, v := range tmpRatelimits.BlackList {
		if strings.Contains(v, "/") { // если подсеть - 192.168.0.0/16
			if _, ipNet, err := net.ParseCIDR(v); err == nil {
				tmpRatelimits.BlackCIDRList = append(tmpRatelimits.BlackCIDRList, *ipNet)
			}
		} else { // если ip - 192.168.32.110
			if ip := net.ParseIP(v); ip != nil {
				tmpRatelimits.BlackIPList = append(tmpRatelimits.BlackIPList, ip)
			}
		}
	}

	tmpRatelimits.BlackIPList = append(tmpRatelimits.BlackIPList, RateLimits.BlackIPList...)
	tmpRatelimits.BlackCIDRList = append(tmpRatelimits.BlackCIDRList, RateLimits.BlackCIDRList...)

	RateLimits = tmpRatelimits

	return
}

func parseBlackList(val []byte) (err error) {
	var tmpRatelimits rateLimits
	err = yaml.Unmarshal(val, &tmpRatelimits)
	if err != nil {
		return
	}

	tmpRatelimits.BlackIPList = make([]net.IP, 0)
	tmpRatelimits.BlackCIDRList = make([]net.IPNet, 0)

	for _, v := range tmpRatelimits.BlackList {
		if strings.Contains(v, "/") { // если подсеть - 192.168.0.0/16
			if _, ipNet, err := net.ParseCIDR(v); err == nil {
				tmpRatelimits.BlackCIDRList = append(tmpRatelimits.BlackCIDRList, *ipNet)
			}
		} else { // если ip - 192.168.32.110
			if ip := net.ParseIP(v); ip != nil {
				tmpRatelimits.BlackIPList = append(tmpRatelimits.BlackIPList, ip)
			}
		}
	}

	RateLimits.BlackList = tmpRatelimits.BlackList
	RateLimits.BlackIPList = tmpRatelimits.BlackIPList

	return
}
