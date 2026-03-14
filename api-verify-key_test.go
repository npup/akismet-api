package akismet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeServer(httpStatus int, responseText string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		w.Write([]byte(responseText))
	}))
	return server
}
func TestVerifyKey_Valid(t *testing.T) {
	server := makeServer(http.StatusOK, bodyValidMessage)
	defer server.Close()

	serverURL := server.URL

	client, err := newClientWithApiBaseURL(context.Background(), "test-key", "http://example.com", serverURL)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := apiEndpoints{
		VerifyKey:    serverURL + "/1.1/verify-key",
		CheckComment: serverURL + "/1.1/comment-check",
		SubmitSpam:   serverURL + "/1.1/submit-spam",
		SubmitHam:    serverURL + "/1.1/submit-ham",
		KeySites:     serverURL + "/1.2/key-sites",
		UsageLimit:   serverURL + "/1.2/usage-limit",
	}
	if client.endpoints != expected {
		t.Errorf("endpoints mismatch:\ngot:  %+v\nwant: %+v", client.endpoints, expected)
	}
}

func TestVerifyKey_ServerDown(t *testing.T) {
	server := makeServer(http.StatusOK, bodyValidMessage)
	server.Close()

	_, err := newClientWithApiBaseURL(context.Background(), "test-key", "http://example.com", server.URL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Alert != nil {
		t.Fatal("expected error to NOT have any akismet alert")
	}
	if err.DebugHelp != "" {
		t.Fatal("expected error to NOT have any akismet debughelp message")
	}
	if err.Err == nil {
		t.Fatal("expected inner error")
	}
	if err.Err.Error() == "" {
		t.Fatal("expected inner error to have text")
	}
}

func TestVerifyKey_Invalid(t *testing.T) {
	server := makeServer(http.StatusOK, bodyInvalidMessage)
	defer server.Close()

	_, err := newClientWithApiBaseURL(context.Background(), "bad-key", "http://example.com", server.URL)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error for invalid key, got nil")
	}
}

// errors also for unknown response text
func TestVerifyKeyStrangeResponse(t *testing.T) {
	server := makeServer(http.StatusOK, "some-unexpected-response")
	defer server.Close()

	_, err := newClientWithApiBaseURL(context.Background(), "bad-key", "http://example.com", server.URL)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error for invalid key, got nil")
	}
}

func TestVerifyKey_RequestBody(t *testing.T) {
	var gotKey, gotBlog string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotKey = r.FormValue("key")
		gotBlog = r.FormValue("blog")
		w.Write([]byte(bodyValidMessage))
	}))
	defer server.Close()

	myKey := "my-key"
	myBlog := "http://myblog.example.com"

	newClientWithApiBaseURL(context.Background(), myKey, myBlog, server.URL)

	if gotKey != myKey {
		t.Errorf("expected key %q, got %q", myKey, gotKey)
	}
	if gotBlog != myBlog {
		t.Errorf("expected blog %q, got %q", myBlog, gotBlog)
	}
}
