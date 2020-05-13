package config

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/yaml.v2"
	"log"
	"time"
)

var (
	etcdClient *clientv3.Client
)

func LoadRemote() {
	var err error
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{C.ETCD.URL},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatalf("cannot connect to etcd - %s", err)
	}

	loadPromptsPlacesTariffs()

	loadRateLimits()

	go WatchRemote()
}

func WatchRemote() {
	pptch := etcdClient.Watch(context.Background(), C.ETCD.BasePath, clientv3.WithPrefix())
	for wresp := range pptch {
		for _, ev := range wresp.Events {
			if ev.IsModify() {
				if string(ev.Kv.Key) == C.ETCD.Paths.PromptsPlacesTariffs {
					err := parsePromptsPlacesTariffs(ev.Kv.Value)
					if err != nil {
						log.Printf("cannot parse key %s - %s\n", C.ETCD.Paths.PromptsPlacesTariffs, err)
					}
				}

				if string(ev.Kv.Key) == C.ETCD.Paths.RateLimits {
					err := parseRateLimits(ev.Kv.Value)
					if err != nil {
						log.Printf("cannot parse key %s - %s\n", C.ETCD.Paths.RateLimits, err)
					}
				}

				if string(ev.Kv.Key) == C.ETCD.Paths.BlackList {
					err := parseBlackList(ev.Kv.Value)
					if err != nil {
						log.Printf("cannot parse key %s - %s\n", C.ETCD.Paths.BlackList, err)
					}
				}
			}
		}
	}
}

func PutNewBlackList(blacklistIPs []string) {
	var (
		blackList struct {
			BlackList []string `json:"blackList" yaml:"blackList"`
		}
	)

	blackList.BlackList = blacklistIPs
	blackListVal, err := yaml.Marshal(blackList)
	if err != nil {
		log.Printf("cannot marshall yaml - %s", err)
	}
	_, err = etcdClient.Put(context.Background(), C.ETCD.Paths.BlackList, string(blackListVal))
	if err != nil {
		log.Printf("cannot put new blacklist values - %s", err)
	}
}
