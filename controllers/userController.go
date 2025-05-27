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

// Register handler
func Register(ctx *gin.Context) {
	//get username , email , password from request body
	type body struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	userCollection := initialializers.DB.Collection("users")

	//check if user already exists
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "User already exists",
		})
		return
	}

	/*TODO:
	  --need to check if the student is in the student database which will be imported from excel file which will be implemented later
	  --if the student is not in the database they have to register from the college email
	  --Email should be used to check if the student is in the database
	  --Student need to register with their college email
	*/

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(b.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error hashing password",
		})
		return
	}

	//create user
	user = models.User{
		ID:        primitive.NewObjectID(),
		Fullname:  b.Fullname,
		Email:     b.Email,
		Password:  string(hashedPassword),
		Role:      "student",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error creating user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created",
		"user":    user,
	})
}

// Login handler
func Login(ctx *gin.Context) {
	// Get email and password from the body
	type body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Check if user exists
	userCollection := initialializers.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid email or password",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Database error",
			})
		}
		return
	}

	// Check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid email or password",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   tokenString,
	})
}

// AdminRegister handler
func AdminRegister(ctx *gin.Context) {
	//get username , email , password from request body
	type body struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	//check if user already exists
	userCollection := initialializers.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "User already exists",
		})
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(b.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error hashing password",
		})
		return
	}
	//create user
	user = models.User{
		ID:        primitive.NewObjectID(),
		Fullname:  b.Fullname,
		Email:     b.Email,
		Password:  string(hashedPassword),
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error creating user",
		})
		return
	}
	//return success
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin user created successfully",
		"user":    user,
	})
}

// AdminLogin handler
func AdminLogin(ctx *gin.Context) {
	//get  email , password from request body
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var b body
	err := ctx.BindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	//check if user exists
	userCollection := initialializers.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid email or password",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Database error",
			})
		}
		return
	}

	//check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid email or password",
		})
		return
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal("Failed to sign token:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate token",
		})
		return
	}

	//return token
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   tokenString,
	})
}
