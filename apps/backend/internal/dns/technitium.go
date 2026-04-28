package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

// TechnitiumClient communicates with the Technitium DNS Server API.
type TechnitiumClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// Stats represents DNS server statistics.
type Stats struct {
	TotalQueries       int64   `json:"total_queries"`
	TotalBlockedQueries int64  `json:"total_blocked_queries"`
	BlockedPercent     float64 `json:"blocked_percent"`
	TotalClients       int     `json:"total_clients"`
	QueryTypes         map[string]int64 `json:"query_types"`
	TopDomains         []DomainStat     `json:"top_domains"`
	TopBlockedDomains  []DomainStat     `json:"top_blocked_domains"`
	TopClients         []ClientStat     `json:"top_clients"`
}

// DomainStat represents a top domain entry.
type DomainStat struct {
	Domain string `json:"domain"`
	Hits   int64  `json:"hits"`
}

// ClientStat represents a top client entry.
type ClientStat struct {
	Client  string `json:"client"`
	Queries int64  `json:"queries"`
}

// QueryLogEntry represents a DNS query log entry.
type QueryLogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	ClientIP     string    `json:"client_ip"`
	Domain       string    `json:"domain"`
	Type         string    `json:"type"`
	ResponseCode string   `json:"response_code"`
	Blocked      bool      `json:"blocked"`
	ResponseTime int       `json:"response_time_ms"`
}

// BlockList represents a DNS blocklist.
type BlockList struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
	Entries int    `json:"entries"`
}

// NewTechnitiumClient creates a DNS client.
func NewTechnitiumClient(apiURL, token string) *TechnitiumClient {
	if apiURL == "" {
		return nil
	}
	return &TechnitiumClient{
		baseURL:  apiURL,
		apiToken: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsAvailable checks if the DNS server is reachable.
func (c *TechnitiumClient) IsAvailable() bool {
	if c == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := c.doRequest(ctx, "/api/dashboard/stats/get", nil)
	return err == nil
}

// GetStats returns DNS dashboard statistics.
func (c *TechnitiumClient) GetStats(ctx context.Context) (*Stats, error) {
	if c == nil {
		return nil, fmt.Errorf("DNS client not available")
	}

	resp, err := c.doRequest(ctx, "/api/dashboard/stats/get", url.Values{
		"type": {"lastDay"},
	})
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		QueryTypes: make(map[string]int64),
	}

	if response, ok := resp["response"].(map[string]interface{}); ok {
		if v, ok := response["totalQueries"].(float64); ok {
			stats.TotalQueries = int64(v)
		}
		if v, ok := response["totalBlockedQueries"].(float64); ok {
			stats.TotalBlockedQueries = int64(v)
		}
		if stats.TotalQueries > 0 {
			stats.BlockedPercent = float64(stats.TotalBlockedQueries) / float64(stats.TotalQueries) * 100
		}
		if v, ok := response["totalClients"].(float64); ok {
			stats.TotalClients = int(v)
		}
	}

	return stats, nil
}

// GetQueryLog retrieves recent DNS query logs.
func (c *TechnitiumClient) GetQueryLog(ctx context.Context, pageNumber int, entriesPerPage int) ([]QueryLogEntry, error) {
	if c == nil {
		return nil, fmt.Errorf("DNS client not available")
	}

	resp, err := c.doRequest(ctx, "/api/queryLogs/list", url.Values{
		"pageNumber":     {fmt.Sprintf("%d", pageNumber)},
		"entriesPerPage": {fmt.Sprintf("%d", entriesPerPage)},
	})
	if err != nil {
		return nil, err
	}

	var entries []QueryLogEntry
	if response, ok := resp["response"].(map[string]interface{}); ok {
		if queryEntries, ok := response["entries"].([]interface{}); ok {
			for _, entry := range queryEntries {
				if e, ok := entry.(map[string]interface{}); ok {
					qe := QueryLogEntry{
						ClientIP: fmt.Sprintf("%v", e["clientIpAddress"]),
						Domain:   fmt.Sprintf("%v", e["question"]),
						Type:     fmt.Sprintf("%v", e["type"]),
					}
					entries = append(entries, qe)
				}
			}
		}
	}

	return entries, nil
}

// GetBlockLists returns configured blocklists.
func (c *TechnitiumClient) GetBlockLists(ctx context.Context) ([]BlockList, error) {
	if c == nil {
		return nil, fmt.Errorf("DNS client not available")
	}

	resp, err := c.doRequest(ctx, "/api/settings/get", nil)
	if err != nil {
		return nil, err
	}

	var blockLists []BlockList
	if response, ok := resp["response"].(map[string]interface{}); ok {
		if lists, ok := response["blockListUrls"].([]interface{}); ok {
			for _, list := range lists {
				if u, ok := list.(string); ok {
					blockLists = append(blockLists, BlockList{
						URL:     u,
						Enabled: true,
					})
				}
			}
		}
	}

	return blockLists, nil
}

func (c *TechnitiumClient) doRequest(ctx context.Context, path string, params url.Values) (map[string]interface{}, error) {
	if params == nil {
		params = url.Values{}
	}
	if c.apiToken != "" {
		params.Set("token", c.apiToken)
	}

	reqURL := c.baseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DNS API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if status, ok := result["status"].(string); ok && status == "error" {
		return nil, fmt.Errorf("DNS API error: %v", result["errorMessage"])
	}

	log.Debug().Str("path", path).Msg("DNS API request")
	return result, nil
}
