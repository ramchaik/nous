package llmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nous/internal/cache"
	"time"
)

type LLMClient interface {
	Predict(ctx context.Context, question string) (*PredictResponse, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL    string
	HTTPClient HTTPClient
	Cache      cache.Cacher
}

type PredictRequest struct {
	Question string `json:"question"`
}

type PredictResponse struct {
	Response string   `json:"response"`
	Steps    []string `json:"steps"`
}

func NewClient(baseURL string, httpClient HTTPClient, cache cache.Cacher) LLMClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		Cache:      cache,
	}
}

func (c *Client) Predict(ctx context.Context, question string) (*PredictResponse, error) {
	// Hash the cache key
	cacheKey := c.Cache.HashKey(fmt.Sprintf("predict:%s", question))

	// Try to get from cache
	cachedResponse, err := c.Cache.GetCompressed(ctx, cacheKey)
	if err == nil {
		var predictResp PredictResponse
		err = json.Unmarshal(cachedResponse, &predictResp)
		if err == nil {
			return &predictResp, nil
		}
	}

	// Cache miss, proceed with API call
	reqBody := PredictRequest{
		Question: question,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/predict", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var predictResp PredictResponse
	err = json.Unmarshal(body, &predictResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	// Cache the response
	jsonResponse, err := json.Marshal(predictResp)
	if err == nil {
		c.Cache.SetCompressed(ctx, cacheKey, jsonResponse, time.Hour) // Cache for 1 hour
	}

	return &predictResp, nil
}
