package events

import (
	"encoding/json"
	"lateslip/models"
)

func NotifyStudent(studentID string, message string) {
	if studentChan, exists := clientManager.studentClients[studentID]; exists {
		studentChan <- message
	}
}

func NotifyAdmins(message string, lateSlip models.LateSlip) {
	msg := map[string]interface{}{
		"type":    "NEW_LATE_SLIP_REQUEST",
		"message": message,
		"data": map[string]interface{}{
			"id":        lateSlip.ID.Hex(),
			"studentId": lateSlip.StudentID.Hex(),
			"reason":    lateSlip.Reason,
			"status":    lateSlip.Status,
			"createdAt": lateSlip.CreatedAt,
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for adminChan := range clientManager.adminClients {
		adminChan <- string(jsonMsg)
	}
}
