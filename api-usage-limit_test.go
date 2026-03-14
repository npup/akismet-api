package akismet

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeClientWithServer(t *testing.T, server *httptest.Server, apiKey string, blogURL string) *Client {
	t.Helper()
	client, err := newClientWithApiBaseURL(context.Background(), apiKey, blogURL, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func makeVerifyThenRespond(verifyResponse string, nextStatus int, nextResponse string) *httptest.Server {
	callCount := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// first call is always verify-key
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(verifyResponse))
			return
		}
		w.WriteHeader(nextStatus)
		w.Write([]byte(nextResponse))
	}))
}

func TestGetUsageLimit_Valid(t *testing.T) {
	responseBody := `{"limit":"1000","usage":423,"percentage":"42.3","throttled":false}`
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	result, err := client.GetUsageLimit(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Limit == nil || *result.Limit != 1000 {
		t.Errorf("expected limit 1000, got %v", result.Limit)
	}
	if result.Usage != 423 {
		t.Errorf("expected usage %d, got %d", 423, result.Usage)
	}
	if result.Percentage != 42.3 {
		t.Errorf("expected percentage %v, got %v", 42.3, result.Percentage)
	}
	if result.Throttled {
		t.Error("expected throttled false, got true")
	}
}

func TestGetUsageLimit_Unlimited(t *testing.T) {
	responseBody := fmt.Sprintf(`{"limit":"%s","usage":0,"percentage":"0","throttled":false}`, propUsageLimitNoLimit)
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	result, err := client.GetUsageLimit(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Limit != nil {
		t.Errorf("expected nil limit for unlimited plan, got %v", *result.Limit)
	}
}

func TestGetUsageLimit_Throttled(t *testing.T) {
	responseBody := `{"limit":"1000","usage":1100,"percentage":"110.0","throttled":true}`
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()

	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	result, err := client.GetUsageLimit(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Limit == nil || *result.Limit != 1000 {
		t.Errorf("expected limit 1000, got %v", result.Limit)
	}
	if !result.Throttled {
		t.Error("expected throttled true, got false")
	}
}

func TestGetUsageLimit_InvalidKey(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, bodyInvalidMessage)
	defer server.Close()
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	_, err := client.GetUsageLimit(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestGetUsageLimit_MalformedResponse(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, "this is not json")
	defer server.Close()
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	_, err := client.GetUsageLimit(context.Background())
	if err == nil {
		t.Fatal("expected error for malformed response, got nil")
	}
}

func TestGetUsageLimit_ServerDown(t *testing.T) {
	server := makeServer(http.StatusOK, bodyValidMessage)
	serverURL := server.URL
	client, clientErr := newClientWithApiBaseURL(context.Background(), "test-key", "http://example.com", serverURL)
	if clientErr != nil {
		t.Fatalf("failed to create client: %v", clientErr)
	}
	server.Close()

	_, err := client.GetUsageLimit(context.Background())
	if err == nil {
		t.Fatal("expected error when server is down, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestGetUsageLimit_RequestParams(t *testing.T) {
	var gotAPIKey string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		gotAPIKey = r.URL.Query().Get("api_key")
		w.Write([]byte(`{"limit":"1000","usage":0,"percentage":"0","throttled":false}`))
	}))
	defer server.Close()
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)

	client.GetUsageLimit(context.Background())

	if gotAPIKey != myApiKey {
		t.Errorf("expected api_key %q, got %q", myApiKey, gotAPIKey)
	}
}
