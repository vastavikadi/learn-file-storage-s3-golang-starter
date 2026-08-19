# 04 — Multipart Uploads

## Why multipart/form-data?

Small structured data is commonly sent as JSON:

```json
{
  "title": "My Video",
  "description": "A cool video"
}
```

Large files are commonly uploaded using:

```text
multipart/form-data
```

Think of multipart as one HTTP request containing several pieces:

```text
HTTP request
├── title
├── description
└── thumbnail file
```

HTML forms that upload files commonly use:

```html
<form enctype="multipart/form-data">
```

## Go support

Go's `net/http` package supports multipart forms.

```go
const maxMemory = 10 << 20

err := r.ParseMultipartForm(maxMemory)
if err != nil {
    // handle error
}
```

### What does `10 << 20` mean?

It is a bit shift.

For this purpose:

```text
10 << 20
=
10 * 1024 * 1024
=
10 MiB approximately
```

## Getting the uploaded file

If the form field is named `thumbnail`:

```go
file, header, err := r.FormFile("thumbnail")
if err != nil {
    // handle error
}
defer file.Close()
```

Conceptually:

```text
file
  → actual uploaded data

header
  → metadata about the upload

err
  → error information
```

## `io.Reader`

An uploaded file can behave like an `io.Reader`.

That means your program can read data from it.

Read the entire file:

```go
data, err := io.ReadAll(file)
```

The result is:

```go
[]byte
```

## Important distinction

`io.ReadAll`:

```text
source
  |
  v
read everything
  |
  v
[]byte in memory
```

`io.Copy`:

```text
source
  |
  | stream/copy
  v
destination
```

For large files, `io.Copy` is often preferable when you simply need to move data from one place to another.

## Simplified upload handler

```go
func handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
    const maxMemory = 10 << 20

    err := r.ParseMultipartForm(maxMemory)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid form", err)
        return
    }

    file, header, err := r.FormFile("thumbnail")
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Unable to get file", err)
        return
    }
    defer file.Close()

    // Validate the file and save it...
}
```
