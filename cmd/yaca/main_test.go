package main

import (
	"errors"
	"testing"
	"yaca/models"
	"yaca/pkg/config"
	"yaca/pkg/logger"
)

// Global mock variables
var (
	mockLoadEnvFunc             func() error
	mockParseArgsFunc           func() models.Args
	mockValidateArgsFunc        func(*models.Args) error
	mockHandleErrorFunc         func(error, string, ...any)
	mockGetZoneIDByNameFunc     func(string) (string, error)
	mockDoesRecordExistOnZoneFunc func(string, string) (string, error)
	mockUpdateRecordOnZoneFunc  func(string, string, models.Record) (bool, error)
	mockCreateRecordOnZoneFunc  func(string, models.Record) (bool, error)
	mockDeleteRecordOnZoneFunc  func(string, string, models.Record) (bool, error)
)

// Track if exit was called
var exitCalled bool
var exitCode int

// Override original functions with mocks
func init() {
	// Initialize config and logger for tests
	config.Load()
	logger.Init()
	
	utilsLoadEnv = func() error { return mockLoadEnvFunc() }
	utilsParseArgs = func() models.Args { return mockParseArgsFunc() }
	utilsValidateArgs = func(args *models.Args) error { return mockValidateArgsFunc(args) }
	utilsHandleError = func(err error, msg string, args ...any) {
		if err != nil {
			// In tests, we'll set flags instead of exiting
			exitCalled = true
			exitCode = 1
			// Optionally panic for tests that expect it
			if mockHandleErrorFunc != nil {
				mockHandleErrorFunc(err, msg, args...)
			}
		}
	}
	clientGetZoneIDByName = func(zoneName string) (string, error) { return mockGetZoneIDByNameFunc(zoneName) }
	clientDoesRecordExistOnZone = func(zoneID, recordName string) (string, error) { return mockDoesRecordExistOnZoneFunc(zoneID, recordName) }
	clientUpdateRecordOnZone = func(zoneID, recordID string, record models.Record) (bool, error) {
		return mockUpdateRecordOnZoneFunc(zoneID, recordID, record)
	}
	clientCreateRecordOnZone = func(zoneID string, record models.Record) (bool, error) {
		return mockCreateRecordOnZoneFunc(zoneID, record)
	}
	clientDeleteRecordOnZone = func(zoneID, recordID string, record models.Record) (bool, error) {
		return mockDeleteRecordOnZoneFunc(zoneID, recordID, record)
	}
}

func resetTestState() {
	exitCalled = false
	exitCode = 0
	mockHandleErrorFunc = nil
}

func TestUpdateRecord(t *testing.T) {
	resetTestState()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args {
		return models.Args{
			Record:   "test.example.com",
			ZoneName: "example.com",
			Target:   "192.168.1.1",
			Type:     "A",
			Proxy:    true,
			Ttl:      3600,
			Delete:   false,
		}
	}
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "test-zone-id", nil }
	mockDoesRecordExistOnZoneFunc = func(zoneID, recordName string) (string, error) { return "test-record-id", nil }
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return true, nil }
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}

	result := run()

	if exitCalled {
		t.Errorf("Exit was called unexpectedly")
	}
	if result != 0 {
		t.Errorf("Expected exit code 0, got %d", result)
	}
}

func TestCreateRecord(t *testing.T) {
	resetTestState()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args {
		return models.Args{
			Record:   "new.example.com",
			ZoneName: "example.com",
			Target:   "192.168.1.2",
			Type:     "A",
			Proxy:    false,
			Ttl:      3600,
			Delete:   false,
		}
	}
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "test-zone-id", nil }
	mockDoesRecordExistOnZoneFunc = func(zoneID, recordName string) (string, error) { return "", nil }
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) { return true, nil }
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}

	result := run()

	if exitCalled {
		t.Errorf("Exit was called unexpectedly")
	}
	if result != 0 {
		t.Errorf("Expected exit code 0, got %d", result)
	}
}

func TestDeleteRecord(t *testing.T) {
	resetTestState()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args {
		return models.Args{
			Record:   "test.example.com",
			ZoneName: "example.com",
			Delete:   true,
		}
	}
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "test-zone-id", nil }
	mockDoesRecordExistOnZoneFunc = func(zoneID, recordName string) (string, error) { return "test-record-id", nil }
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return true, nil }

	result := run()

	if exitCalled {
		t.Errorf("Exit was called unexpectedly")
	}
	if result != 0 {
		t.Errorf("Expected exit code 0, got %d", result)
	}
}

func TestGetZoneIDByNameFails(t *testing.T) {
	resetTestState()
	
	// Set up mock to panic when error handler is called
	mockHandleErrorFunc = func(err error, msg string, args ...any) {
		panic(err)
	}
	
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args { return models.Args{} }
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "", errors.New("test error") }

	run()
}

func TestDoesRecordExistOnZoneFails(t *testing.T) {
	resetTestState()
	
	// Set up mock to panic when error handler is called
	mockHandleErrorFunc = func(err error, msg string, args ...any) {
		panic(err)
	}
	
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args { return models.Args{} }
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "test-zone-id", nil }
	mockDoesRecordExistOnZoneFunc = func(zoneID, recordName string) (string, error) { return "", errors.New("test error") }

	run()
}

func TestDeleteNonExistentRecord(t *testing.T) {
	resetTestState()
	
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args {
		return models.Args{
			Record:   "nonexistent.example.com",
			ZoneName: "example.com",
			Delete:   true,
		}
	}
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "test-zone-id", nil }
	mockDoesRecordExistOnZoneFunc = func(zoneID, recordName string) (string, error) { return "", nil } // Record doesn't exist
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) {
		return false,
		errors.New("should not be called")
	}

	result := run()

	if exitCalled {
		t.Errorf("Exit was not expected to be called")
	}
	// Should return 1 for trying to delete non-existent record
	if result != 1 {
		t.Errorf("Expected exit code 1, got %d", result)
	}
}
