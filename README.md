# gwaihir

CLI tool for managing Cloudflare DNS records. Part of the `arnor` infrastructure management suite.

## Install

```bash
go install github.com/dukerupert/gwaihir@latest
```

Or build from source:

```bash
go build -o gwaihir .
```

## Configuration

Set credentials in `~/.dotfiles/.env` (preferred) or `.env` in the current directory:

```
CLOUDFLARE_API_TOKEN=your_api_token
CLOUDFLARE_ACCOUNT_ID=your_account_id
```

## Usage

```bash
# Verify token
gwaihir ping

# List zones
gwaihir zone list

# List DNS records
gwaihir dns list --zone example.com

# Create a record
gwaihir dns create --zone example.com --name app.example.com --type A --content 1.2.3.4
gwaihir dns create --zone example.com --name app.example.com --type CNAME --content target.com --proxied

# Edit a record
gwaihir dns edit --zone example.com --id <record-id> --name app.example.com --type A --content 5.6.7.8

# Delete a record
gwaihir dns delete --zone example.com --id <record-id>
```
