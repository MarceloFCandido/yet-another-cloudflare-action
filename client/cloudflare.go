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

func DoesRecordExistOnZone(zoneID, recordName string) (bool, error) {
	fmt.Printf("Checking if record '%s' exists on zone '%s'\n", recordName, zoneID)

	client := GetSingletonClient()

	page, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return false, fmt.Errorf("failed to list DNS records: %w", err)
	}

	for _, record := range page.Result {
		if record.Name == recordName {
			return true, nil
		}
	}

	return false, nil
}

func CreateRecordOnZone(zoneID string, record models.Record) (bool, error) {
	fmt.Printf("Creating record '%s' on zone '%s' of type %s\n", record.Record, zoneID, record.Type)

	client := GetSingletonClient()

	var body dns.RecordNewParamsBody

	switch record.Type {
	case "A":
		body = dns.RecordNewParamsBody{
			Name:    cloudflare.F(record.Record),
			Content: cloudflare.F(record.Target),
			Proxied: cloudflare.F(record.Proxy),
			TTL:     cloudflare.F(dns.TTL(record.Ttl)),
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeA),
		}
	case "CNAME":
		body = dns.RecordNewParamsBody{
			Name:    cloudflare.F(record.Record),
			Content: cloudflare.F(record.Target),
			Proxied: cloudflare.F(record.Proxy),
			TTL:     cloudflare.F(dns.TTL(record.Ttl)),
			Type:    cloudflare.F(dns.RecordNewParamsBodyTypeCNAME),
		}
	default:
		return false, fmt.Errorf("unsupported record type: %s", record.Type)
	}

	_, err := client.DNS.Records.New(context.TODO(), dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body:   body,
	})
	if err != nil {
		return false, fmt.Errorf("failed to create DNS record: %w", err)
	}

	return true, nil
}
