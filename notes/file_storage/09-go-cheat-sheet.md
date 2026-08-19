# 09 — Go Cheat Sheet

## `io.Reader`

Represents something that can be read from.

Uploaded files commonly behave like readers.

```go
func process(reader io.Reader) {
    // read from reader
}
```

## `io.ReadAll`

Reads everything into memory:

```go
data, err := io.ReadAll(file)
if err != nil {
    // handle error
}
```

Result:

```go
[]byte
```

Use this when you actually need the complete contents in memory.

## `io.Copy`

Copies data from a source to a destination:

```go
_, err := io.Copy(destination, source)
```

Useful for large files because you can copy data without first storing the entire file in a byte slice.

## `os.Create`

Creates a file:

```go
file, err := os.Create(path)
if err != nil {
    // handle error
}
defer file.Close()
```

## `filepath.Join`

Builds a filesystem path:

```go
path := filepath.Join(root, filename)
```

## `r.ParseMultipartForm`

Parses multipart/form-data:

```go
const maxMemory = 10 << 20

err := r.ParseMultipartForm(maxMemory)
```

## `r.FormFile`

Gets an uploaded file:

```go
file, header, err := r.FormFile("thumbnail")
```

## `base64.StdEncoding.EncodeToString`

Converts bytes to Base64:

```go
encoded := base64.StdEncoding.EncodeToString(data)
```

## `mime.ParseMediaType`

Extracts the media type:

```go
mediaType, _, err := mime.ParseMediaType(
    header.Get("Content-Type"),
)
```

## Imports used in this part of the course

```go
import (
    "encoding/base64"
    "io"
    "mime"
    "net/http"
    "os"
    "path/filepath"
)
```

## Typical upload flow

```go
func handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
    const maxMemory = 10 << 20

    if err := r.ParseMultipartForm(maxMemory); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid form", err)
        return
    }

    file, header, err := r.FormFile("thumbnail")
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Unable to get file", err)
        return
    }
    defer file.Close()

    mediaType, _, err := mime.ParseMediaType(
        header.Get("Content-Type"),
    )
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid media type", err)
        return
    }

    if mediaType != "image/png" && mediaType != "image/jpeg" {
        respondWithError(
            w,
            http.StatusBadRequest,
            "Unsupported image type",
            nil,
        )
        return
    }

    path := filepath.Join(cfg.assetsRoot, "123.png")

    output, err := os.Create(path)
    if err != nil {
        respondWithError(
            w,
            http.StatusInternalServerError,
            "Unable to create file",
            err,
        )
        return
    }
    defer output.Close()

    if _, err := io.Copy(output, file); err != nil {
        respondWithError(
            w,
            http.StatusInternalServerError,
            "Unable to save file",
            err,
        )
        return
    }
}
```

This is a learning skeleton, not a drop-in complete Tubely solution.
