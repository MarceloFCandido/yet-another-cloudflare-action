package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/dns"
	"github.com/cloudflare/cloudflare-go/v3/option"
	"github.com/cloudflare/cloudflare-go/v3/zones"
)

var once sync.Once
var client *cloudflare.Client

func GetSingletonClient() (*cloudflare.Client) {
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
		Name: cloudflare.F(recordName),
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

