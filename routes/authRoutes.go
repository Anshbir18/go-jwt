package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/Anshbir18/go-jwt/controllers"
)


func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("user/signup",controller.Signup())
	incomingRoutes.POST("user/login",controller.Login())
}