# akismet-api

A Go client library for the [Akismet](https://akismet.com) spam detection API.
To read up on the Akismet concept, visit their guide: [https://akismet.com/developers/getting-started/](https://akismet.com/developers/getting-started/).

## Requirements

An Akismet API key is required. You can get one somewhere at [akismet.com](https://akismet.com/).

## Features

- Check comments for spam
- Report false negatives (missed spam) and false positives (ham)
- Check API key usage limits
- List sites using your API key
- Full access to Akismet response metadata: discard hints, alert codes, debug help, recheck-after

## Installation

```
go get akismet-go
```

## Usage

### Creating a client

```go
client, err := akismet.NewClient(ctx, "your-api-key", "https://yourblog.com")
if err != nil {
    log.Fatal(err)
}
```

`NewClient` verifies the API key against Akismet on creation and returns an `*AkismetError` if the key is invalid or the request fails.

### Checking a comment

```go
comment := akismet.NewComment("Great post!", akismet.CommentTypeComment).
    WithUserIP("1.2.3.4").
    WithUserAgent("Mozilla/5.0").
    WithAuthor("Jane Doe").
    WithAuthorEmail("jane@example.com")

result, err := client.CheckComment(ctx, comment)
if err != nil {
    log.Fatal(err)
}

if result.IsSpam {
    if result.Discard {
        // blatant spam — safe to discard without saving
    } else {
        // spam — save for review
    }
}

if result.RecheckAfter != nil {
    // Akismet requests a recheck after this duration;
    // resubmit the comment with WithRecheckReason("recheck")
}
```

`CheckCommentResult` fields:

| Field          | Type             | Description                                                              |
| -------------- | ---------------- | ------------------------------------------------------------------------ |
| `IsSpam`       | `bool`           | True if Akismet considers the comment spam                               |
| `Discard`      | `bool`           | True if Akismet considers it blatant spam safe to discard without saving |
| `RecheckAfter` | `*time.Duration` | If set, Akismet requests a recheck after this duration                   |
| `AkismetGUID`  | `string`         | Unique request identifier, useful when contacting Akismet support        |

### Reporting spam and ham

```go
// Report a comment that was not caught (false negative)
err := client.ReportSpam(ctx, comment)

// Report a comment that was wrongly flagged (false positive)
err := client.ReportHam(ctx, comment)
```

### Getting usage limits

```go
result, err := client.GetUsageLimit(ctx)
if err != nil {
    log.Fatal(err)
}

if result.Limit == nil {
    fmt.Println("unlimited plan")
} else {
    fmt.Printf("%d / %d calls used (%.1f%%)\n", result.Usage, *result.Limit, result.Percentage)
}

if result.Throttled {
    fmt.Println("requests are currently being throttled due to overuse")
}
```

`UsageLimitResult` fields:

| Field        | Type      | Description                                              |
| ------------ | --------- | -------------------------------------------------------- |
| `Limit`      | `*int`    | Monthly API call allowance; nil if the plan has no limit |
| `Usage`      | `int`     | API calls made since the start of the current month      |
| `Percentage` | `float64` | Percentage of the monthly limit consumed                 |
| `Throttled`  | `bool`    | True if Akismet is throttling requests due to overuse    |

### Listing sites

```go
result, err := client.GetKeySites(ctx, nil) // nil = use defaults
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%d total sites\n", result.Total)
for _, site := range result.Sites {
    fmt.Printf("%s: %d calls, %d spam\n", site.Site, site.APICalls, site.Spam)
}
```

Pass a `*KeySitesParams` to filter or paginate:

```go
result, err := client.GetKeySites(ctx, &akismet.KeySitesParams{
    Month:  "2026-01",         // YYYY-MM; defaults to current month
    Filter: "example.com",     // filter by site URL substring
    Order:  akismet.KeySitesOrderSpam,
    Limit:  10,
    Offset: 20,
})
```

`KeySitesOrder` constants: `KeySitesOrderTotal`, `KeySitesOrderSpam`, `KeySitesOrderHam`, `KeySitesOrderMissedSpam`, `KeySitesOrderFalsePositives`.

`KeySiteEntry` fields:

| Field            | Type     | Description                                           |
| ---------------- | -------- | ----------------------------------------------------- |
| `Site`           | `string` | Site URL                                              |
| `APICalls`       | `int`    | Total API calls                                       |
| `Spam`           | `int`    | Comments correctly identified as spam                 |
| `Ham`            | `int`    | Comments correctly identified as ham                  |
| `MissedSpam`     | `int`    | Spam not caught (false negatives)                     |
| `FalsePositives` | `int`    | Ham incorrectly flagged as spam                       |
| `IsRevoked`      | `bool`   | True if the site's access to the key has been revoked |

## Comments

`NewComment(content, type)` returns a `*Comment`. All fields beyond content and type are optional — the more context you provide, the more accurate Akismet's evaluation will be. `UserIP` and `UserAgent` are especially valuable.

Builder methods: `WithUserIP`, `WithUserAgent`, `WithAuthor`, `WithAuthorEmail`, `WithAuthorURL`, `WithReferrer`, `WithPermalink`, `WithDateGMT`, `WithPostModifiedGMT`, `WithBlogLang`, `WithBlogCharset`, `WithUserRole`, `WithIsTest`, `WithRecheckReason`, `WithCommentParent`, `WithHoneypot`, `WithPostContextTags`.

`CommentType` constants: `CommentTypeComment`, `CommentTypeForumPost`, `CommentTypeBlogPost`, `CommentTypeContactForm`, `CommentTypeSignup`, `CommentTypeMessage`, `CommentTypeReply`, `CommentTypeTweet`, `CommentTypePingback`, `CommentTypeTrackback`.

Set `WithUserRole("administrator")` to bypass spam checking entirely for trusted users.

## Error handling

All methods return `*AkismetError` (or nil on success). It implements the `error` interface and unwraps via `errors.As`/`errors.Is`.

```go
result, err := client.CheckComment(ctx, comment)
if err != nil {
    if err.Alert != nil {
        fmt.Printf("alert %d: %s\n", err.Alert.Code, err.Alert.Description)
    }
    if err.DebugHelp != "" {
        fmt.Println("debug:", err.DebugHelp)
    }
    log.Fatal(err.Err) // underlying Go error
}
```

`AkismetError` fields:

| Field       | Type     | Description                                                     |
| ----------- | -------- | --------------------------------------------------------------- |
| `Err`       | `error`  | Underlying Go error                                             |
| `DebugHelp` | `string` | Value of the `X-Akismet-debug-help` response header, if present |
| `Alert`     | `*Alert` | Akismet account alert, if present (see below)                   |

`Alert` fields: `Code int`, `Message string`, `Description string`. Known alert codes and their descriptions are listed in `AlertDescriptionsByCode`.
