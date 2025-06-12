package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	name           string
	requestBody    map[string]interface{}
	expectedStatus int
	expectedBody   map[string]interface{}
}

func TestNewGatewayHandler(t *testing.T) {
	// Test cases organized by category
	testCases := []struct {
		category string
		cases    []testCase
	}{
		{
			category: "Valid Requests",
			cases: []testCase{
				{
					name: "Basic valid request",
					requestBody: map[string]interface{}{
						"id":     "123",
						"method": "test",
						"params": map[string]interface{}{
							"prompt": "Hello, world!",
						},
					},
					expectedStatus: http.StatusOK,
					expectedBody: map[string]interface{}{
						"status":  "ok",
						"message": "Request processed by SafeCtx",
					},
				},
				{
					name: "Request with sensitive data",
					requestBody: map[string]interface{}{
						"id":     "456",
						"method": "test",
						"params": map[string]interface{}{
							"prompt":     "Hello",
							"password":   "secret123",
							"api_key":    "key123",
							"other_data": "safe data",
						},
					},
					expectedStatus: http.StatusOK,
					expectedBody: map[string]interface{}{
						"status":  "ok",
						"message": "Request processed by SafeCtx",
					},
				},
			},
		},
		{
			category: "Invalid Requests",
			cases: []testCase{
				{
					name: "Missing required fields",
					requestBody: map[string]interface{}{
						"id": "123",
						// Missing method and params
					},
					expectedStatus: http.StatusBadRequest,
				},
				{
					name: "Invalid JSON format",
					requestBody: map[string]interface{}{
						"id":     "123",
						"method": "test",
						"params": "not a map", // Invalid type
					},
					expectedStatus: http.StatusBadRequest,
				},
			},
		},
		{
			category: "Security Checks",
			cases: []testCase{
				{
					name: "SQL injection attempt",
					requestBody: map[string]interface{}{
						"id":     "123",
						"method": "test",
						"params": map[string]interface{}{
							"prompt": "DROP TABLE users;",
						},
					},
					expectedStatus: http.StatusForbidden,
				},
				{
					name: "System command attempt",
					requestBody: map[string]interface{}{
						"id":     "123",
						"method": "test",
						"params": map[string]interface{}{
							"prompt": "system command: shutdown",
						},
					},
					expectedStatus: http.StatusForbidden,
				},
			},
		},
	}

	for _, category := range testCases {
		t.Run(category.category, func(t *testing.T) {
			for _, tc := range category.cases {
				t.Run(tc.name, func(t *testing.T) {
					// Setup
					body, err := json.Marshal(tc.requestBody)
					if err != nil {
						t.Fatalf("Failed to marshal request body: %v", err)
					}

					req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
					req.Header.Set("Content-Type", "application/json")
					rr := httptest.NewRecorder()

					// Execute
					handler := NewGatewayHandler()
					handler.ServeHTTP(rr, req)

					// Assert status code
					if status := rr.Code; status != tc.expectedStatus {
						t.Errorf("handler returned wrong status code: got %v want %v",
							status, tc.expectedStatus)
					}

					// If we expect a successful response, verify the response body
					if tc.expectedStatus == http.StatusOK && tc.expectedBody != nil {
						var response map[string]interface{}
						if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
							t.Fatalf("Failed to unmarshal response: %v", err)
						}

						// Check expected fields
						for key, expectedValue := range tc.expectedBody {
							if value, exists := response[key]; !exists || value != expectedValue {
								t.Errorf("Response missing or incorrect field %s: got %v want %v",
									key, value, expectedValue)
							}
						}

						// Verify sensitive data was redacted
						if params, ok := response["data"].(map[string]interface{}); ok {
							for _, sensitiveKey := range []string{"password", "api_key"} {
								if value, exists := params[sensitiveKey]; exists && value != "[REDACTED]" {
									t.Errorf("Sensitive data not redacted: %s", sensitiveKey)
								}
							}
						}
					}
				})
			}
		})
	}
}
