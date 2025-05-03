package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/Anshbir18/go-jwt/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type signedDetails struct{
	Email string
	First_name string
	Last_name string
	Uuid string
	User_type string
	jwt.StandardClaims

}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllToken(email string, first_name string, last_name string, user_type string, uuid string) (signedToken string,signedRefreshToken string,err error){
	claims := &signedDetails{
		Email: email,
		First_name: first_name,
		Last_name: last_name,
		Uuid: uuid,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			//with all the tokens we need to set the expiry time
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refeshClaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	refreshToken,err := jwt.NewWithClaims(jwt.SigningMethodHS256,refeshClaims).SignedString([]byte(SECRET_KEY))

	if err!=nil{
		log.Panic(err)
		return
	}
	return token,refreshToken,err
}

func UpdateAllToken(signedToken string, signedRefreshToken string, userId string){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	// primitive.D vs bson.M
	// primitive.D (alias for []bson.E)
	// An ordered slice of key/value pairs (bson.E{Key, Value}). Order matters if you need to preserve it (e.g. multi‑stage aggregation pipelines).

	// bson.M (alias for map[string]interface{})
	// An unordered map of keys to values. More concise when order doesn’t matter (e.g., simple filters).

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	upsert := true
	filter := bson.M{"user_id":userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	defer cancel()

	if err!=nil{
		log.Panic(err)
		return
	}
	return
}