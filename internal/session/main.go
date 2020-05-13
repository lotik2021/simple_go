package session

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/go-pg/pg/v9"
)

func Create(ctx common.Context, userID int) (err error) {
	devices := []*UserDevice{}
	err = ctx.DB.Model(&UserDevice{}).Where(
		"device_id = ? AND user_id = ?", ctx.DeviceID, userID).Select(&devices)
	if err != nil || len(devices) == 0 {
		_, err = ctx.DB.Model(&UserDevice{DeviceID: ctx.DeviceID, UserID: userID}).Insert()
	}
	return err
}

func Find(ctx common.Context, userID int) (*UserDevice, error) {
	devices := []*UserDevice{}
	err := ctx.DB.Model(&UserDevice{}).Where(
		"device_id = ? AND user_id = ?", ctx.DeviceID, userID).Select(&devices)
	if err != nil {
		return nil, models.NewUnauthorizedError()
	}
	if len(devices) == 0 {
		return nil, models.NewUnauthorizedError()
	}
	return devices[0], nil
}

func Logout(ctx common.Context, userID int) (err error) {

	_, err = ctx.DB.Model(&UserDevice{DeviceID: ctx.DeviceID, UserID: userID}).Where("device_id = ?device_id and user_id = ?user_id").Delete()

	return
}

func FindActiveDeviceIDsByUserID(ctx common.Context) (deviceIDs []string, err error) {
	err = ctx.DB.Model((*UserDevice)(nil)).ColumnExpr("array_agg(device_id)").Where("user_id = ?", ctx.UserID).Select(pg.Array(&deviceIDs))
	return
}
