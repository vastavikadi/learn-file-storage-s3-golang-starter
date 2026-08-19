# 11 — Troubleshooting

## Error: go-sqlite3 requires cgo to work

Install GCC.

macOS:

```bash
brew install gcc
```

Linux:

```bash
sudo apt install gcc
```

Check:

```bash
go env CGO_ENABLED
```

If it prints `0`:

```bash
go env -w CGO_ENABLED=1
```

## Thumbnail upload succeeds but image doesn't appear

Check:

1. Is the server running?
2. Did the request succeed?
3. Is the MIME type allowed?
4. Is the file actually stored?
5. Is `thumbnail_url` correct?
6. Is the browser displaying a cached version?

## Browser cache

A browser can cache an image.

You may upload a new image but still see an old version.

When debugging, consider:

- Hard refresh
- Opening the image URL directly
- Using a new/private browser window

## File not found

If using filesystem storage, verify:

```text
/assets/<videoID>.<extension>
```

exists.

For example:

```bash
ls -lh assets/123.png
```

## Database reference mismatch

If the database says:

```text
/assets/123.png
```

but the file is:

```text
/assets/456.png
```

the browser cannot retrieve the expected file.

The database reference and physical/object-storage location must agree.

## Upload validation failure

Check the Content-Type.

Allowed examples:

```text
image/png
image/jpeg
```

Rejected example:

```text
application/pdf
```

## General debugging strategy

When something fails, trace the pipeline:

```text
Browser
  ↓
HTTP request
  ↓
multipart parser
  ↓
uploaded file
  ↓
validation
  ↓
storage
  ↓
database reference
  ↓
HTTP response
  ↓
browser retrieves asset
```

Find the first stage where reality differs from what you expect.
