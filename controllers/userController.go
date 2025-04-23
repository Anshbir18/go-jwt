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
)


var useCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()
func HashPassword()  {}
func VerifyPassword()  {}
func Signup()  {}
func Login()  {}

//only an admin can access this route and get the users data
func GetUser()  {} 

func GetUsers()  {}
