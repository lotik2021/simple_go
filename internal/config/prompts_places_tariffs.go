package config

import (
	"context"
	"gopkg.in/yaml.v2"
	"log"
)

var PromptsPlacesTariffs promptsPlacesTariffs

type Lang struct {
	RU string `yaml:"ru" json:"ru"`
}

type promptsPlacesTariffs struct {
	All                       Lang `yaml:"all" json:"all"`
	Bus                       Lang `yaml:"bus" json:"bus"`
	Rail                      Lang `yaml:"rail" json:"rail"`
	Flight                    Lang `yaml:"flight" json:"flight"`
	BusAutomaticSeatSelection Lang `yaml:"busAutomaticSeatSelection" json:"bus_automatic_seat_selection"`
}

func loadPromptsPlacesTariffs() {
	resp, err := etcdClient.Get(context.Background(), C.ETCD.Paths.PromptsPlacesTariffs)
	if err != nil {
		log.Fatalf("cannot read key from etcd - %s", err)
	}

	for _, v := range resp.Kvs {
		if string(v.Key) == C.ETCD.Paths.PromptsPlacesTariffs {
			err = parsePromptsPlacesTariffs(v.Value)
			if err != nil {
				log.Fatalf("cannot read key from etcd - %s", err)
			}
		}
	}
}

func parsePromptsPlacesTariffs(val []byte) (err error) {
	var tmpPT promptsPlacesTariffs
	err = yaml.Unmarshal(val, &tmpPT)
	if err != nil {
		return
	}

	PromptsPlacesTariffs = tmpPT

	return
}
