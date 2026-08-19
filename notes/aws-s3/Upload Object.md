# Upload Object to S3

> This lesson covers the practical process of uploading a file (object) to an Amazon S3 bucket using the AWS Management Console, accessing it via its public URL, and verifying the upload using the AWS CLI. It demonstrates the fundamental workflow of storing and retrieving objects in S3.

## Prerequisites
- AWS account with S3 access
- AWS CLI installed and configured (aws configure)
- S3 bucket created with public access enabled
- Basic shell knowledge (redirection, cat command)
- Understanding of HTTP/HTTPS and URLs

## Key Concepts

### S3 Object Upload

Uploading a file to S3 creates an 'object' in a bucket. An object consists of the file data (value), a key (name/path), metadata, and optional version ID. The upload process stores the file in S3's distributed storage system with high durability (99.999999999% - eleven 9's).

**Examples:**
- Upload boots-image-horizontal.png via AWS Console drag-and-drop or 'Upload' button
- Use default settings: no encryption, standard storage class, no additional metadata

> ⚠️ **Gotchas:**
> - Default settings mean the object inherits bucket-level permissions - if bucket is public, object is public
> - Default storage class is S3 Standard (not Intelligent-Tiering, Glacier, etc.)
> - Default encryption is 'None' unless bucket has default encryption enabled
> - File name becomes the object key - spaces and special characters are allowed but can cause URL encoding issues

### S3 Object URL Structure

Every object in a public bucket has a publicly accessible URL. The URL format depends on the bucket's region and whether it uses virtual-hosted-style or path-style addressing. For modern buckets, virtual-hosted-style is standard.

**Examples:**
- Virtual-hosted-style: https://bucket-name.s3.region.amazonaws.com/object-key
- Path-style (deprecated for new buckets): https://s3.region.amazonaws.com/bucket-name/object-key
- Example: https://my-boot-bucket.s3.us-east-1.amazonaws.com/boots-image-horizontal.png

> ⚠️ **Gotchas:**
> - URL only works if bucket and object are public (no ACL blocking, bucket policy allows public read)
> - Region in URL must match bucket's region - wrong region returns 301 redirect or 400 error
> - Object key is case-sensitive and must match exactly including path separators
> - Special characters in key must be URL-encoded (spaces → %20, etc.)

### AWS CLI S3 Commands

The AWS CLI provides high-level S3 commands (s3) that simplify common operations like ls, cp, mv, rm, sync. These are built on top of the lower-level API commands (s3api). The 'aws s3 ls' command lists objects in a bucket or prefixes.

**Examples:**
- aws s3 ls s3://bucket-name -- lists all objects in bucket root
- aws s3 ls s3://bucket-name/prefix/ -- lists objects under prefix
- aws s3 ls s3://bucket-name --recursive -- lists all objects recursively
- aws s3 ls s3://bucket-name --human-readable --summarize -- shows sizes in human format with total

> ⚠️ **Gotchas:**
> - Must use 's3://' prefix for bucket references in CLI
> - Output redirection (>) captures stdout but not stderr - errors won't appear in file
> - Default output shows: LastModifiedDate Size Key (no bucket name in output)
> - Requires AWS credentials configured (aws configure) with appropriate IAM permissions

### Shell Redirection and Verification

Shell redirection (>) writes command output to a file, overwriting existing content. The 'cat' command displays file contents. This is a standard Unix pattern for capturing and verifying CLI output.

**Examples:**
- aws s3 ls s3://my-bucket > /tmp/bucket_contents.txt -- captures listing to file
- cat /tmp/bucket_contents.txt -- displays captured output
- aws s3 ls s3://my-bucket 2>&1 | tee /tmp/bucket_contents.txt -- captures both stdout and stderr

> ⚠️ **Gotchas:**
> - > overwrites file silently - use >> to append
> - If command fails, file may be empty or contain partial output
> - Always verify with cat after redirection to confirm capture worked

## Mental Models

### 💡 S3 as a Giant Hash Map

Think of S3 as a massive distributed hash map (dictionary) where the bucket is the top-level namespace and the object key is the key. The value is the file data + metadata. There are no real directories - 'folders' are just key prefixes with '/' delimiters. The 'aws s3 ls' command with --recursive scans all keys matching a prefix.

**Analogy:** Like a Python dict: bucket = dict, key = 'images/logo.png', value = {data: bytes, metadata: {...}, tags: {...}}. The '/' in keys is just a character - S3 doesn't enforce hierarchy.

### 💡 Object URL as a Direct Pointer

An S3 object URL is essentially a direct pointer to that specific object in S3's global namespace. When you paste it in a browser, the browser makes an HTTP GET request to S3, which checks permissions and returns the object data with appropriate headers (Content-Type, Content-Length, ETag, Cache-Control).

**Analogy:** Like a permanent, globally unique link to a file in a shared drive - but the 'drive' is distributed across multiple availability zones and the link works from anywhere on the internet (if public).

### 💡 CLI as a Remote Control

The AWS CLI is a remote control for your AWS account. Commands like 'aws s3 ls' send API requests to S3's control plane, which responds with data. The CLI formats this response as text. Redirection (>) captures that formatted text locally. You're not 'in' the bucket - you're querying it over HTTPS.

**Analogy:** Like using a TV remote - you press 'list' (ls), the TV (S3) sends back the channel guide, and you can write that guide to a notepad (file redirection).

## Code Examples

### Upload File via AWS CLI (Alternative to Console)

The 'aws s3 cp' command uploads files. By default it uses the bucket's default settings. You can override storage class, ACL, metadata, encryption, etc. This is useful for automation and scripts.

```bash
# Upload a local file to S3 bucket (alternative to console upload)  aws s3 cp boots-image-horizontal.png s3://my-bucket-name/boots-image-horizontal.png  # Upload with specific storage class  aws s3 cp boots-image-horizontal.png s3://my-bucket-name/boots-image-horizontal.png --storage-class STANDARD_IA  # Upload with public-read ACL (if bucket allows)  aws s3 cp boots-image-horizontal.png s3://my-bucket-name/boots-image-horizontal.png --acl public-read
```

### List Bucket Contents with Various Options

The 'aws s3 ls' command has several useful flags. --recursive is essential for seeing all objects since S3 is flat (no real folders). --human-readable makes sizes readable (KB, MB, GB). --summarize shows total objects and total size at the end.

```bash
# Basic list (shows objects in root only)  aws s3 ls s3://my-bucket-name  # Recursive list (shows all objects in all 'folders')  aws s3 ls s3://my-bucket-name --recursive  # Human-readable sizes with summary  aws s3 ls s3://my-bucket-name --recursive --human-readable --summarize  # List specific prefix (folder)  aws s3 ls s3://my-bucket-name/images/  # Save to file and verify  aws s3 ls s3://my-bucket-name --recursive > /tmp/bucket_contents.txt  cat /tmp/bucket_contents.txt
```

### Generate Presigned URL for Private Object Access

Presigned URLs grant temporary access to private objects. This is the secure way to share private files without making the bucket public. The URL contains a cryptographic signature that expires after the specified time.

```bash
# Generate a presigned URL valid for 1 hour (3600 seconds)  aws s3 presign s3://my-bucket-name/boots-image-horizontal.png --expires-in 3600  # Output: https://my-bucket-name.s3.region.amazonaws.com/boots-image-horizontal.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=...  # This URL works even for private objects!
```

### Upload with Metadata and Tags

S3 objects can have custom metadata (key-value pairs) and tags. Metadata is returned with HEAD/GET requests. Tags are used for cost allocation, lifecycle policies, and access control. Both are set at upload time or added later.

```bash
# Upload with custom metadata  aws s3 cp boots-image-horizontal.png s3://my-bucket-name/boots-image-horizontal.png --metadata '{
```

### Verify Object Exists and Get Details

The 's3api head-object' command retrieves object metadata without downloading the file. This is efficient for checking existence, size, content-type, and custom metadata. The ETag is typically an MD5 hash (for single-part uploads) useful for integrity verification.

```bash
# Check if object exists (returns metadata or error)  aws s3api head-object --bucket my-bucket-name --key boots-image-horizontal.png  # Get object metadata only  aws s3api head-object --bucket my-bucket-name --key boots-image-horizontal.png --query '{ContentLength: ContentLength, ContentType: ContentType, LastModified: LastModified, ETag: ETag, Metadata: Metadata}' --output table  # Download object to verify content  aws s3 cp s3://my-bucket-name/boots-image-horizontal.png /tmp/downloaded.png  # Compare checksums  sha256sum boots-image-horizontal.png /tmp/downloaded.png
```

## Common Mistakes

- ❌ Forgetting to replace BUCKET_NAME placeholder with actual bucket name in CLI commands
- ❌ Using path-style URL (s3.region.amazonaws.com/bucket/key) instead of virtual-hosted-style (bucket.s3.region.amazonaws.com/key) for new buckets
- ❌ Not configuring AWS CLI credentials before running commands (run 'aws configure' first)
- ❌ Assuming 'aws s3 ls' shows all objects without --recursive flag (only shows root-level by default)
- ❌ Thinking S3 has real directories - it doesn't, keys with '/' are just naming conventions
- ❌ Not URL-encoding special characters when constructing object URLs manually
- ❌ Using '>' redirection without verifying the file was created and has content
- ❌ Confusing bucket policies (resource-based) with IAM policies (identity-based) for public access

## Key Takeaways

- ✅ Uploading to S3 creates an object with a key (name), data, metadata, and optional version ID
- ✅ Default upload settings: Standard storage class, no encryption, no custom metadata, inherits bucket permissions
- ✅ Public objects are accessible via HTTPS URLs in format: https://bucket.s3.region.amazonaws.com/key
- ✅ AWS CLI 'aws s3 ls s3://bucket' lists objects; use --recursive to see all keys
- ✅ Shell redirection (>) captures command output to a file for verification or logging
- ✅ Always verify uploads by checking the object URL in browser AND using CLI to list bucket contents
- ✅ S3 is a flat key-value store - 'folders' are just key prefixes with '/' delimiter

## Practice Questions

**Q1:** What is the default storage class when uploading an object to S3 without specifying one?

<details><summary>Answer</summary>

S3 Standard (also called STANDARD). This is the default for all uploads unless you explicitly choose another class like STANDARD_IA, GLACIER, or INTELLIGENT_TIERING.

</details>

**Q2:** Why does 'aws s3 ls s3://my-bucket' only show some objects but not others in 'subdirectories'?

<details><summary>Answer</summary>

By default, 'aws s3 ls' only lists objects at the root level (no '/' in key) and common prefixes (simulated folders). To see all objects recursively, you must use the --recursive flag: 'aws s3 ls s3://my-bucket --recursive'.

</details>

**Q3:** What happens if you run 'aws s3 ls s3://my-bucket > /tmp/list.txt' and the bucket doesn't exist?

<details><summary>Answer</summary>

The command will fail with an error printed to stderr (not captured in the file), and /tmp/list.txt will be created but empty (or contain only partial output if some succeeded). Always check the file with 'cat' and check exit code ($?) to verify success.

</details>

**Q4:** How do you construct the public URL for an object named 'images/photo.jpg' in bucket 'my-app-assets' in region 'us-west-2'?

<details><summary>Answer</summary>

https://my-app-assets.s3.us-west-2.amazonaws.com/images/photo.jpg (virtual-hosted-style, which is the modern standard). The key 'images/photo.jpg' is used as-is in the path.

</details>

**Q5:** What does the ETag returned by S3 represent for a simple (non-multipart) upload?

<details><summary>Answer</summary>

For single-part uploads, the ETag is the MD5 hash of the object content (in hex). For multipart uploads, it's a different hash (MD5 of concatenated part MD5s with part count). It can be used to verify file integrity after download.

</details>

## Concept Relationships

- **S3 Object Upload** → **S3 Object URL Structure**: Uploading an object with public permissions makes it accessible via a predictable URL structure based on bucket name, region, and object key
- **AWS CLI S3 Commands** → **Shell Redirection and Verification**: CLI commands produce stdout that can be captured via shell redirection for verification, logging, or further processing
- **S3 Object Upload** → **AWS CLI S3 Commands**: Objects uploaded via Console can be listed, downloaded, and managed via CLI - both interfaces operate on the same S3 API

---
*Generated by Boot.dev Notes*