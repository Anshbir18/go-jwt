package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Anshbir18/go-jwt/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	models "github.com/Anshbir18/go-jwt/models"
	"github.com/Anshbir18/go-jwt/helper"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()
func HashPassword()  {}
func VerifyPassword()  {}
func Signup()  {}
func Login()  {}

//only an admin can access this route and get the users data

func GetUsers()  {}
func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
