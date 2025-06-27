package main

import (
	"context"
	"lateslip/controllers"
	"lateslip/events"
	"log/slog"
	"os"

	"lateslip/initialializers"
	"lateslip/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	initialializers.LoadEnvVariables()
	initialializers.ConnectToDB()
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	// Index for schedules collection
	schedulesCollection := db.Collection("schedules")
	_, err := schedulesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "day", Value: 1},
			{Key: "start_time", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	// Index for students collection
	studentsCollection := db.Collection("students")
	_, err = studentsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "level", Value: 1}},
	})
	if err != nil {
		return err
	}

	slog.Info("Indexes created successfully")
	return nil
}

func main() {
	// Initialize logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Ensure indexes
	if err := ensureIndexes(context.Background(), initialializers.DB); err != nil {
		slog.Error("Failed to create indexes", "error", err)
		os.Exit(1)
	}
	r := gin.Default()

	r.Use(cors.Default())

	r.Use(middleware.RequestIDMiddleware())

	userRoutes := r.Group("/")
	{
		userRoutes.POST("/student/register", controllers.Register)
		userRoutes.POST("/student/login", controllers.Login)
		userRoutes.POST("/admin/register", controllers.AdminRegister)
		userRoutes.POST("/admin/login", controllers.AdminLogin)

	}

	// Add WebSocket routes
	studentRoutes := r.Group("/student").Use(middleware.AuthMiddleware(), middleware.RequireRole("student"))
	{
		studentRoutes.POST("/requestLateSlip", controllers.RequestLateSlip)
		studentRoutes.GET("/studentLateslips", controllers.GetStudentsLateslip)
		// Replace SSE with WebSocket endpoint for students
		studentRoutes.GET("/ws", events.WebSocketHandler)
	}

	adminRoutes := r.Group("/admin").Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		adminRoutes.PUT("/lateslips/approve", controllers.ApproveLateSlip)
		adminRoutes.GET("/lateslips", controllers.GetAllLateSlips)
		adminRoutes.POST("/uploadStudentData", controllers.UploadStudentData)
		adminRoutes.GET("/lateslips/pending", controllers.GetAllPendingLateSlip)
		adminRoutes.PUT("/lateslips/reject", controllers.RejectLateSlip)
		adminRoutes.POST("/uploadScheduleData", controllers.UploadScheduleData)
		// Replace SSE with WebSocket endpoint for admins
		adminRoutes.GET("/ws", events.WebSocketHandler)
	}

	events.StartScheduleNotifier()

	r.Run(":8000")
}
