# akismet-go

A Go client library for the [Akismet](https://akismet.com) spam detection API.

## Features

- Check comments for spam
- Report false negatives (missed spam) and false positives (ham)
- Check API key usage limits
- List sites using your API key
- Full access to Akismet response metadata: discard hints, alert codes, debug help, recheck-after

## Usage

```go
client, err := akismet.NewClient("your-api-key", "https://yourblog.com")
if err != nil {
    log.Fatal(err) // key verification failed
}

comment := akismet.NewComment("Great post!", akismet.CommentTypeComment).
    WithUserIP("1.2.3.4").
    WithUserAgent("Mozilla/5.0").
    WithAuthor("Jane Doe").
    WithAuthorEmail("jane@example.com")

result, akErr := client.CheckComment(comment)
if akErr != nil {
    log.Fatal(akErr)
}

if result.IsSpam {
    // handle spam
}
```

`NewClient` verifies the API key against Akismet on creation and returns an error if it is invalid.

The `Comment` type supports all fields documented in the Akismet API, including honeypot fields, post context tags, and recheck reasons. All fields beyond content and type are optional.

## Requirements

An Akismet API key is required. You can get one at [akismet.com](https://akismet.com).
