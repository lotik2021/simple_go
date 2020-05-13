package ratelimit

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net"
	"strings"
	"time"
)

var parseClient *gorequest.SuperAgent

func init() {
	parseClient = common.DefaultRequest.Clone().Timeout(config.C.Push.RequestTimeout)
}

type Method string

const (
	MethodAuthCreate         Method = "authCreate"
	MethodSearchAsync        Method = "searchAsync"
	MethodSaveSelectedRoutes Method = "saveSelectedRoutes"
	MethodSendSms            Method = "sendSms"
)

func UpdateBlackList() {
	for {
		time.Sleep(15 * time.Minute)
		_, body, err := parseClient.Clone().Get("https://check.torproject.org/cgi-bin/TorBulkExitList.py?ip=1.1.1.1").End()
		if err != nil {
			log.Printf("%s\n", err)
		}

		sBlackList := strings.Split(body, "\n")

		config.PutNewBlackList(sBlackList)
	}
}

func Apply(ctx common.Context, method Method) (err error) {

	// если выключен учёт лимитов или ip из "своих" - ничего не делаем
	if !config.RateLimits.Enabled || isWhiteListIP(ctx.RealIP) {
		return
	}

	if !config.RateLimits.Enabled || isBlackListIP(ctx.RealIP) {
		err = models.NewRateLimitError()
		return err
	}

	// redis key ip:method (192.168.19.100:sendSms)
	key := fmt.Sprintf("%s:%s", ctx.RealIP, method)

	// redis key deviceID:method (deviceID:sendSms)
	deviceKey := fmt.Sprintf("%s:%s", ctx.DeviceID, method)

	res := ctx.RedisClient.Incr(key)
	if err = res.Err(); err != nil {
		return
	}

	deviceRes := ctx.RedisClient.Incr(deviceKey)
	if err = res.Err(); err != nil {
		return
	}

	// если до этого ключа не было в redis - надо поставить ему время,
	// когда истечёт ограничение (expiration time)
	// договорились что 24 часа
	if res.Val() == 1 {
		expireRes := ctx.RedisClient.Expire(key, time.Hour*24)
		if err = expireRes.Err(); err != nil {
			return
		}
	}

	if deviceRes.Val() == 1 {
		expireRes := ctx.RedisClient.Expire(deviceKey, time.Hour*24)
		if err = expireRes.Err(); err != nil {
			return
		}
	}

	var methodLimit int

	switch method {
	case MethodAuthCreate:
		methodLimit = config.RateLimits.AuthCreate
	case MethodSendSms:
		methodLimit = config.RateLimits.SendSms
	case MethodSearchAsync:
		methodLimit = config.RateLimits.SearchAsync
	case MethodSaveSelectedRoutes:
		methodLimit = config.RateLimits.SaveSelectedRoutes
	}

	if int(res.Val()) > methodLimit {
		err = models.NewRateLimitError()
	}

	return
}

// проверка того, что ip из "своих" подсетей
func isWhiteListIP(ip string) (ok bool) {
	for _, v := range config.RateLimits.WhiteIPList {
		if v.Equal(net.ParseIP(ip)) {
			ok = true
			return
		}
	}

	for _, v := range config.RateLimits.WhiteCIDRList {
		if v.Contains(net.ParseIP(ip)) {
			ok = true
			return
		}
	}

	return
}

func isBlackListIP(ip string) (ok bool) {
	for _, v := range config.RateLimits.BlackIPList {
		if v.Equal(net.ParseIP(ip)) {
			ok = true
			return
		}
	}

	for _, v := range config.RateLimits.BlackCIDRList {
		if v.Contains(net.ParseIP(ip)) {
			ok = true
			return
		}
	}

	return
}
