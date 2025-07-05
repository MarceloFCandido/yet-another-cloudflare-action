package client

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/option"
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
