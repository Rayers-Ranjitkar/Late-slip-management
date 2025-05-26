package controllers

import (
	"lateslip/initialializers"
	"lateslip/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	userCollection := initialializers.DB.Collection("users")

	//check if user already exists
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err == nil{
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(b.Password), 10)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	//return success
	ctx.JSON(http.StatusOK, gin.H{"message": "User created", "user": user})

}

func Login(ctx *gin.Context) {
    // Get email and password from the body
    type body struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    var b body
    err := ctx.BindJSON(&b)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Check if user exists
    userCollection := initialializers.DB.Collection("users")
    var user models.User
    err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    // Check if password is correct
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password))
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
    })
	
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
		log.Fatal("Failed to sign token:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Return token
    ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
        "token": tokenString,
    })
}