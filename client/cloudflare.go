package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cloudflare/cloudflare-go/v3"
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
