# 10 — Commands Cheat Sheet

## Go

Start the app:

```bash
go run .
```

Check Go version:

```bash
go version
```

Check CGO:

```bash
go env CGO_ENABLED
```

Enable CGO:

```bash
go env -w CGO_ENABLED=1
```

## Environment

Copy environment configuration:

```bash
cp .env.example .env
```

## Course samples

Download:

```bash
./samplesdownload.sh
```

List:

```bash
ls samples
```

Inspect binary file:

```bash
xxd samples/boots-image-horizontal.png
```

Inspect first 8 bytes:

```bash
xxd -l 8 samples/boots-image-horizontal.png
```

## SQLite

Check version:

```bash
sqlite3 --version
```

Open database:

```bash
sqlite3 tubely.db
```

List/query users:

```sql
SELECT *
FROM users;
```

Exit:

```text
.exit
```

## Boot.dev CLI

Run tests:

```bash
bootdev run <lesson-id>
```

Submit:

```bash
bootdev run -s <lesson-id>
```

Configure a different base URL:

```bash
bootdev config base_url <url>
```

## Useful filesystem commands

List files:

```bash
ls
```

List a directory:

```bash
ls assets
```

Check a file:

```bash
ls -lh assets/123.png
```

## Useful mental shortcut

```text
go run .
    → start app

sqlite3 tubely.db
    → inspect database

xxd file
    → inspect raw binary data

bootdev run
    → test

bootdev run -s
    → submit
```
