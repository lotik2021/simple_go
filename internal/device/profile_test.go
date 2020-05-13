package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"fmt"
	"testing"
)

func TestGetOneWithSettingsAndFavorites(t *testing.T) {
	ctx := common.NewInternalContext()
	ctx.DeviceID = "31a093ea1af258e6"
	ctx.UserID = 2215

	res, err := GetOneWithSettingsAndFavorites(ctx)
	fmt.Println(res, err)
}
