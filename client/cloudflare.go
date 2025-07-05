package client

import (
	"os"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/option"
)

func GetClient() (*cloudflare.Client) {
  client := cloudflare.NewClient(
		option.WithAPIEmail(os.Getenv("CLOUDFLARE_API_EMAIL")),
		option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")),
	)

	return client
}
