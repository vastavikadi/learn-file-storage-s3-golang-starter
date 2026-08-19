# 08 — Video Streaming and Scale

## Why videos are harder

Images might be a few MB.

Videos can be:

```text
50 MB
500 MB
2 GB
10 GB+
```

Loading an entire large video into memory is wasteful.

Bad mental model:

```text
2 GB video
   ↓
load all into RAM
   ↓
process
```

Better:

```text
video
  ↓
small chunks
  ↓
client
  ↓
play while more data arrives
```

## What is streaming?

Streaming means processing/sending data progressively rather than waiting for the entire file.

Benefits:

- Lower memory usage
- Faster startup
- Better user experience
- Better handling of large assets
- Potentially lower bandwidth waste

## Local filesystem vs S3

### Filesystem

Good for:

- Local development
- Simple single-server applications
- Learning storage concepts

### S3/object storage

Good for:

- Multiple application servers
- Large-scale applications
- Centralized storage
- Large numbers of files
- Cloud deployments

## Distributed application problem

With local storage:

```text
             Load Balancer
             /     |                 /      |             Server A Server B Server C
          |        |        |
        Disk A   Disk B   Disk C
```

Files are tied to individual machines.

With centralized object storage:

```text
             Load Balancer
             /     |            Server A Server B Server C
            \      |      /
             \     |     /
                S3
```

Every server can access the same objects.

## Core architecture principle

```text
Database
  → metadata + references

Object storage
  → actual large files

Application
  → authentication + business logic + APIs
```
