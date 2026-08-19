# 05 — Storage Evolution

The course intentionally moves through several storage strategies.

```text
In-memory
    ↓
Base64 in database
    ↓
Filesystem
    ↓
AWS S3
```

Each step teaches why the next approach can be better.

## 1. In-memory storage

Conceptually:

```text
video ID
   ↓
thumbnail bytes in RAM
```

Example:

```go
var videoThumbnails map[int]thumbnail
```

### Problem

RAM is temporary.

```text
Server starts
    ↓
Upload image
    ↓
Image stored in RAM
    ↓
Server restarts
    ↓
Image disappears
```

Useful for learning or temporary data, but not reliable persistent storage.

## 2. Base64 in SQLite

Binary data can be converted into text using Base64.

```go
encoded := base64.StdEncoding.EncodeToString(data)
```

Import:

```go
import "encoding/base64"
```

Then a data URL can be created:

```text
data:<media-type>;base64,<data>
```

Example:

```text
data:image/png;base64,iVBORw0KGgo...
```

### Why it works

The database can store the resulting string.

### Why it is usually a bad idea

#### CPU

Encoding and decoding requires extra work.

#### Storage

Base64 increases the size of the representation.

#### Database performance

Relational databases are optimized for structured data, not giant binary blobs.

#### Caching

Dedicated file/object storage is generally better suited to serving large assets.

## 3. Filesystem

Store the actual file:

```text
/assets/123.png
```

Store only a reference in SQLite:

```text
thumbnail_url = /assets/123.png
```

This is a much cleaner separation.

## 4. S3

At larger scale, cloud object storage such as AWS S3 is useful.

It allows multiple application servers to share centralized storage.

```text
             Load Balancer
             /     |                 /      |             Server A Server B Server C
            \      |      /
             \     |     /
                AWS S3
```

## Key principle

> Store metadata/reference information in the database; store large files in storage designed for large files.
