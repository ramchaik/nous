package llmclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nous/internal/cache"
	"strings"
	"time"

	"github.com/texttheater/golang-levenshtein/levenshtein"
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

type CachedResponse struct {
	Question string          `json:"question"`
	Response PredictResponse `json:"response"`
}

type QuestionIndex struct {
	Hash     string `json:"hash"`
	Question string `json:"question"`
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
	normalizedQuestion := strings.ToLower(strings.TrimSpace(question))
	questionHash := c.generateQuestionHash(normalizedQuestion)

	// Try to get from cache or find similar question
	cachedResponse, err := c.findCachedResponse(ctx, questionHash, normalizedQuestion)
	if err == nil {
		return &cachedResponse.Response, nil
	}

	// No cached response found, proceed with API call
	predictResp, err := c.makeAPICall(ctx, question)
	if err != nil {
		return nil, err
	}

	// Cache the response
	cachedResponse = &CachedResponse{
		Question: normalizedQuestion,
		Response: *predictResp,
	}
	if err := c.cacheResponse(ctx, questionHash, cachedResponse); err != nil {
		fmt.Printf("Error caching response: %v\n", err)
	}

	return predictResp, nil
}

func (c *Client) findCachedResponse(ctx context.Context, questionHash, normalizedQuestion string) (*CachedResponse, error) {
	// First, try to get the exact match
	cachedResp, err := c.getCachedResponse(ctx, questionHash)
	if err == nil {
		return cachedResp, nil
	}

	// If no exact match, try to find similar questions
	similarResp, err := c.findSimilarQuestion(ctx, normalizedQuestion)
	if err == nil {
		return similarResp, nil
	}

	// If no similar questions found, check uncompressed values
	return c.findSimilarUncompressedResponse(ctx, normalizedQuestion)
}

func (c *Client) generateQuestionHash(question string) string {
	hash := sha256.Sum256([]byte(question))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (c *Client) findSimilarQuestion(ctx context.Context, question string) (*CachedResponse, error) {
	allQuestions, err := c.getAllQuestions(ctx)
	if err != nil {
		return nil, err
	}

	for _, indexedQuestion := range allQuestions {
		similarity := calculateSimilarity(question, indexedQuestion.Question)
		if similarity >= 0.9 { // 90% similarity threshold
			return c.getCachedResponse(ctx, indexedQuestion.Hash)
		}
	}

	return nil, fmt.Errorf("no similar questions found")
}

func (c *Client) findSimilarUncompressedResponse(ctx context.Context, question string) (*CachedResponse, error) {
	allResponses, err := c.getAllUncompressedResponses(ctx)
	if err != nil {
		return nil, err
	}

	for _, cachedResp := range allResponses {
		similarity := calculateSimilarity(question, cachedResp.Question)
		if similarity >= 0.95 { // 95% similarity threshold
			return &cachedResp, nil
		}
	}

	return nil, fmt.Errorf("no similar uncompressed responses found")
}

func (c *Client) getAllQuestions(ctx context.Context) ([]QuestionIndex, error) {
	indexData, err := c.Cache.Get(ctx, "question_index")
	if err != nil {
		return nil, err
	}

	var questionIndex []QuestionIndex
	err = json.Unmarshal(indexData, &questionIndex)
	if err != nil {
		return nil, err
	}

	return questionIndex, nil
}

func (c *Client) getAllUncompressedResponses(ctx context.Context) ([]CachedResponse, error) {
	pattern := "predict:*"
	cachedData, err := c.Cache.GetAllValues(ctx, pattern)
	if err != nil {
		return nil, err
	}

	var allResponses []CachedResponse
	for _, data := range cachedData {
		var cachedResp CachedResponse
		err = json.Unmarshal(data, &cachedResp)
		if err == nil {
			allResponses = append(allResponses, cachedResp)
		}
	}

	return allResponses, nil
}

func (c *Client) getCachedResponse(ctx context.Context, questionHash string) (*CachedResponse, error) {
	cacheKey := fmt.Sprintf("predict:%s", questionHash)
	cachedData, err := c.Cache.GetUncompressed(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var cachedResp CachedResponse
	err = json.Unmarshal(cachedData, &cachedResp)
	if err != nil {
		return nil, err
	}

	return &cachedResp, nil
}

func (c *Client) cacheResponse(ctx context.Context, questionHash string, cachedResp *CachedResponse) error {
	// Cache the response
	jsonResponse, err := json.Marshal(cachedResp)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("predict:%s", questionHash)
	err = c.Cache.SetCompressed(ctx, cacheKey, jsonResponse, time.Hour) // Cache for 1 hour
	if err != nil {
		return err
	}

	// Update the question index
	return c.updateQuestionIndex(ctx, questionHash, cachedResp.Question)
}

func (c *Client) updateQuestionIndex(ctx context.Context, questionHash, question string) error {
	indexData, err := c.Cache.Get(ctx, "question_index")
	var questionIndex []QuestionIndex
	if err == nil {
		err = json.Unmarshal(indexData, &questionIndex)
		if err != nil {
			return err
		}
	}

	// Add new question to index
	questionIndex = append(questionIndex, QuestionIndex{Hash: questionHash, Question: question})

	// Save updated index
	updatedIndexData, err := json.Marshal(questionIndex)
	if err != nil {
		return err
	}

	return c.Cache.Set(ctx, "question_index", updatedIndexData, time.Hour*24*30) // Cache for 30 days
}

func (c *Client) makeAPICall(ctx context.Context, question string) (*PredictResponse, error) {
	reqBody := PredictRequest{Question: question}
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

	return &predictResp, nil
}

func calculateSimilarity(s1, s2 string) float64 {
	distance := levenshtein.DistanceForStrings([]rune(s1), []rune(s2), levenshtein.DefaultOptions)
	maxLen := float64(max(len(s1), len(s2)))
	return 1 - float64(distance)/maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
