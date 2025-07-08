package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"yaca/models"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

var once sync.Once
var client *cloudflare.Client

func GetSingletonClient() *cloudflare.Client {
	if client == nil {
		once.Do(func() {
			fmt.Println("Initializing singleton Cloudflare client...")

			client = cloudflare.NewClient(
				option.WithAPIEmail(os.Getenv("CLOUDFLARE_API_EMAIL")),
				option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")),
			)
		})
	} else {
		fmt.Println("Using existing Cloudflare client instance.")
	}

	return client
}

func GetZoneIDByName(zoneName string) (string, error) {
	fmt.Printf("Retrieving zone ID for zone name: %s\n", zoneName)

	client := GetSingletonClient()

	page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Name: cloudflare.F(zoneName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	if len(page.Result) == 0 {
		return "", fmt.Errorf("no zone found with name: %s", zoneName)
	}

	return page.Result[0].ID, nil
}

func DoesRecordExistOnZone(zoneID, recordName string) (string, error) {
	fmt.Printf("Checking if record '%s' exists on zone '%s'\n", recordName, zoneID)

	client := GetSingletonClient()

	page, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list DNS records: %w", err)
	}

	for _, record := range page.Result {
		if record.Name == recordName {
			return record.ID, nil
		}
	}

	return "", nil
}

func CreateRecordOnZone(zoneID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID: zoneID,
		Record: record,
	}, "Creating")
}

func UpdateRecordOnZone(zoneID, recordID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID:   zoneID,
		RecordID: recordID,
		Record:   record,
	}, "Updating")
}

func DeleteRecordOnZone(zoneID, recordID string, record models.Record) (bool, error) {
	return handleRecord(context.TODO(), models.RecordData{
		ZoneID:   zoneID,
		RecordID: recordID,
		Record:   record,
	}, "Deleting")
}

func handleRecord(ctx context.Context, recordData models.RecordData, operation string) (bool, error) {
	fmt.Printf("%s record '%s' on zone '%s' of type %s\n", operation, recordData.Record.Record, recordData.ZoneID, recordData.Record.Type)

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
		return false, fmt.Errorf("unsupported operation: %s", operation)
	}

	if err != nil {
		return false, fmt.Errorf("failed to %s DNS record: %w", map[string]string{
			"Creating": "create",
			"Updating": "update",
			"Deleting": "delete",
		}[operation], err)
	}

	return true, nil
}
