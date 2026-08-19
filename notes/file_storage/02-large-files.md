# 02 — Large Files

## Small structured data

Traditional databases are excellent for values such as:

```text
user_id   → 123
is_active → true
email     → user@example.com
```

Example SQL:

```sql
SELECT *
FROM users;
```

## Large files

Large files are blobs of binary data represented by formats such as:

```text
photo.png
video.mp4
song.mp3
document.pdf
```

A useful beginner rule:

> If it normally exists as its own file on your computer, it probably belongs in file/object storage.

## Why large files are special

Large files can:

1. Consume lots of storage.
2. Consume lots of bandwidth.
3. Require significant CPU or memory during processing.
4. Be accessed very frequently.
5. Become performance bottlenecks.

## Binary data

An image is not stored internally as readable text. It is a sequence of bytes.

Inspect a file with:

```bash
xxd samples/boots-image-horizontal.png
```

Inspect only the first 8 bytes:

```bash
xxd -l 8 samples/boots-image-horizontal.png
```

A PNG begins with the well-known signature:

```text
89 50 4e 47 0d 0a 1a 0a
```

This is part of how software identifies the file format.

## Sample files

Download course samples:

```bash
./samplesdownload.sh
```

List them:

```bash
ls samples
```

You should see files such as:

```text
boots-image-horizontal.png
boots-image-vertical.png
boots-video-horizontal.mp4
boots-video-vertical.mp4
is-bootdev-for-you.pdf
```

## Mental model

```text
Structured data
      |
      v
   Database

Large binary files
      |
      v
Filesystem / Object Storage
```
