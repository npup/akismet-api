package akismet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetKeySites_Valid(t *testing.T) {
	responseBody := `{
		"http://some.example.com": {"api_calls":150,"spam":30,"ham":120,"missed_spam":2,"false_positives":1,"is_revoked":false},
		"http://other.example.com":   {"api_calls":50,"spam":5,"ham":45,"missed_spam":0,"false_positives":0,"is_revoked":false},
		"limit": 500,
		"offset": 0,
		"total": 2
	}`
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.GetKeySites(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
	if result.Limit != 500 {
		t.Errorf("expected limit 500, got %d", result.Limit)
	}
	if result.Offset != 0 {
		t.Errorf("expected offset 0, got %d", result.Offset)
	}
	if len(result.Sites) != 2 {
		t.Fatalf("expected 2 sites, got %d", len(result.Sites))
	}

	sitesByURL := make(map[string]KeySiteEntry)
	for _, s := range result.Sites {
		sitesByURL[s.Site] = s
	}

	example := sitesByURL["http://some.example.com"]
	if example.APICalls != 150 {
		t.Errorf("expected api_calls 150, got %d", example.APICalls)
	}
	if example.Spam != 30 {
		t.Errorf("expected spam 30, got %d", example.Spam)
	}
	if example.Ham != 120 {
		t.Errorf("expected ham 120, got %d", example.Ham)
	}
	if example.MissedSpam != 2 {
		t.Errorf("expected missed_spam 2, got %d", example.MissedSpam)
	}
	if example.FalsePositives != 1 {
		t.Errorf("expected false_positives 1, got %d", example.FalsePositives)
	}
	if example.IsRevoked {
		t.Error("expected is_revoked false, got true")
	}
}

func TestGetKeySites_IsRevoked(t *testing.T) {
	responseBody := `{
		"http://revoked.example.com": {"api_calls":0,"spam":0,"ham":0,"missed_spam":0,"false_positives":0,"is_revoked":true},
		"limit": 500,
		"offset": 0,
		"total": 1
	}`
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.GetKeySites(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result.Sites) != 1 {
		t.Fatalf("expected 1 site, got %d", len(result.Sites))
	}
	if !result.Sites[0].IsRevoked {
		t.Error("expected is_revoked true, got false")
	}
}

func TestGetKeySites_Pagination(t *testing.T) {
	responseBody := `{
		"http://example.com": {"api_calls":10,"spam":1,"ham":9,"missed_spam":0,"false_positives":0,"is_revoked":false},
		"limit": 10,
		"offset": 20,
		"total": 100
	}`
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, responseBody)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.GetKeySites(context.Background(), &KeySitesParams{Limit: 10, Offset: 20})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Limit != 10 {
		t.Errorf("expected limit 10, got %d", result.Limit)
	}
	if result.Offset != 20 {
		t.Errorf("expected offset 20, got %d", result.Offset)
	}
	if result.Total != 100 {
		t.Errorf("expected total 100, got %d", result.Total)
	}
}

func TestGetKeySites_InvalidKey(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, bodyInvalidMessage)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	_, err := client.GetKeySites(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestGetKeySites_MalformedResponse(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, "this is not json")
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	_, err := client.GetKeySites(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for malformed response, got nil")
	}
}

func TestGetKeySites_ServerDown(t *testing.T) {
	server := makeServer(http.StatusOK, bodyValidMessage)
	serverURL := server.URL
	client, clientErr := newClientWithApiBaseURL(context.Background(), "test-key", "http://example.com", serverURL)
	if clientErr != nil {
		t.Fatalf("failed to create client: %v", clientErr)
	}
	server.Close()

	_, err := client.GetKeySites(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when server is down, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestGetKeySites_QueryParams(t *testing.T) {
	var gotAPIKey, gotMonth, gotFilter, gotOrder, gotLimit, gotOffset string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		q := r.URL.Query()
		gotAPIKey = q.Get("api_key")
		gotMonth = q.Get("month")
		gotFilter = q.Get("filter")
		gotOrder = q.Get("order")
		gotLimit = q.Get("limit")
		gotOffset = q.Get("offset")
		w.Write([]byte(`{"limit":10,"offset":5,"total":0}`))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	client.GetKeySites(context.Background(), &KeySitesParams{
		Month:  "2026-01",
		Filter: "example",
		Order:  KeySitesOrderSpam,
		Limit:  10,
		Offset: 5,
	})

	if gotAPIKey != "my-apikey" {
		t.Errorf("expected api_key %q, got %q", "my-apikey", gotAPIKey)
	}
	if gotMonth != "2026-01" {
		t.Errorf("expected month %q, got %q", "2026-01", gotMonth)
	}
	if gotFilter != "example" {
		t.Errorf("expected filter %q, got %q", "example", gotFilter)
	}
	if gotOrder != "spam" {
		t.Errorf("expected order %q, got %q", "spam", gotOrder)
	}
	if gotLimit != "10" {
		t.Errorf("expected limit %q, got %q", "10", gotLimit)
	}
	if gotOffset != "5" {
		t.Errorf("expected offset %q, got %q", "5", gotOffset)
	}
}

func TestGetKeySites_NilParams_OnlySendsAPIKey(t *testing.T) {
	var gotQuery string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		gotQuery = r.URL.RawQuery
		w.Write([]byte(`{"limit":500,"offset":0,"total":0}`))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	client.GetKeySites(context.Background(), nil)

	// Only api_key should be present — no month, filter, order, limit, offset
	if gotQuery != "api_key=my-apikey" {
		t.Errorf("expected only api_key in query, got %q", gotQuery)
	}
}
