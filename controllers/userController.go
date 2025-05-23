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

	helper "github.com/Anshbir18/go-jwt/helpers"
	models "github.com/Anshbir18/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()
func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string , providedPassword string) (bool , string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!= nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check=false
	}
	return check, msg
}

func Signup()gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	}

	validationError := validate.Struct(user)
	if validationError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":validationError.Error()})
	}
	password := HashPassword(*user.Password)

	user.Password=&password
	count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the email"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})
		defer cancel()
		if err!= nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the phone number"})
		}

		if count >0{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
		}

		user.Created_at=time.Now().UTC()
		user.Updated_at=time.Now().UTC()
		user.ID=primitive.NewObjectID()
		hexID := user.ID.Hex()
		user.User_id = &hexID
		token, refreshToken, _ := helper.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *user.User_id)
		user.Token=&token
		user.Refresh_token=&refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		defer cancel()
		if insertErr != nil {
			msg := fmt.Sprintf("error occured while inserting the user in the database")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data":resultInsertionNumber})
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancle = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User
		var foundUser models.User

		if err :=c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		err := userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancle()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the email"})
			return
		}

		passwordIsValid,msg :=VerifyPassword(*user.Password, *foundUser.Password)


		if passwordIsValid != true{
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		}

		token,refreshToekn,err := helper.GenerateAllToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id)
		
		helper.UpdateAllToken(token,refreshToekn,*foundUser.User_id)

		err = userCollection.FindOne(ctx,bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}


//this is only accessable by admin only
func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
	matchStage, groupStage, projectStage})
defer cancel()
if err!=nil{
	c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
}
var allusers []bson.M
if err = result.All(ctx, &allusers); err!=nil{
	log.Fatal(err)
}
c.JSON(http.StatusOK, allusers[0])}
}


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
