# 06 — Filesystem Storage

## Why use the filesystem?

A filesystem is designed to store files directly.

Example:

```text
/assets/
    123.png
    456.jpg
    789.png
```

The database can store:

```text
thumbnail_url = /assets/123.png
```

The actual image stays on disk.

## Creating a path

Use `filepath.Join`:

```go
path := filepath.Join(cfg.assetsRoot, "123.png")
```

Avoid manually building filesystem paths with string concatenation.

## Creating a file

Use `os.Create`:

```go
output, err := os.Create(path)
if err != nil {
    // handle error
}
defer output.Close()
```

## Copying uploaded data

```go
_, err = io.Copy(output, file)
if err != nil {
    // handle error
}
```

Conceptually:

```text
multipart upload
      |
      | io.Copy
      v
/assets/123.png
```

## File extension

The application can use the MIME type to determine the appropriate extension.

For example:

```text
image/png  → .png
image/jpeg → .jpg
```

Then create a unique filename based on the video ID:

```text
<videoID>.<extension>
```

Example:

```text
123.png
```

## Serving assets

If the server exposes `/assets`, then:

```text
/assets/123.png
```

can be requested by the browser.

The database does not need to contain the image bytes.

## Why local filesystem eventually has limits

Imagine multiple servers:

```text
Server A → disk A
Server B → disk B
Server C → disk C
```

A file written to Server A's disk is not automatically available on B or C.

That is one reason object storage becomes important for distributed applications.
