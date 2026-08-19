# 01 — Course Overview

## What is this course about?

Web applications often need to handle files such as:

- Images
- Videos
- Audio
- PDFs
- Documents
- Static assets
- User uploads

These files can be much larger than ordinary database values.

The course uses **Tubely**, a SaaS-style application for managing video assets.

## Technologies

- **Go** — backend/application code
- **SQLite** — small structured data
- **Filesystem** — local large-file storage
- **AWS S3** — scalable object storage

## Learning goals

By the end, you should understand:

1. The difference between structured data and large files.
2. How file uploads work.
3. How Go handles multipart uploads.
4. Why storing large blobs directly in a database is usually a bad idea.
5. How to store files on a filesystem.
6. Why object storage such as S3 is useful at scale.
7. The basics of video streaming.

## Setup

Copy the environment file:

```bash
cp .env.example .env
```

Run the Go application:

```bash
go run .
```

The local app uses a default base URL similar to:

```text
http://localhost:8091
```

## Windows

The course recommends **WSL 2** for Windows users.

## AWS warning

The course uses AWS. Even if the course aims to stay within the free tier, AWS can charge for incorrect or unexpected usage.

When finished, remove resources you no longer need.

## Boot.dev CLI

Run tests:

```bash
bootdev run <lesson-id>
```

Submit:

```bash
bootdev run -s <lesson-id>
```

## CGO troubleshooting

If you see:

```text
go-sqlite3 requires cgo to work
```

Install GCC.

macOS:

```bash
brew install gcc
```

Linux:

```bash
sudo apt install gcc
```

Check CGO:

```bash
go env CGO_ENABLED
```

Enable it if necessary:

```bash
go env -w CGO_ENABLED=1
```

## Tubely architecture

At a high level:

```text
                 Tubely
                   |
          +--------+--------+
          |                 |
       Go App            SQLite
          |
          +---- Filesystem
          |
          +---- AWS S3
```

SQLite handles structured metadata. Filesystem/S3 handles large assets.
