package dialog

import (
	"bitbucket.movista.ru/maas/maasapi/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func Add(api *echo.Group) {
	dialogs := api.Group("/dialogs", middlewares.Token)

	dialogs.POST("/createSession", createSession)
	dialogs.POST("/getSession", getSession)
	dialogs.POST("/updateSessionV2", updateSession)
	dialogs.POST("/getSessionData", getSessionData)
	dialogs.POST("/updateLastSessionData", updateLastSessionData)
	dialogs.POST("/createSessionData", createSessionData)
	dialogs.POST("/getProfile", getProfile)
	dialogs.POST("/getProfileFavorites", getDeviceFavorites)
	dialogs.POST("/updateProfile", updateProfile)
	dialogs.POST("/getUserActionCount", getUserActionCount)
	dialogs.POST("/getUserPlaceCount", getUserPlaceCount)
	dialogs.POST("/getGooglePlaceInfo", getGooglePlaceInfo)

	return
}
