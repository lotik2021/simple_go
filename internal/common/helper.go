package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"github.com/google/uuid"
	"strings"
)

func FetchDeviceCategoryAndType(deviceId, deviceInfo, os string) (dCategory, dType string) {
	dCategory = "web"
	dType = constant.WEB

	// если android - будет uuid.Nil
	_, err := uuid.Parse(deviceId)

	if (strings.Contains(deviceInfo, "iPhone") || strings.Contains(deviceInfo, "iPad")) && err == nil {
		dCategory = "mobile"
		dType = constant.IOS

		return
	}

	if strings.Contains(os, "Android") && err != nil {
		dCategory = "mobile"
		dType = constant.Android

		return
	}

	return
}
