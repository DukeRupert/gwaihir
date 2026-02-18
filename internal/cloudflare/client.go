package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.cloudflare.com/client/v4"

type Client struct {
	Token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		httpClient: &http.Client{},
	}
}

// cfError represents an error returned by the Cloudflare API.
type cfError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// cfResultInfo contains pagination info from the Cloudflare API.
type cfResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

// cfResponse is the base response wrapper for all Cloudflare API calls.
type cfResponse struct {
	Success    bool            `json:"success"`
	Errors     []cfError       `json:"errors"`
	ResultInfo *cfResultInfo   `json:"result_info,omitempty"`
	Result     json.RawMessage `json:"result"`
}

// do performs an HTTP request against the Cloudflare API.
// method is the HTTP method, path is appended to baseURL,
// payload is marshaled as JSON for the request body (nil for no body),
// and out receives the unmarshaled "result" field from the response.
func (c *Client) do(method, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var cfResp cfResponse
	if err := json.Unmarshal(raw, &cfResp); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return fmt.Errorf("cloudflare error: %s (code %d)", cfResp.Errors[0].Message, cfResp.Errors[0].Code)
		}
		return fmt.Errorf("cloudflare request failed")
	}

	if out != nil && cfResp.Result != nil {
		if err := json.Unmarshal(cfResp.Result, out); err != nil {
			return fmt.Errorf("unmarshaling result: %w", err)
		}
	}

	return nil
}

// doPaginated performs paginated GET requests, collecting all results.
// path should not include query parameters — page and per_page are added automatically.
func (c *Client) doPaginated(path string, out any) error {
	var body io.Reader

	var allResults []json.RawMessage

	page := 1
	for {
		url := fmt.Sprintf("%s%s?page=%d&per_page=50", baseURL, path, page)

		req, err := http.NewRequest("GET", url, body)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("http request: %w", err)
		}

		raw, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}

		var cfResp cfResponse
		if err := json.Unmarshal(raw, &cfResp); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}

		if !cfResp.Success {
			if len(cfResp.Errors) > 0 {
				return fmt.Errorf("cloudflare error: %s (code %d)", cfResp.Errors[0].Message, cfResp.Errors[0].Code)
			}
			return fmt.Errorf("cloudflare request failed")
		}

		// Parse the result array and append individual items
		var items []json.RawMessage
		if err := json.Unmarshal(cfResp.Result, &items); err != nil {
			return fmt.Errorf("unmarshaling result array: %w", err)
		}
		allResults = append(allResults, items...)

		// Check if there are more pages
		if cfResp.ResultInfo == nil || page >= cfResp.ResultInfo.TotalPages {
			break
		}
		page++
	}

	// Marshal all collected results back and unmarshal into the output type
	collected, err := json.Marshal(allResults)
	if err != nil {
		return fmt.Errorf("marshaling collected results: %w", err)
	}

	if err := json.Unmarshal(collected, out); err != nil {
		return fmt.Errorf("unmarshaling collected results: %w", err)
	}

	return nil
}

// Zone represents a Cloudflare DNS zone.
type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ListZones returns all zones in the account.
func (c *Client) ListZones() ([]Zone, error) {
	var zones []Zone
	if err := c.doPaginated("/zones", &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

// DNSRecord represents a Cloudflare DNS record.
type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// GetZoneID resolves a domain name to its Cloudflare zone ID.
func (c *Client) GetZoneID(domain string) (string, error) {
	var zones []Zone
	if err := c.do("GET", "/zones?name="+domain, nil, &zones); err != nil {
		return "", err
	}
	if len(zones) == 0 {
		return "", fmt.Errorf("no zone found for domain: %s", domain)
	}
	return zones[0].ID, nil
}

// ListRecords returns all DNS records for the given zone ID.
func (c *Client) ListRecords(zoneID string) ([]DNSRecord, error) {
	var records []DNSRecord
	if err := c.doPaginated("/zones/"+zoneID+"/dns_records", &records); err != nil {
		return nil, err
	}
	return records, nil
}

// CreateRecord creates a DNS record in the given zone.
func (c *Client) CreateRecord(zoneID string, record DNSRecord) (DNSRecord, error) {
	var result DNSRecord
	if err := c.do("POST", "/zones/"+zoneID+"/dns_records", record, &result); err != nil {
		return DNSRecord{}, err
	}
	return result, nil
}

// EditRecord updates a DNS record by ID in the given zone.
func (c *Client) EditRecord(zoneID, recordID string, record DNSRecord) (DNSRecord, error) {
	var result DNSRecord
	if err := c.do("PUT", "/zones/"+zoneID+"/dns_records/"+recordID, record, &result); err != nil {
		return DNSRecord{}, err
	}
	return result, nil
}

// VerifyToken verifies the API token is valid.
func (c *Client) VerifyToken() error {
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	if err := c.do("GET", "/user/tokens/verify", nil, &result); err != nil {
		return err
	}

	if result.Status != "active" {
		return fmt.Errorf("token status: %s", result.Status)
	}

	return nil
}
