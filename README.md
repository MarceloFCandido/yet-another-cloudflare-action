# Yet Another Cloudflare Action (YACA)

This GitHub Action simplifies the management of Cloudflare DNS records directly from your workflows. It allows you to create, update, and delete DNS records within your Cloudflare zones.

## Features

- **Create DNS Records:** Automatically creates a new DNS record if it doesn't exist.
- **Update DNS Records:** Updates an existing DNS record with new information.
- **Delete DNS Records:** Explicitly deletes a specified DNS record.

## Usage

To use this action in your GitHub Actions workflow, you need to provide the necessary inputs for the DNS record you want to manage.

### Inputs

| Input       | Description                                            | Required | Default   |
|-------------|--------------------------------------------------------|----------|-----------|
| `record`    | The full record name (e.g., `www.example.com`).        | `true`   |           |
| `zone-name` | The Cloudflare zone name (e.g., `example.com`).        | `true`   |           |
| `delete`    | Set to `true` to delete the record.                    | `true`   | `false`   |
| `target`    | The target IP address or hostname for the record.      | `false`  |           |
| `type`      | The type of DNS record (e.g., `A`, `CNAME`).           | `false`  |           |
| `proxy`     | Whether to enable Cloudflare proxy for the record.     | `false`  | `false`   |
| `ttl`       | The Time-To-Live (TTL) for the record in seconds.      | `false`  | `3600`    |

### Examples

Here are examples demonstrating how to use the action for creating, updating, and deleting DNS records, similar to how the action is tested internally.

#### Cloudflare API Credentials

Before using the action, ensure you have set up your Cloudflare API credentials as GitHub Secrets.

- `CLOUDFLARE_API_EMAIL`: Your Cloudflare account email.
- `CLOUDFLARE_API_TOKEN`: Your Cloudflare API Token with sufficient permissions to edit DNS zones.

#### Create a DNS Record

This example creates a `CNAME` record pointing `your-record.example.com` to `www.bing.com`.

```yaml
- name: Create DNS Record
  uses: marcelofcandido/yet-another-cloudflare-action@master
  env:
    CLOUDFLARE_API_EMAIL: ${{ secrets.CLOUDFLARE_API_EMAIL }}
    CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
  with:
    record: your-record.example.com
    zone-name: your-zone.com
    target: www.bing.com
    type: CNAME
    proxy: true
    ttl: 300
```

#### Update a DNS Record

This example updates the `CNAME` record `your-record.example.com` to point to `www.google.com`.

```yaml
- name: Update DNS Record
  uses: marcelofcandido/yet-another-cloudflare-action@master
  env:
    CLOUDFLARE_API_EMAIL: ${{ secrets.CLOUDFLARE_API_EMAIL }}
    CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
  with:
    record: your-record.example.com
    zone-name: your-zone.com
    target: www.google.com
    type: CNAME
    proxy: true
    ttl: 300
```

#### Delete a DNS Record

This example deletes the `your-record.example.com` DNS record.

```yaml
- name: Delete DNS Record
  uses: marcelofcandido/yet-another-cloudflare-action@master
  env:
    CLOUDFLARE_API_EMAIL: ${{ secrets.CLOUDFLARE_API_EMAIL }}
    CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
  with:
    record: your-record.example.com
    zone-name: your-zone.com
    delete: true
```