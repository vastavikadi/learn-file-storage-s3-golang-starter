# 07 — MIME Types

## What is a MIME type?

A MIME type is a web-friendly description of a file format.

Format:

```text
type/subtype
```

Examples:

```text
image/png
image/jpeg
video/mp4
audio/mp3
text/html
application/pdf
```

## Content-Type

When a browser uploads a file, the request can contain:

```http
Content-Type: image/png
```

This tells the server what kind of media is being sent.

## Parsing the media type in Go

Use:

```go
mediaType, _, err := mime.ParseMediaType(
    header.Get("Content-Type"),
)
```

Import:

```go
import "mime"
```

## Validate thumbnail uploads

If Tubely only allows PNG and JPEG thumbnails:

```go
if mediaType != "image/png" &&
   mediaType != "image/jpeg" {
    respondWithError(
        w,
        http.StatusBadRequest,
        "Unsupported image type",
        nil,
    )
    return
}
```

A PDF should be rejected:

```text
application/pdf
```

because it is not an allowed thumbnail type.

## Important distinction

MIME type answers:

> What format is this file intended to be?

It does not describe:

- File size
- Security level
- Creation date

## Security reminder

Do not blindly trust user uploads.

At minimum, consider validating:

- Authentication
- Ownership
- MIME type
- File size
- Storage path/name
- Application-specific rules
