package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"yaca/models"
	"yaca/pkg/logger"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

var once sync.Once
var client *cloudflare.Client

var GetSingletonClient = getSingletonClient

func getSingletonClient() *cloudflare.Client {
	if client == nil {
		once.Do(func() {
			logger.Debug("Initializing singleton Cloudflare client")

			apiEmail := os.Getenv("CLOUDFLARE_API_EMAIL")
			apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
			
			// Log that we're using credentials without exposing them
			logger.Debug("Creating Cloudflare client",
				slog.Bool("has_email", apiEmail != ""),
				slog.Bool("has_token", apiToken != ""))

			client = cloudflare.NewClient(
				option.WithAPIEmail(apiEmail),
				option.WithAPIToken(apiToken),
			)
		})
	} else {
		logger.Debug("Using existing Cloudflare client instance")
	}

	return client
}

var GetZoneIDByName = getZoneIDByName

func getZoneIDByName(zoneName string) (string, error) {
	logger.Debug("Retrieving zone ID",
		slog.String("zone_name", zoneName))

	client := GetSingletonClient()

	page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Name: cloudflare.F(zoneName),
	})
	if err != nil {
		logger.Error("Failed to list zones",
			slog.String("zone_name", zoneName),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	if len(page.Result) == 0 {
		logger.Warn("No zone found",
			slog.String("zone_name", zoneName))
		return "", fmt.Errorf("no zone found with name: %s", zoneName)
	}

	zoneID := page.Result[0].ID
	logger.Debug("Zone ID retrieved",
		slog.String("zone_id", zoneID), // Will be masked
		slog.String("zone_name", zoneName))

	return zoneID, nil
}

var DoesRecordExistOnZone = doesRecordExistOnZone

func doesRecordExistOnZone(zoneID, recordName string) (string, error) {
	logger.Debug("Checking record existence",
		slog.String("zone_id", zoneID), // Will be masked
		slog.String("record_name", recordName))

	client := GetSingletonClient()

	page, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		logger.Error("Failed to list DNS records",
			slog.String("zone_id", zoneID),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to list DNS records: %w", err)
	}

	for _, record := range page.Result {
		if record.Name == recordName {
			logger.Debug("Record found",
				slog.String("record_id", record.ID), // Will be masked
				slog.String("record_name", recordName))
			return record.ID, nil
		}
	}

	logger.Debug("Record not found",
		slog.String("record_name", recordName))
	return "", nil
}

var CreateRecordOnZone = createRecordOnZone

func createRecordOnZone(zoneID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID: zoneID,
		Record: record,
	}, "Creating")
}

var UpdateRecordOnZone = updateRecordOnZone

func updateRecordOnZone(zoneID, recordID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID:   zoneID,
		RecordID: recordID,
		Record:   record,
	}, "Updating")
}

var DeleteRecordOnZone = deleteRecordOnZone

func deleteRecordOnZone(zoneID, recordID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID:   zoneID,
		RecordID: recordID,
		Record:   record,
	}, "Deleting")
}

var handleRecord = handleRecordImpl

func handleRecordImpl(ctx context.Context, recordData models.RecordData, operation string) (bool, error) {
	// Log operation with appropriate details
	logger.Info("DNS operation started",
		slog.String("operation", operation),
		slog.String("record_name", recordData.Record.Record),
		slog.String("record_type", recordData.Record.Type),
		slog.String("zone_id", recordData.ZoneID), // Will be masked
		slog.Bool("proxied", recordData.Record.Proxy),
		slog.Float64("ttl", recordData.Record.Ttl))
	
	// Log target only for non-delete operations and mask if it's an IP
	if operation != "Deleting" && recordData.Record.Target != "" {
		logger.Debug("Record target",
			slog.String("target_ip", recordData.Record.Target),
			slog.String("operation", operation))
	}

	client := GetSingletonClient()
	var err error

	switch operation {
	case "Creating":
		body := dns.RecordNewParamsBody{
			Name:    cloudflare.F(recordData.Record.Record),
			Content: cloudflare.F(recordData.Record.Target),
			Proxied: cloudflare.F(recordData.Record.Proxy),
			TTL:     cloudflare.F(dns.TTL(recordData.Record.Ttl)),
		}
		switch recordData.Record.Type {
		case "A":
			body.Type = cloudflare.F(dns.RecordNewParamsBodyTypeA)
		case "CNAME":
			body.Type = cloudflare.F(dns.RecordNewParamsBodyTypeCNAME)
		default:
			logger.Error("Unsupported record type",
				slog.String("type", recordData.Record.Type))
			return false, fmt.Errorf("unsupported record type: %s", recordData.Record.Type)
		}
		_, err = client.DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(recordData.ZoneID),
			Body:   body,
		})
	case "Updating":
		body := dns.RecordEditParamsBody{
			Name:    cloudflare.F(recordData.Record.Record),
			Content: cloudflare.F(recordData.Record.Target),
			Proxied: cloudflare.F(recordData.Record.Proxy),
			TTL:     cloudflare.F(dns.TTL(recordData.Record.Ttl)),
		}
		switch recordData.Record.Type {
		case "A":
			body.Type = cloudflare.F(dns.RecordEditParamsBodyTypeA)
		case "CNAME":
			body.Type = cloudflare.F(dns.RecordEditParamsBodyTypeCNAME)
		default:
			logger.Error("Unsupported record type",
				slog.String("type", recordData.Record.Type))
			return false, fmt.Errorf("unsupported record type: %s", recordData.Record.Type)
		}
		_, err = client.DNS.Records.Edit(ctx, recordData.RecordID, dns.RecordEditParams{
			ZoneID: cloudflare.F(recordData.ZoneID),
			Body:   body,
		})
	case "Deleting":
		_, err = client.DNS.Records.Delete(ctx, recordData.RecordID, dns.RecordDeleteParams{
			ZoneID: cloudflare.F(recordData.ZoneID),
		})
	default:
		logger.Error("Unsupported operation",
			slog.String("operation", operation))
		return false, fmt.Errorf("unsupported operation: %s", operation)
	}

	if err != nil {
		logger.Error("DNS operation failed",
			slog.String("operation", operation),
			slog.String("record_name", recordData.Record.Record),
			slog.String("error", err.Error()))
		return false, fmt.Errorf("failed to %s DNS record: %w", map[string]string{
			"Creating": "create",
			"Updating": "update",
			"Deleting": "delete",
		}[operation], err)
	}

	logger.Info("DNS operation completed successfully",
		slog.String("operation", operation),
		slog.String("record_name", recordData.Record.Record),
		slog.String("record_type", recordData.Record.Type))

	return true, nil
}
