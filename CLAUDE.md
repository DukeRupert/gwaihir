# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**gwaihir** is a Go CLI tool for managing DNS records via the Cloudflare API. It can be used standalone or as a component of the `arnor` infrastructure management suite. Part of a family of tools that includes `shadowfax` (Porkbun DNS) — mirror its structure and conventions.

## Build & Run Commands

```bash
go build -o gwaihir .        # Build binary
go run .                      # Run directly
go test ./...                 # Run all tests
go test ./internal/cloudflare # Run tests for a specific package
go vet ./...                  # Lint
```

## Architecture

```
main.go              → Entry point
cmd/root.go          → Cobra setup, env loading, client initialization
cmd/dns.go           → DNS subcommands (create, list, delete, edit)
internal/cloudflare/
  client.go          → Cloudflare API client
```

**CLI framework:** Cobra (`github.com/spf13/cobra`)
**Env loading:** godotenv (`github.com/joho/godotenv`)

## Configuration

Credentials loaded from `~/.dotfiles/.env` with fallback to `.env` in the current directory.

Required env vars: `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID`

## Cloudflare API Notes

- Base URL: `https://api.cloudflare.com/client/v4`
- Auth: Bearer token in `Authorization` header
- Zone ID must be resolved by domain name before any DNS operation — implement `getZoneID(domain)` with per-command caching
- Record `name` is the **full** record name including root domain (e.g. `myapp.example.com`)
- TTL of `1` means "automatic" in Cloudflare
- `proxied` defaults to `false` but must be explicitly set in request body
- Responses wrap results in a `result` field with a `success` boolean
- Paginated results — implement pagination handling in `ListRecords`
- Token verification: `GET /user/tokens/verify`
