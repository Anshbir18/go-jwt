package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/Anshbir18/go-jwt/controllers"
	middleware "github.com/Anshbir18/go-jwt/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate()) // -----> this is user to check if the user is authenticated or not
	incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/users/:user_id",controller.GetUser())
}