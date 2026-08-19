# 12 — Beginner Glossary

## Blob

A large block of binary data.

## Binary

Data represented as bytes rather than ordinary human-readable text.

## Byte

A basic unit of digital data, usually containing 8 bits.

## Database

A system for storing and querying structured data.

## Filesystem

The operating system's system for organizing files and directories on storage devices.

## Object storage

Storage designed around objects/files rather than relational rows.

AWS S3 is a major example.

## S3

Amazon Web Services' object storage service.

## Metadata

Information describing something.

For a video:

```text
title
description
owner
created_at
thumbnail_url
```

## MIME type

A web-friendly description of file format.

Examples:

```text
image/png
image/jpeg
video/mp4
application/pdf
```

## Multipart/form-data

An HTTP encoding commonly used when sending files through forms.

## Base64

An encoding that represents binary data as text.

## Data URL

A URL that contains the data itself.

Example:

```text
data:image/png;base64,...
```

## Streaming

Sending or processing data progressively instead of loading the entire thing first.

## `io.Reader`

A Go interface representing something you can read data from.

## `io.Copy`

A Go function for copying data between readers/writers.

## CGO

Go's mechanism for interacting with C code.

SQLite drivers such as `go-sqlite3` can require CGO.

## Persistence

Data that survives process/server restarts.

Filesystem and database storage are persistent in normal operation; ordinary RAM is not.

## Cache

Temporarily stored data used to make repeated access faster.

## Content-Type

An HTTP header describing the type of content being sent.

Example:

```http
Content-Type: image/png
```
