# 03 — Database and Videos

## SQLite

Tubely uses SQLite for structured application data.

SQLite is a relational database stored in a file rather than requiring a separate database server.

The database file is:

```text
tubely.db
```

Open it:

```bash
sqlite3 tubely.db
```

Inspect users:

```sql
SELECT *
FROM users;
```

Exit:

```text
.exit
```

Check SQLite installation:

```bash
sqlite3 --version
```

## Tubely video model

A video has three important areas:

```text
Video
├── Metadata
│   ├── title
│   ├── description
│   └── other information
├── Thumbnail
│   └── image file
└── Video
    └── actual video file
```

## Draft workflow

Creating a video draft initially creates a database record containing metadata.

The files can be uploaded separately.

This separation is useful because metadata and large binary assets have different storage requirements.

## API flow

The course tests login using:

```http
POST /api/login
```

Then authenticated video access:

```http
GET /api/videos
Authorization: Bearer <token>
```

A video record might conceptually look like:

```json
{
  "title": "Boots, an Emote Story",
  "description": "A short film about the many faces of Boots",
  "thumbnail_url": "/assets/123.png"
}
```

The important point is that `thumbnail_url` can be a reference to a file rather than the file's binary contents.
