package main

import (
	akismet "akismet-api/src"
	"context"
	"fmt"
	"time"

	"log"
	"os"

	"github.com/joho/godotenv"
)

var akismetClient *akismet.Client

func main() {
	godotenv.Load()
	apiKey := os.Getenv("AKISMET_API_KEY")
	blogURL := os.Getenv("BLOG_URL")

	log.Printf("Using API key: %s\n", apiKey)

	client, err := akismet.NewClient(context.Background(), apiKey, blogURL)
	if err != nil {
		log.Fatal(err)
	}
	akismetClient = client

	checkComment()

}

func checkComment() {

	//comment := akismet.NewComment("akismet-guaranteed-spam", akismet.CommentTypeComment).
	comment := akismet.NewComment("Tack för idag, det var trevligt", akismet.CommentTypeComment).
		WithIsTest(true).
		WithUserIP("127.0.0.1").
		WithUserAgent("Mozilla/5.0").
		WithAuthor("John Doe").
		//WithAuthorEmail("akismet-guaranteed-spam@example.com").
		WithAuthorEmail("jd@example.com").
		WithAuthorURL("http://example.com/authors/jd").
		WithReferrer("https://google.com").
		WithPermalink("https://example.com/post/123").
		WithDateGMT(time.Now()).
		WithPostModifiedGMT(time.Now()).
		WithBlogLang("en-US").
		WithBlogCharset("utf-8").
		WithUserRole("user") // or "administrator" etc
		//WithPostContextTags([]string{"hello", "friend"}).
		//WithRecheckReason("user edited their comment").
		//WithHoneypot("pot", "i assure you! i am njet a botskij!")

	//fmt.Printf("Checking comment %+v\n", comment)

	result, err := akismetClient.CheckComment(context.Background(), comment)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("comment check result: %v\n", result)

}
