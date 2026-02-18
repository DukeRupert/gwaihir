# gwaihir — Project Scope

## Overview

`gwaihir` is a lightweight CLI tool for managing DNS records via the Cloudflare API. It is designed to be used standalone or as a component of the `arnor` infrastructure management suite.

Named after the Lord of the Eagles — carries messages swiftly across vast distances, appears when needed most.

## Repository

`github.com/fireflysoftware/gwaihir`

## Status

**Planned**

## Technology

- **Language:** Go
- **CLI framework:** Cobra
- **Key dependencies:** `github.com/spf13/cobra`, `github.com/joho/godotenv`

## Configuration

Credentials are loaded from `~/.dotfiles/.env` with fallback to `.env` in the current directory.

Required environment variables:

```env
CLOUDFLARE_API_TOKEN=your_api_token
CLOUDFLARE_ACCOUNT_ID=your_account_id
```

> Use a scoped API token with `Zone:DNS:Edit` and `Zone:Zone:Read` permissions. Avoid using the global API key.

## Project Structure

```
gwaihir/
├── main.go
├── cmd/
│   ├── root.go       # env loading, cobra setup, client init
│   └── dns.go        # all dns subcommands
├── internal/
│   └── cloudflare/
│       └── client.go # Cloudflare API client
├── .env.example
├── ROADMAP.md
└── README.md
```

## Commands

### v1.0.0

| Command | Description |
|---|---|
| `ping` | Verify credentials and return account info |
| `dns create` | Create a DNS record |
| `dns list` | List all records for a zone |
| `dns delete` | Delete a record by ID |
| `dns edit` | Edit a record by ID |
| `zone list` | List all zones in the account |

### v1.1.0

| Command | Description |
|---|---|
| `dns list-by-type` | List records filtered by type |
| `dns delete-by-type` | Delete records by type and name |

### v2.0.0 — Quality of Life

- `--output json` flag for machine-readable output
- `--quiet` flag for scripting
- Shell autocompletion (bash, zsh, fish)
- Config file support via Viper (`~/.config/gwaihir/config.yaml`)

### Flags — `dns create`

| Flag | Required | Default | Description |
|---|---|---|---|
| `--zone` | yes | — | Zone name (e.g. `example.com`) |
| `--name` | yes | — | Full record name (e.g. `myapp.example.com`) |
| `--type` | yes | — | Record type (`A`, `CNAME`, `MX`, etc.) |
| `--content` | yes | — | Record value (IP or target domain) |
| `--ttl` | no | `600` | TTL in seconds (`1` = auto in Cloudflare) |
| `--proxied` | no | `false` | Whether to proxy through Cloudflare |

## Key Differences from shadowfax

Cloudflare's API has some important differences from Porkbun's that the implementing agent should be aware of:

- Cloudflare uses **zones** rather than domains. A zone has an ID that must be looked up before creating records
- Record `name` is the **full** record name including the root domain (e.g. `myapp.example.com`), not just the subdomain
- TTL of `1` means "automatic" in Cloudflare
- The `--proxied` flag controls whether traffic routes through Cloudflare's CDN/proxy — important for business sites
- Authentication uses a Bearer token in the `Authorization` header, not a request body field
- Base URL: `https://api.cloudflare.com/client/v4`
- Responses wrap results in a `result` field and include a `success` boolean

## API Reference

Key endpoints:

```
GET  /zones?name={domain}           # look up zone ID by domain name
GET  /zones/{zone_id}/dns_records   # list records
POST /zones/{zone_id}/dns_records   # create record
PUT  /zones/{zone_id}/dns_records/{id}  # edit record
DELETE /zones/{zone_id}/dns_records/{id} # delete record
```

## Integration

When used as part of `arnor`, this package is imported as an internal client library. The standalone CLI remains available for direct use.

### Example usage in scripts

```bash
# Create A record
gwaihir dns create --zone example.com --name myapp.example.com --type A --content 1.2.3.4

# Create www CNAME
gwaihir dns create --zone example.com --name www.example.com --type CNAME --content example.com

# List all records for a zone
gwaihir dns list --zone example.com

# List all zones
gwaihir zone list
```

## Notes for Implementing Agents

- Zone ID must be resolved before any DNS operation — implement a `getZoneID(domain)` helper in the client that looks up the zone by name and caches it for the duration of the command
- Cloudflare returns paginated results — implement pagination handling in `ListRecords`
- The `proxied` field defaults to `false` for A and CNAME records but must be explicitly set in the request body
- Scoped API tokens are preferred over global API keys — the `.env.example` should reflect this
- Mirror the structure and conventions of `shadowfax` as closely as possible for consistency across the suite
- The `ping` command should call `GET /user/tokens/verify` to confirm the token is valid
