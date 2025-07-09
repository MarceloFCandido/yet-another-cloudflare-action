package main

import (
	"errors"
	"testing"
	"yaca/models"
)

// Global mock variables
var (
	mockLoadEnvFunc           func() error
	mockParseArgsFunc         func() models.Args
	mockValidateArgsFunc      func(*models.Args) error
	mockPanicOnErrorFunc      func(error)
	mockGetZoneIDByNameFunc   func(string) (string, error)
	mockDoesRecordExistOnZoneFunc func(string, string) (string, error)
	mockUpdateRecordOnZoneFunc func(string, string, models.Record) (bool, error)
	mockCreateRecordOnZoneFunc func(string, models.Record) (bool, error)
	mockDeleteRecordOnZoneFunc func(string, string, models.Record) (bool, error)
)

// Override original functions with mocks
func init() {
	utilsLoadEnv = func() error { return mockLoadEnvFunc() }
	utilsParseArgs = func() models.Args { return mockParseArgsFunc() }
	utilsValidateArgs = func(args *models.Args) error { return mockValidateArgsFunc(args) }
	utilsPanicOnError = func(err error) { 
		if err != nil {
			panic(err)
		}
	}
	clientGetZoneIDByName = func(zoneName string) (string, error) { return mockGetZoneIDByNameFunc(zoneName) }
	clientDoesRecordExistOnZone = func(zoneID, recordName string) (string, error) { return mockDoesRecordExistOnZoneFunc(zoneID, recordName) }
	clientUpdateRecordOnZone = func(zoneID, recordID string, record models.Record) (bool, error) { return mockUpdateRecordOnZoneFunc(zoneID, recordID, record) }
	clientCreateRecordOnZone = func(zoneID string, record models.Record) (bool, error) { return mockCreateRecordOnZoneFunc(zoneID, record) }
	clientDeleteRecordOnZone = func(zoneID, recordID string, record models.Record) (bool, error) { return mockDeleteRecordOnZoneFunc(zoneID, recordID, record) }
}

func TestUpdateRecord(t *testing.T) {
	var exitCode int
	osExit = func(code int) { exitCode = code }

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
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }

	main()

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestCreateRecord(t *testing.T) {
	var exitCode int
	osExit = func(code int) { exitCode = code }

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
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) { return true, nil }
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }

	main()

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestDeleteRecord(t *testing.T) {
	var exitCode int
	osExit = func(code int) { exitCode = code }

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
	mockUpdateRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }
	mockCreateRecordOnZoneFunc = func(zoneID string, record models.Record) (bool, error) { return false, errors.New("should not be called") }
	mockDeleteRecordOnZoneFunc = func(zoneID, recordID string, record models.Record) (bool, error) { return true, nil }

	main()

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestGetZoneIDByNameFails(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	mockLoadEnvFunc = func() error { return nil }
	mockParseArgsFunc = func() models.Args { return models.Args{} }
	mockValidateArgsFunc = func(args *models.Args) error { return nil }
	mockGetZoneIDByNameFunc = func(zoneName string) (string, error) { return "", errors.New("test error") }

	main()
}

func TestDoesRecordExistOnZoneFails(t *testing.T) {
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

	main()
}
