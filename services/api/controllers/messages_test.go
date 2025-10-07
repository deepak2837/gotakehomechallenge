package controllers

import (
	"fmt"
	"net/http"
	"testing"

	"take-home-challenge/models"
	"take-home-challenge/services/api/controllers/payloads"
)

func TestMessagesController_Routes(t *testing.T) {
	message := &models.Message{
		Created:   nil,
		DeletedAt: nil,
		UserIDs:   []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"tag": "employee",
		},
	}
	url := fmt.Sprintf("%s/api/v2/messages/", server.URL)
	req, _ := http.NewRequest(http.MethodPost, url, toBody(message))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	message = &models.Message{
		Created:   nil,
		DeletedAt: nil,
		UserIDs:   []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"tag": "employee",
		},
	}
	url = fmt.Sprintf("%s/api/v2/messages/", server.URL)
	req, _ = http.NewRequest(http.MethodPost, url, toBody(message))
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	message = &models.Message{
		Created:   nil,
		DeletedAt: nil,
		UserIDs:   []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"tag": "employee",
		},
	}
	url = fmt.Sprintf("%s/api/v2/messages/", server.URL)
	req, _ = http.NewRequest(http.MethodPost, url, toBody(message))
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	url = fmt.Sprintf("%s/api/v2/messages/delete", server.URL)
	data := &payloads.MessagesMarkDeleted{
		IDs:         []int{3, 2},
		DeletedWhen: "2025-06-01",
	}
	req, _ = http.NewRequest(http.MethodPost, url, toBody(data))
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read message by ID
	url = fmt.Sprintf("%s/api/v2/messages/%d", server.URL, 1)
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	readMessage := fromBody[models.Message](resp)
	if readMessage.ID != 1 {
		t.Fatalf("Expected message ID %d, got %d", 1, readMessage.ID)
	}
}

func TestMessagesController_ReadMessage_NotFound(t *testing.T) {
	// Test reading a non-existent message
	url := fmt.Sprintf("%s/api/v2/messages/99999", server.URL)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestMessagesController_ReadMessage_InvalidID(t *testing.T) {
	// Test with invalid message ID formats
	testCases := []struct {
		messageID    string
		expectedCode int
	}{
		{"abc", http.StatusBadRequest},        // Non-numeric
		{"-1", http.StatusBadRequest},         // Negative number
		{"0", http.StatusBadRequest},          // Zero
	}

	for _, tc := range testCases {
		url := fmt.Sprintf("%s/api/v2/messages/%s", server.URL, tc.messageID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		resp, _ := http.DefaultClient.Do(req)
		if resp.StatusCode != tc.expectedCode {
			t.Fatalf("Expected status code %d for messageID '%s', got %d", tc.expectedCode, tc.messageID, resp.StatusCode)
		}
	}
}

func TestMessagesController_ReadMessage_DeletedMessage(t *testing.T) {
	// Create a message
	message := &models.Message{
		UserIDs:   []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"tag": "test",
		},
	}
	url := fmt.Sprintf("%s/api/v2/messages/", server.URL)
	req, _ := http.NewRequest(http.MethodPost, url, toBody(message))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	createdMessage := fromBody[models.Message](resp)

	// Delete the message
	deleteData := &payloads.MessagesMarkDeleted{
		IDs:         []int{createdMessage.ID},
		DeletedWhen: "2025-10-08",
	}
	url = fmt.Sprintf("%s/api/v2/messages/delete", server.URL)
	req, _ = http.NewRequest(http.MethodPost, url, toBody(deleteData))
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Try to read the deleted message - should return 404
	url = fmt.Sprintf("%s/api/v2/messages/%d", server.URL, createdMessage.ID)
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code %d for deleted message, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
