# Large File Storage — Beginner Notes

A structured set of notes for learning and revising the large-file-storage course.

## Contents

1. `01-course-overview.md` — course goals, setup, Tubely architecture
2. `02-large-files.md` — large files vs structured data, binary files, hexdumps
3. `03-database-and-videos.md` — SQLite, video metadata, and the Tubely data model
4. `04-multipart-uploads.md` — multipart/form-data and Go upload handling
5. `05-storage-evolution.md` — in-memory → Base64 → filesystem → S3
6. `06-filesystem.md` — saving files to disk with Go
7. `07-mime-types.md` — MIME types and upload validation
8. `08-video-streaming-and-scale.md` — scaling, S3, and video streaming concepts
9. `09-go-cheat-sheet.md` — important Go APIs and code snippets
10. `10-commands-cheat-sheet.md` — terminal, SQLite, and Boot.dev commands
11. `11-troubleshooting.md` — common setup and upload problems
12. `12-glossary.md` — beginner-friendly terminology
13. `13-final-review.md` — condensed revision notes and mental models

## Recommended learning order

Read the files in numerical order. After each section, try the commands and small code examples yourself.

## Core principle

> Store information about a file in your database; store the actual large file in storage designed for files.
