# Serverless Architecture & AWS S3 Fundamentals

> This lesson introduces the serverless architecture paradigm, clarifies the 'serverless' misnomer, distinguishes between serverless compute and serverless storage, and provides hands-on practice creating and configuring an AWS S3 bucket with public read access via bucket policy.

## Prerequisites
- Basic understanding of cloud computing concepts
- AWS account with console access
- AWS CLI installed and configured with credentials
- Understanding of JSON syntax for policies
- Basic knowledge of HTTP/REST APIs (S3 uses HTTP)

## Key Concepts

### Serverless Architecture

Serverless is an architectural pattern where the cloud provider manages server provisioning, scaling, maintenance, and security. The term is misleading — servers still exist, but operational responsibility shifts to the provider (AWS, Google Cloud, Azure). Developers focus on code/logic rather than infrastructure.

**Examples:**
- AWS Lambda — serverless compute (run code without managing servers)
- Google Cloud Functions — serverless compute
- Azure Functions — serverless compute
- AWS S3 — serverless storage (store files without managing storage servers)

> ⚠️ **Gotchas:**
> - Serverless ≠ no servers. It means 'servers are someone else's problem.'
> - Serverless compute (Lambda) and serverless storage (S3) are different categories — don't conflate them.
> - Cold starts can affect serverless compute performance (not applicable to S3 storage).
> - Vendor lock-in is a consideration when adopting serverless services.

### Serverless Compute vs Serverless Storage

Serverless comes in two primary flavors. Serverless compute executes code on-demand (functions). Serverless storage provides durable, scalable object/file storage without managing underlying disks, RAID, or scaling logic. S3 pioneered serverless storage (launched 2006, before Lambda in 2014).

**Examples:**
- Compute: AWS Lambda, Cloudflare Workers, Vercel Edge Functions
- Storage: AWS S3, Google Cloud Storage, Azure Blob Storage, Cloudflare R2
- Database: DynamoDB (serverless NoSQL), Aurora Serverless (relational)

> ⚠️ **Gotchas:**
> - S3 is not a filesystem — it's an object store with a flat namespace (keys/prefixes), not directories.
> - S3 operations are HTTP API calls, not local syscalls (open/read/write). Latency and retry logic matter.
> - Serverless storage still has consistency models (S3: strong read-after-write for new objects, eventual for overwrites/deletes).

### AWS S3 Bucket Creation & Configuration

An S3 bucket is a top-level container for objects. Bucket names must be globally unique across all AWS accounts. Configuration choices at creation time affect security, cost, and functionality.

**Examples:**
- Bucket naming: tubely-56841 (must be globally unique, lowercase, no underscores, 3-63 chars)
- Block Public Access: UNCHECKED for this lesson (allows bucket policy to grant public read)
- Versioning: OFF (saves storage cost, but no object recovery from accidental overwrite/delete)
- Default Encryption: ON with SSE-S3 (AWS-managed keys) — encrypts objects at rest automatically
- Object Lock: OFF (prevents WORM compliance; enable only for regulatory requirements)

> ⚠️ **Gotchas:**
> - Bucket names are globally unique — not just per account. 'tubely' alone will fail.
> - Block Public Access settings override bucket policies. If 'Block all public access' is ON, the public-read policy won't work.
> - Versioning cannot be fully disabled once enabled — only suspended. Plan carefully.
> - SSE-S3 (AES-256) vs SSE-KMS: SSE-S3 uses AWS-managed keys (free), SSE-KMS uses customer-managed keys (costs $1/key/month + API calls).

### S3 Bucket Policy for Public Read Access

A bucket policy is a resource-based IAM policy attached to a bucket. This policy grants s3:GetObject permission to Principal '*' (anyone) on all objects (Resource: arn:aws:s3:::BUCKET_NAME/*), enabling direct URL access to files. It does NOT grant s3:ListBucket, so users cannot enumerate objects — they must know the exact object key/URL.

**Examples:**
- Policy JSON structure: Version, Statement array with Effect, Principal, Action, Resource
- Resource ARN format: arn:aws:s3:::bucket-name/* (the /* means all objects in bucket)
- Principal: '*' means anonymous/unauthenticated access allowed
- Action: s3:GetObject only — no PutObject, DeleteObject, ListBucket

> ⚠️ **Gotchas:**
> - This policy makes objects publicly readable via direct URL (https://bucket-name.s3.region.amazonaws.com/key).
> - No listing permission means no directory browsing — security through obscurity for object keys.
> - Bucket policies are limited to 20 KB. For complex permissions, use IAM roles/policies instead.
> - Public access via bucket policy requires 'Block Public Access' to be OFF at both account and bucket level.
> - Consider CloudFront + OAC (Origin Access Control) instead of public buckets for production — better security, caching, and HTTPS.

### AWS CLI Verification

The AWS CLI (Command Line Interface) provides programmatic access to AWS services. `aws s3 ls` lists all buckets in the configured account/region, confirming bucket creation succeeded.

**Examples:**
- aws s3 ls — lists all buckets
- aws s3 ls s3://tubely-56841 — lists objects in specific bucket
- aws s3 mb s3://tubely-56841 — create bucket via CLI (alternative to console)

> ⚠️ **Gotchas:**
> - Requires AWS credentials configured (~/.aws/credentials or environment variables).
> - Default region matters — buckets are regional but globally named. CLI uses default region unless --region specified.
> - Output format: '2024-01-15 10:30:45 tubely-56841' (creation date + bucket name).

## Mental Models

### 💡 Serverless = 'Managed Infrastructure as a Service'

Think of serverless like renting a fully furnished, maintained apartment vs. buying a house. You don't fix the roof (server maintenance), you don't add rooms when guests arrive (auto-scaling), and you don't install locks (security patches) — the landlord (AWS) handles it. You just live there (run your code/store your files).

### 💡 S3 as an Infinite, Durable Key-Value Store

Mental model: S3 = global, HTTP-accessible hash map with infinite capacity. Keys = object names (with / for pseudo-folders). Values = file bytes + metadata. No directories, no inodes, no filesystem hierarchy — just keys and values. Durability: 99.999999999% (11 9's) — expect to lose 1 object per 10,000 years per 10M objects.

### 💡 Bucket Policy = 'Rules Posted on the Bucket's Front Door'

A bucket policy is like a sign on a storage unit: 'Anyone with a key (URL) can read contents, but no one can see the inventory list.' It's evaluated at request time. Principal '*' = 'any person.' Action 'GetObject' = 'read only.' Resource '/*' = 'all units inside.'

## Code Examples

### S3 Bucket Policy - Public Read Access

This bucket policy allows anyone (Principal: '*') to read objects (s3:GetObject) from the specified bucket. The Resource ARN uses /* to match all object keys. Version '2012-10-17' is the current policy language version.

```json
  {  
```

### AWS CLI - List Buckets

Lists all S3 buckets in the current AWS account/region. Verifies bucket creation. Output shows creation timestamp and bucket name.

```bash
aws s3 ls
```

### AWS CLI - Create Bucket (Alternative to Console)

Creates bucket via CLI. 'mb' = make bucket. Region specification recommended. Note: This doesn't configure Block Public Access, versioning, encryption, or object lock — those require additional commands.

```bash
aws s3 mb s3://tubely-56841 --region us-east-1
```

### AWS CLI - Configure Block Public Access Off

Disables all four Block Public Access settings via CLI. Required for bucket policy public read to work. Equivalent to unchecking 'Block all public access' in console.

```bash
aws s3api put-public-access-block --bucket tubely-56841 --public-access-block-configuration BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false
```

### AWS CLI - Enable Default Encryption (SSE-S3)

Enables default encryption with AWS-managed keys (SSE-S3 / AES256). All new objects encrypted at rest automatically. Free tier.

```bash
aws s3api put-bucket-encryption --bucket tubely-56841 --server-side-encryption-configuration 'Rules=[{ApplyServerSideEncryptionByDefault:{SSEAlgorithm:AES256}}]'
```

### AWS CLI - Disable Versioning

Ensures versioning is off (suspended). Note: Once versioning is enabled, it can only be suspended, never fully disabled.

```bash
aws s3api put-bucket-versioning --bucket tubely-56841 --versioning-configuration Status=Suspended
```

## Common Mistakes

- ❌ Leaving 'Block all public access' ON while adding a public bucket policy — policy appears valid but requests return 403 Forbidden.
- ❌ Using underscores in bucket names — invalid. Only lowercase letters, numbers, dots, hyphens allowed.
- ❌ Forgetting bucket names are globally unique — 'tubely' or 'my-bucket' will always be taken.
- ❌ Confusing s3:GetObject (read object) with s3:ListBucket (list objects) — they are separate permissions.
- ❌ Assuming S3 is a filesystem — trying to 'mkdir' or 'rename' folders. S3 has no rename; copy + delete required.
- ❌ Not configuring AWS CLI credentials before running `aws s3 ls` — results in 'Unable to locate credentials' error.
- ❌ Using SSE-KMS (customer-managed keys) when SSE-S3 (AWS-managed) suffices — adds cost and complexity unnecessarily.

## Key Takeaways

- ✅ Serverless means 'you don't manage servers' — not 'no servers exist.' Operational burden shifts to cloud provider.
- ✅ S3 was the first mainstream serverless service (2006) — serverless storage, not compute. Lambda came later (2014).
- ✅ S3 buckets are globally named, regional resources. Names must be unique across ALL AWS accounts worldwide.
- ✅ Block Public Access settings OVERRIDE bucket policies. Must be disabled for public-read policies to work.
- ✅ Bucket policies grant permissions at the bucket level. This policy grants anonymous s3:GetObject on all objects — direct URL access only, no listing.
- ✅ Default encryption (SSE-S3) is free and should be enabled by default. SSE-KMS costs money — use only when key control/audit required.
- ✅ AWS CLI is essential for automation and verification. `aws s3 ls` confirms bucket existence; `aws s3api` handles advanced config.
- ✅ Production workloads should use CloudFront + OAC instead of public S3 buckets — better security, performance, and cost control.

## Practice Questions

**Q1:** What does 'serverless' actually mean in cloud computing?

<details><summary>Answer</summary>

Serverless means the cloud provider manages server provisioning, scaling, patching, and maintenance. Servers still exist — developers just don't operate them. It's 'servers are someone else's problem.'

</details>

**Q2:** What is the difference between serverless compute and serverless storage? Give one example of each.

<details><summary>Answer</summary>

Serverless compute runs code on-demand without managing servers (e.g., AWS Lambda). Serverless storage stores data without managing storage infrastructure (e.g., AWS S3). S3 launched in 2006; Lambda in 2014.

</details>

**Q3:** Why must S3 bucket names be globally unique across all AWS accounts?

<details><summary>Answer</summary>

Because S3 bucket names form part of the global DNS namespace (e.g., bucket.s3.amazonaws.com). DNS names must be globally unique, so bucket names must be too.

</details>

**Q4:** You create a bucket policy allowing public s3:GetObject, but requests return 403 Forbidden. What is the most likely cause?

<details><summary>Answer</summary>

Block Public Access settings are enabled (at account or bucket level). These settings override bucket policies and must be disabled for public access to work.

</details>

**Q5:** What permissions does this bucket policy grant, and what does it NOT grant? { Effect: Allow, Principal: *, Action: s3:GetObject, Resource: arn:aws:s3:::my-bucket/* }

<details><summary>Answer</summary>

Grants: Anonymous read access to any object in 'my-bucket' via direct URL. Does NOT grant: Listing objects (s3:ListBucket), writing (PutObject), deleting (DeleteObject), or accessing bucket metadata.

</details>

**Q6:** What is the difference between SSE-S3 and SSE-KMS encryption in S3?

<details><summary>Answer</summary>

SSE-S3 uses AWS-managed keys (AES-256), free, automatic. SSE-KMS uses customer-managed keys in AWS KMS, costs $1/key/month + API call fees, provides key rotation control and audit trail. Use SSE-S3 unless regulatory requirements demand SSE-KMS.

</details>

**Q7:** How would you verify an S3 bucket exists using the AWS CLI?

<details><summary>Answer</summary>

Run `aws s3 ls` to list all buckets, or `aws s3 ls s3://bucket-name` to check a specific bucket. Requires configured AWS credentials.

</details>

## Concept Relationships

- **Serverless Architecture** → **AWS S3**: S3 is a foundational serverless storage service — the first widely adopted serverless offering (2006).
- **Serverless Architecture** → **AWS Lambda**: Lambda is the canonical serverless compute service (2014). Often paired with S3 for event-driven architectures (S3 triggers Lambda).
- **S3 Bucket Creation** → **Block Public Access Settings**: Block Public Access must be disabled at creation (or after) for public bucket policies to take effect.
- **Bucket Policy** → **Block Public Access Settings**: Bucket policies are evaluated AFTER Block Public Access. If BPA blocks public policies, the policy is ignored.
- **Default Encryption (SSE-S3)** → **S3 Object Storage**: Default encryption automatically applies SSE-S3 to all new objects. Existing objects unaffected unless rewritten.
- **AWS CLI** → **S3 Bucket Verification**: CLI provides programmatic verification (`aws s3 ls`) and configuration alternative to console.
- **Public Bucket Policy** → **CloudFront + OAC**: Public bucket policies are a simple but less secure pattern. Production uses CloudFront with Origin Access Control for private buckets + signed URLs.

---
*Generated by Boot.dev Notes*