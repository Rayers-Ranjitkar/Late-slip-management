package controllers

import (
	"context"
	"lateslip/initialializers"
	"lateslip/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequestLateSlip(c *gin.Context) {
	//TODO: need to check if student's late slip request limit is reached or not (max is 4 per semester)
	//If the limit is reached, return an error response

	//get student ID from context and reason from request body
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	studentID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	type requestBody struct {
		Reason string `json:"reason" binding:"required"`
	}
	var body requestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//create a new late slip request
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlip := models.LateSlip{
		ID:        primitive.NewObjectID(),
		StudentID: studentID,
		Reason:    body.Reason,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//insert the late slip into the database
	lateSlipCollection := initialializers.DB.Collection("lateslips")
	_, err = lateSlipCollection.InsertOne(ctx, lateSlip)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create late slip"})
		return
	}

	//TODO: send notification to admin
	// This could be done via email, push notification, etc.

	//return the late slip
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Late slip created successfully", "lateSlip": lateSlip})

}

func ApproveLateSlip(c *gin.Context) {
	// Get late slip ID from URL params
	lateSlipID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid late slip ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	var lateSlip models.LateSlip
	err = lateSlipCollection.FindOne(ctx, bson.M{"_id": lateSlipID}).Decode(&lateSlip)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch late slip",
		})
		return
	}
	//TODO: check if student's lateslip limit is reached or not (max is 4 per semester)
	// This could involve checking the number of approved late slips for the student

	//update the late slip status
	lateSlip.Status = "approved"
	lateSlip.UpdatedAt = time.Now()

	//update the late slip in the database
	_, err = lateSlipCollection.UpdateOne(ctx, bson.M{"_id": lateSlipID}, bson.M{"$set": lateSlip})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update late slip"})
		return
	}

	//TODO: send notification to student
	// This could be done via email, push notification, etc.

	//return the late slip
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Late slip approved successfully", "lateSlip": lateSlip})

}

func GetAllLateSlips(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	cursor, err := lateSlipCollection.Find(ctx, bson.M{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch late slips"})
		return
	}
	defer cursor.Close(ctx)

	var lateSlips []models.LateSlip
	if err = cursor.All(ctx, &lateSlips); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode late slips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "lateSlips": lateSlips})
}
