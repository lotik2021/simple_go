package session

import (
	"testing"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	var (
		userID   = 10
		deviceID = "asd"
	)

	ctx := common.NewInternalContext()
	ctx.DeviceID = deviceID

	err := Create(ctx, userID)
	assert.Nil(t, err)

	uds, err := Find(ctx, userID)
	if err != nil {
		assert.Equal(t, models.NewUnauthorizedError(), err)
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, uds)
	assert.Equal(t, userID, uds.UserID)
	assert.Equal(t, ctx.DeviceID, uds.DeviceID)

	err = Logout(ctx, userID)
	assert.Nil(t, err)

	_, err = ctx.DB.Model(&uds).Exec("DELETE FROM ?TableName WHERE id = ?id")
	assert.Nil(t, err)
}
