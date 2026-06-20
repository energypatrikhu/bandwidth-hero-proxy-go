package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type flareSolverrRequest struct {
	Cmd        string `json:"cmd"`
	URL        string `json:"url"`
	MaxTimeout int    `json:"maxTimeout"`
}

type flareSolverrCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain,omitempty"`
}

type flareSolverrSolution struct {
	URL       string               `json:"url"`
	Status    int                  `json:"status"`
	Cookies   []flareSolverrCookie `json:"cookies"`
	UserAgent string               `json:"userAgent"`
	Response  string               `json:"response"`
}

type flareSolverrResponse struct {
	Status   string               `json:"status"`
	Message  string               `json:"message"`
	Solution flareSolverrSolution `json:"solution"`
}

// SolveWithFlareSolverr asks the FlareSolverr instance (BHP_FLARESOLVERR_URL)
// to load targetURL in a real browser, solving any Cloudflare/JS challenge
// along the way, and returns the resulting cookies + User-Agent.
func SolveWithFlareSolverr(targetURL string, timeout time.Duration) (*flareSolverrSolution, error) {
	if strings.TrimSpace(BHP_FLARESOLVERR_URL) == "" {
		return nil, fmt.Errorf("BHP_FLARESOLVERR_URL is not configured")
	}

	reqPayload := flareSolverrRequest{
		Cmd:        "request.get",
		URL:        targetURL,
		MaxTimeout: int(timeout / time.Millisecond),
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal flaresolverr request: %v", err)
	}

	endpoint := strings.TrimRight(BHP_FLARESOLVERR_URL, "/") + "/v1"

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build flaresolverr request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach flaresolverr at %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read flaresolverr response: %v", err)
	}

	var fsResp flareSolverrResponse
	if err := json.Unmarshal(body, &fsResp); err != nil {
		return nil, fmt.Errorf("failed to parse flaresolverr response: %v", err)
	}

	if fsResp.Status != "ok" {
		return nil, fmt.Errorf("flaresolverr returned error status %q: %s", fsResp.Status, fsResp.Message)
	}

	return &fsResp.Solution, nil
}
