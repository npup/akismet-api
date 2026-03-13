package akismet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockComment(text string, userIp string, userAgent string) *Comment {
	return NewComment(text, CommentTypeComment).
		WithUserIP(userIp).
		WithUserAgent(userAgent)
}

func TestReportSpam_Success(t *testing.T) {
	server := makeVerifyThenRespond(BODY_VALID_MESSAGE, http.StatusOK, BODY_REPORT_SUCCESS_MESSAGE)
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	err := client.ReportSpam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestReportHam_Success(t *testing.T) {
	server := makeVerifyThenRespond(BODY_VALID_MESSAGE, http.StatusOK, BODY_REPORT_SUCCESS_MESSAGE)
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	err := client.ReportHam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestReportSpam_HitsCorrectEndpoint(t *testing.T) {
	var gotPath string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(BODY_VALID_MESSAGE))
			return
		}
		gotPath = r.URL.Path
		w.Write([]byte(BODY_REPORT_SUCCESS_MESSAGE))
	}))
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	client.ReportSpam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))

	if gotPath != "/1.1/submit-spam" {
		t.Errorf("expected path /1.1/submit-spam, got %s", gotPath)
	}
}

func TestReportHam_HitsCorrectEndpoint(t *testing.T) {
	var gotPath string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(BODY_VALID_MESSAGE))
			return
		}
		gotPath = r.URL.Path
		w.Write([]byte(BODY_REPORT_SUCCESS_MESSAGE))
	}))
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	client.ReportHam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))

	if gotPath != "/1.1/submit-ham" {
		t.Errorf("expected path /1.1/submit-ham, got %s", gotPath)
	}
}

func TestReport_UnexpectedResponse(t *testing.T) {
	server := makeVerifyThenRespond(BODY_VALID_MESSAGE, http.StatusOK, "something went wrong")
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	err := client.ReportSpam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))
	if err == nil {
		t.Fatal("expected error for unexpected response, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestReport_ServerDown(t *testing.T) {
	server := makeServer(http.StatusOK, BODY_VALID_MESSAGE)
	serverURL := server.URL

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"

	client, clientErr := newClientWithApiBaseURL(context.Background(), myApiKey, myBlogURL, serverURL)
	if clientErr != nil {
		t.Fatalf("failed to create client: %v", clientErr)
	}
	server.Close()

	// mock data
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	err := client.ReportSpam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))
	if err == nil {
		t.Fatal("expected error when server is down, got nil")
	}
	if err.Err == nil {
		t.Fatal("expected inner error, got nil")
	}
}

func TestReport_RequestBody(t *testing.T) {
	var gotAPIKey, gotUserIP, gotUserAgent string
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte(BODY_VALID_MESSAGE))
			return
		}
		r.ParseForm()
		gotAPIKey = r.FormValue("api_key")
		gotUserIP = r.FormValue("user_ip")
		gotUserAgent = r.FormValue("user_agent")
		w.Write([]byte(BODY_REPORT_SUCCESS_MESSAGE))
	}))
	defer server.Close()

	// mock data
	myApiKey := "my-apikey"
	myBlogURL := "http://my-blog.example.com"
	myCommentText := "some comment"
	myUserIp := "1.2.3.4"
	myUserAgent := "Mozilla/5.0"

	client := makeClientWithServer(t, server, myApiKey, myBlogURL)
	client.ReportSpam(context.Background(), mockComment(myCommentText, myUserIp, myUserAgent))

	if gotAPIKey != myApiKey {
		t.Errorf("expected api_key %q, got %q", myApiKey, gotAPIKey)
	}
	if gotUserIP != myUserIp {
		t.Errorf("expected user_ip %q, got %q", myUserIp, gotUserIP)
	}
	if gotUserAgent != myUserAgent {
		t.Errorf("expected user_agent %q, got %q", myUserAgent, gotUserAgent)
	}
}
