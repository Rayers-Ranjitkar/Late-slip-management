package events

import (
	"context"
	"encoding/json"
	"lateslip/initialializers"
	"lateslip/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var clientManager = NewClientManager()

// WebSocketHandler handles WebSocket connections for both admin and student clients
func WebSocketHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	isAdmin := role == "admin"

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Create new client
	client := NewClient(conn, userID, isAdmin, clientManager)

	// Register client with manager
	clientManager.Register(client)

	// Send initial connection message
	initialMessage := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: "CONNECTED",
		Data: gin.H{
			"message": "WebSocket connection established",
			"time":    time.Now(),
		},
	}

	client.SendJSON(initialMessage)
}

func StartScheduleNotifier() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			checkUpcomingClasses()
		}
	}()
}

func checkUpcomingClasses() {
	ctx := context.Background()
	scheduleCollection := initialializers.DB.Collection("schedules")

	// Get current time
	now := time.Now()
	// Check for classes starting in next 15 minutes
	fifteenMinsFromNow := now.Add(15 * time.Minute)

	// Find schedules for current day and time
	filter := bson.M{
		"day": now.Weekday().String(),
		"startTime": bson.M{
			"$gte": now.Format("15:04"),
			"$lte": fifteenMinsFromNow.Format("15:04"),
		},
	}

	cursor, err := scheduleCollection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var schedules []models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		return
	}

	// Send notifications for each upcoming class
	for _, schedule := range schedules {
		notification := map[string]interface{}{
			"type": "CLASS_REMINDER",
			"data": map[string]interface{}{
				"moduleCode": schedule.ModuleCode,
				"moduleName": schedule.ModuleName,
				"startTime":  schedule.StartTime,
				"room":       schedule.RoomName,
				"message":    "Your class starts in 15 minutes",
			},
		}

		// Convert notification to JSON
		jsonMsg, err := json.Marshal(notification)
		if err != nil {
			continue
		}

		// Get all students enrolled in this class
		// You'll need to implement this based on your data model
		studentIDs := getStudentsForClass(schedule.Level)

		// Send notification to each student
		for _, studentID := range studentIDs {
			if client, exists := clientManager.studentClients[studentID]; exists {
				client.Send <- jsonMsg
			}
		}
	}
}

func getStudentsForClass(level string) []string {
	ctx := context.Background()
	studentCollection := initialializers.DB.Collection("students")

	filter := bson.M{"level": level} // MongoDB field names are case-sensitive

	cursor, err := studentCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding students: %v", err)
		return []string{} // Return empty slice instead of nil
	}
	defer cursor.Close(ctx)

	var students []models.Student
	if err := cursor.All(ctx, &students); err != nil {
		log.Printf("Error decoding students: %v", err)
		return []string{}
	}

	// Extract student IDs
	studentIDs := make([]string, 0, len(students))
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID.Hex())
	}

	return studentIDs
}
