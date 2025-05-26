package main

import (
	"lateslip/controllers"
	"lateslip/initialializers"

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
        
    }
	

	r.Run()

}