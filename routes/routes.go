package routes

import (
	"token_based_auth/handlers"

	"github.com/gin-gonic/gin"
)

func NewRoute() *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	v1 := app.Group("v1")
	v2 := app.Group("v2")

	device := v1.Group("device")
	device.GET("/all", handlers.GetAllDevices)
	device.POST("/add_device", handlers.AddNewDevice)
	device.DELETE("/delete_device", handlers.DeleteDevices)
	device.PATCH("/update_device", handlers.UdpateNewDevice)

	user_kc := v1.Group("user")

	user_kc.POST("/login", handlers.LoginUser)
	user_kc.POST("/create", handlers.CreateNewUser)
	user_kc.GET("/", handlers.GetAllUsers)
	user_kc.POST("/update", handlers.UpdateUser)
	user_kc.DELETE("/delete", handlers.DeleteUserByID)

	user_es := v2.Group("user")

	user_es.POST("/create", handlers.CreateNewUserES)
	user_es.GET("/", handlers.GetAllUsersES)
	user_es.POST("/update", handlers.UpdateUserES)
	user_es.DELETE("/delete", handlers.DeleteUsersES)
	return app
}
