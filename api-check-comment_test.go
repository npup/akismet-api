package akismet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckComment_Ham(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, bodyHamResponse)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.CheckComment(context.Background(), mockComment("hello, world!", "1.2.3.4", "Mozilla/5.0"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.IsSpam {
		t.Error("expected IsSpam false, got true")
	}
	if result.Discard {
		t.Error("expected Discard false, got true")
	}
}

func TestCheckComment_Spam(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, bodySpamResponse)
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.CheckComment(context.Background(), mockComment("buy cheap meds", "1.2.3.4", "Mozilla/5.0"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !result.IsSpam {
		t.Error("expected IsSpam true, got false")
	}
}

func TestCheckComment_Discard(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		w.Header().Set(akismetHeaders.ProTip, headerProtipDiscardResponse)
		w.Write([]byte(bodySpamResponse))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.CheckComment(context.Background(), mockComment("worst spam ever", "1.2.3.4", "Mozilla/5.0"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !result.IsSpam {
		t.Error("expected IsSpam true, got false")
	}
	if !result.Discard {
		t.Error("expected Discard true, got false")
	}
}

func TestCheckComment_RecheckAfter(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		w.Header().Set(akismetHeaders.RecheckAfter, "30")
		w.Write([]byte(bodyHamResponse))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.CheckComment(context.Background(), mockComment("some comment", "1.2.3.4", "Mozilla/5.0"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.RecheckAfter == nil {
		t.Fatal("expected RecheckAfter to be set, got nil")
	}
	if *result.RecheckAfter != 30*time.Second {
		t.Errorf("expected RecheckAfter 30s, got %v", *result.RecheckAfter)
	}
}

func TestCheckComment_AkismetGUID(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		w.Header().Set(akismetHeaders.GUID, "abc-123-guid")
		w.Write([]byte(bodyHamResponse))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	result, err := client.CheckComment(context.Background(), mockComment("some comment", "1.2.3.4", "Mozilla/5.0"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.AkismetGUID != "abc-123-guid" {
		t.Errorf("expected AkismetGUID %q, got %q", "abc-123-guid", result.AkismetGUID)
	}
}

func TestCheckComment_UnexpectedResponse(t *testing.T) {
	server := makeVerifyThenRespond(bodyValidMessage, http.StatusOK, "something unexpected")
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	_, err := client.CheckComment(context.Background(), mockComment("some comment", "1.2.3.4", "Mozilla/5.0"))
	if err == nil {
		t.Fatal("expected error for unexpected response, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestCheckComment_ServerDown(t *testing.T) {
	server := makeServer(http.StatusOK, bodyValidMessage)
	serverURL := server.URL
	client, clientErr := newClientWithApiBaseURL(context.Background(), "test-key", "http://example.com", serverURL)
	if clientErr != nil {
		t.Fatalf("failed to create client: %v", clientErr)
	}
	server.Close()

	_, err := client.CheckComment(context.Background(), mockComment("some comment", "1.2.3.4", "Mozilla/5.0"))
	if err == nil {
		t.Fatal("expected error when server is down, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestCheckComment_HitsCorrectEndpoint(t *testing.T) {
	var gotPath string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		gotPath = r.URL.Path
		w.Write([]byte(bodyHamResponse))
	}))
	defer server.Close()

	client := makeClientWithServer(t, server, "my-apikey", "http://my-blog.example.com")
	client.CheckComment(context.Background(), mockComment("some comment", "1.2.3.4", "Mozilla/5.0"))

	if gotPath != "/1.1/comment-check" {
		t.Errorf("expected path /1.1/comment-check, got %s", gotPath)
	}
}

func TestCheckComment_RequestBody(t *testing.T) {
	var gotAPIKey, gotUserIP, gotUserAgent, gotCommentContent string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(bodyValidMessage))
			return
		}
		r.ParseForm()
		gotAPIKey = r.FormValue("api_key")
		gotUserIP = r.FormValue("user_ip")
		gotUserAgent = r.FormValue("user_agent")
		gotCommentContent = r.FormValue("comment_content")
		w.Write([]byte(bodyHamResponse))
	}))
	defer server.Close()

	myApiKey := "my-apikey"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"
	myCommentText := "hello world"

	client := makeClientWithServer(t, server, myApiKey, "http://my-blog.example.com")
	client.CheckComment(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))

	if gotAPIKey != myApiKey {
		t.Errorf("expected api_key %q, got %q", myApiKey, gotAPIKey)
	}
	if gotUserIP != myUserIp {
		t.Errorf("expected user_ip %q, got %q", myUserIp, gotUserIP)
	}
	if gotUserAgent != myUserAgent {
		t.Errorf("expected user_agent %q, got %q", myUserAgent, gotUserAgent)
	}
	if gotCommentContent != myCommentText {
		t.Errorf("expected comment_content %q, got %q", myCommentText, gotCommentContent)
	}
}
