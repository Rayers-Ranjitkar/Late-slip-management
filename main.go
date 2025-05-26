package main

import (
	"lateslip/controllers"
	"lateslip/initialializers"
	"lateslip/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initialializers.LoadEnvVariables()
	initialializers.ConnectToDB()
}

func main() {
	r := gin.Default();

    userRoutes := r.Group("/user")
    {
        userRoutes.POST("/register", controllers.Register)
		userRoutes.POST("/login", controllers.Login)
        
    }

	studentRoutes := r.Group("/student").Use(middleware.AuthMiddleware(), middleware.RequireRole("student"))
	studentRoutes.GET("/profile", func(c *gin.Context) {
		userID := c.GetString("user_id")
		role := c.GetString("role")
		c.JSON(200, gin.H{"message": "Welcome to your profile", "user_id": userID})
		c.JSON(200, gin.H{"message": "Welcome to your student profile", "user_id": userID, "role": role})
	})
		
	
	

	r.Run()

}