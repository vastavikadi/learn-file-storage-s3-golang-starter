# 13 — Final Review

## The one-sentence summary

> Databases are best suited to structured metadata, while large binary files are usually better stored in filesystems or object storage such as S3.

## The storage progression

```text
In-memory
    ↓
Easy but temporary

Base64 in database
    ↓
Persistent but inefficient

Filesystem
    ↓
Good for simple/single-server applications

S3 / object storage
    ↓
Designed for scalable distributed storage
```

## Upload pipeline

```text
User
 ↓
Browser
 ↓
multipart/form-data
 ↓
Go HTTP handler
 ↓
Parse form
 ↓
Get file
 ↓
Read Content-Type
 ↓
Validate
 ↓
Check authorization/ownership
 ↓
Save file
 ↓
Store file URL/reference in DB
 ↓
Return metadata
```

## Database/file separation

Good design:

```text
SQLite
├── video_id
├── title
├── description
└── thumbnail_url
             |
             v
        /assets/123.png
             |
             v
        Actual image
```

At scale:

```text
SQLite / PostgreSQL
       |
       | reference
       v
      S3
       |
       ├── image
       ├── thumbnail
       └── video
```

## Three Go operations to remember

### Read everything

```go
data, err := io.ReadAll(file)
```

### Copy/stream

```go
_, err := io.Copy(destination, source)
```

### Create file

```go
output, err := os.Create(path)
```

## MIME validation

```go
mediaType, _, err := mime.ParseMediaType(
    header.Get("Content-Type"),
)

if mediaType != "image/png" &&
   mediaType != "image/jpeg" {
    // reject
}
```

## Questions you should be able to answer

### Why not store a 500 MB video in a normal relational database column?

Because relational databases are optimized primarily for structured data, while large binary objects can negatively affect storage, performance, backups, queries, and operational complexity.

### Why is in-memory storage bad for permanent files?

RAM contents disappear when the process/server restarts.

### Why is Base64 not ideal?

It adds encoding/decoding work and increases the amount of storage required.

### Why is a filesystem better?

Filesystems are designed to store and retrieve files efficiently.

### Why is S3 useful?

It provides centralized object storage that can be shared across multiple application servers and scales much better than tying assets to one machine.

### What is multipart/form-data?

An HTTP encoding commonly used to send files and other form fields together.

### What does a MIME type describe?

The web-friendly format/type of a piece of content.

### Why stream video?

Because videos can be extremely large, and streaming avoids unnecessarily loading the entire file before playback.

## Final mental model

```text
             APPLICATION
                  |
        +---------+---------+
        |                   |
   Structured data       Large assets
        |                   |
        v                   v
    Database         Filesystem / S3
        |                   |
        +------- reference-+
```

Keep this architecture in mind as the course moves into S3 and more advanced large-file handling.
