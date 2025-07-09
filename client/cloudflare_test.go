package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yaca/models"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

func setupMockServer(t *testing.T, handler http.Handler) (*httptest.Server, *cloudflare.Client) {
	server := httptest.NewTLSServer(handler)

	client := cloudflare.NewClient(
		option.WithAPIToken("test-token"),
		option.WithAPIEmail("test-email"),
		option.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
		option.WithBaseURL(server.URL),
	)
	return server, client
}

func TestGetZoneIDByName(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"result": [
				{
					"id": "test-zone-id",
					"name": "example.com"
				}
			],
			"success": true,
			"errors": [],
			"messages": []
		}`)
	})

	server, cfClient := setupMockServer(t, handler)
	defer server.Close()

	// Replace the singleton client with the mock client
	client = cfClient

	zoneID, err := GetZoneIDByName("example.com")
	if err != nil {
		t.Errorf("GetZoneIDByName() returned an error: %v", err)
	}

	if zoneID != "test-zone-id" {
		t.Errorf("GetZoneIDByName() returned incorrect zone ID, got: %s, want: %s", zoneID, "test-zone-id")
	}
}

func TestDoesRecordExistOnZone(t *testing.T) {
	t.Run("should return record ID when record exists", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"result": [
					{
						"id": "test-record-id",
						"name": "test.example.com"
					}
				],
				"success": true,
				"errors": [],
				"messages": []
			}`)
		})

		server, cfClient := setupMockServer(t, handler)
		defer server.Close()

		client = cfClient

		recordID, err := DoesRecordExistOnZone("test-zone-id", "test.example.com")
		if err != nil {
			t.Errorf("DoesRecordExistOnZone() returned an error: %v", err)
		}

		if recordID != "test-record-id" {
			t.Errorf("DoesRecordExistOnZone() returned incorrect record ID, got: %s, want: %s", recordID, "test-record-id")
		}
	})

	t.Run("should return empty string when record does not exist", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"result": [],
				"success": true,
				"errors": [],
				"messages": []
			}`)
		})

		server, cfClient := setupMockServer(t, handler)
		defer server.Close()

		client = cfClient

		recordID, err := DoesRecordExistOnZone("test-zone-id", "test.example.com")
		if err != nil {
			t.Errorf("DoesRecordExistOnZone() returned an error: %v", err)
		}

		if recordID != "" {
			t.Errorf("DoesRecordExistOnZone() returned incorrect record ID, got: %s, want: %s", recordID, "")
		}
	})
}

func TestHandleRecord(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"result": {},
			"success": true,
			"errors": [],
			"messages": []
		}`)
	})

	server, cfClient := setupMockServer(t, handler)
	defer server.Close()

	client = cfClient

	record := models.Record{
		Record: "test.example.com",
		Type:   "A",
		Target: "127.0.0.1",
		Proxy:  true,
		Ttl:    3600,
	}

	t.Run("should create record", func(t *testing.T) {
		_, err := CreateRecordOnZone("test-zone-id", record)
		if err != nil {
			t.Errorf("CreateRecordOnZone() returned an error: %v", err)
		}
	})

	t.Run("should update record", func(t *testing.T) {
		_, err := UpdateRecordOnZone("test-zone-id", "test-record-id", record)
		if err != nil {
			t.Errorf("UpdateRecordOnZone() returned an error: %v", err)
		}
	})

	t.Run("should delete record", func(t *testing.T) {
		_, err := DeleteRecordOnZone("test-zone-id", "test-record-id", record)
		if err != nil {
			t.Errorf("DeleteRecordOnZone() returned an error: %v", err)
		}
	})
}

