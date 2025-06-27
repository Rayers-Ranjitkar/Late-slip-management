package events

import (
	"context"
	"encoding/json"
	"lateslip/initialializers"
	"lateslip/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	// Set WebSocket parameters
	upgrader.EnableCompression = true
	upgrader.HandshakeTimeout = 10 * time.Second

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Set connection parameters
	conn.SetReadLimit(32 * 1024)
	conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		log.Printf("Pong received from client %s", userID)
		return nil
	})

	// Create new client
	client := NewClient(conn, userID, isAdmin, clientManager)

	// Register client and ensure cleanup
	clientManager.Register(client)
	defer clientManager.Unregister(client)

	// Send initial connection message
	initialMessage := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: "CONNECTED",
		Data: gin.H{
			"message": "WebSocket connection established",
			"time":    time.Now(),
			"userId":  userID,
		},
	}

	if err := client.SendJSON(initialMessage); err != nil {
		log.Printf("Error sending initial message: %v", err)
		return
	}

	// Create done channel for cleanup
	done := make(chan struct{})
	defer close(done)

	// Start ping/pong
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := client.writeControl(
					websocket.PingMessage,
					[]byte{},
					time.Now().Add(10*time.Second),
				); err != nil {
					log.Printf("Ping failed for client %s: %v", userID, err)
					return
				}
				log.Printf("Ping sent to client %s", userID)
			case <-done:
				return
			}
		}
	}()

	// Start message pumps in separate goroutines
	go client.WritePump()
	client.ReadPump() // This blocks until connection closes
}

func StartScheduleNotifier() {
	ticker := time.NewTicker(5 * time.Second)
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
	fifteenMinsFromNow := now.Add(60 * time.Minute)

	log.Printf("Checking for classes between %s and %s on %s",
		now.Format("15:04"),
		fifteenMinsFromNow.Format("15:04"),
		now.Weekday().String())

	// Find schedules for current day and time
	filter := bson.M{
		"day": now.Weekday().String(),
		"start_time": bson.M{
			"$gte": now.Format("15:04"),
			"$lte": fifteenMinsFromNow.Format("15:04"),
		},
	}

	cursor, err := scheduleCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding schedules: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var schedules []models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		log.Printf("Error decoding schedules: %v", err)
		return
	}

	if len(schedules) == 0 {
		log.Printf("No upcoming classes found")
		return
	}

	log.Printf("Found %d upcoming classes", len(schedules))

	// Send notifications for each upcoming class
	for _, schedule := range schedules {
		// Validate schedule data
		if schedule.ModuleCode == "" || schedule.StartTime == "" {
			log.Printf("Invalid schedule data: %+v", schedule)
			continue
		}

		notification := map[string]interface{}{
			"type": "CLASS_REMINDER",
			"data": map[string]interface{}{
				"moduleCode": schedule.ModuleCode,
				"moduleName": schedule.ModuleName,
				"startTime":  schedule.StartTime,
				"room":       schedule.RoomName,
				"message":    "Your class starts in 15 minutes",
				"level":      schedule.Level,
			},
		}

		// Convert notification to JSON
		jsonMsg, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Error marshaling notification for schedule %s: %v", schedule.ModuleCode, err)
			continue
		}

		// Get all students enrolled in this class
		studentIDs := getStudentsForClass(schedule.Level)
		if len(studentIDs) == 0 {
			log.Printf("No students found for level %s", schedule.Level)
			continue
		}

		sentCount := 0
		failedCount := 0

		// Send notification to each student
		for _, studentID := range studentIDs {
			if client, exists := clientManager.studentClients[studentID]; exists {
				select {
				case client.Send <- jsonMsg:
					sentCount++
				default:
					// Channel is full or blocked
					log.Printf("Failed to send notification to student %s: channel full", studentID)
					failedCount++
				}
			} else {
				log.Printf("Student %s not connected", studentID)
				failedCount++
			}
		}

		log.Printf("Notification stats for %s: sent=%d, failed=%d, total students=%d",
			schedule.ModuleCode, sentCount, failedCount, len(studentIDs))
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

func (c *Client) writeMessage(messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.Conn.WriteMessage(messageType, data)
}

func (c *Client) writeControl(messageType int, data []byte, deadline time.Time) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.Conn.WriteControl(messageType, data, deadline)
}
