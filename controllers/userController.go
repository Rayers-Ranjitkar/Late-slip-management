package controllers

import (
	"lateslip/initialializers"
	"lateslip/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context){
	//get username , email , password from request body
	type body struct{
		Fullname string `json:"fullname"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil{
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	userCollection := initialializers.DB.Collection("users")

	//check if user already exists
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err == nil{
		ctx.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(b.Password), 10)
	if err != nil{
		ctx.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	//create user
	user = models.User{
		ID:primitive.NewObjectID(),
		Fullname: b.Fullname,
		Email: b.Email,
		Password: string(hashedPassword),
		Role: "student",
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil{
		ctx.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	//return success
	ctx.JSON(200, gin.H{"message": "User created", "user": user})

}

func Login (ctx *gin.Context){
	//get email and password from the body
	type body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil{
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	//check if user exists
	userCollection := initialializers.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err != nil{
		ctx.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}

	//check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password))
	if err != nil{
		ctx.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}
	//generate token
	
	//return token
}