package user

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
)

// TODO
func Logout(ctx common.Context) (token string, err error) {
	userID := ctx.Token.GetUserID()
	if userID == 0 {
		err = models.NewUnauthorizedError()
		return
	}

	err = session.Logout(ctx, userID)
	if err != nil {
		return
	}

	t := models.NewDeviceToken(models.DeviceClaims{
		DeviceID: ctx.DeviceID,
	})

	token = t.String()

	ctx.Token = t

	return
}
