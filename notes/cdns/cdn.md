## CNDs ( AMAZON CloudFront, Cloudflare, Akamai, etc.)
- CDNs (Content Delivery Networks) are a network of servers that deliver web content to users based on their geographic location. They help improve website performance, reduce latency, and enhance user experience by caching content closer to the end-users.
- CDNs work by distributing content across multiple servers located in different regions. When a user requests content, the CDN routes the request to the nearest server, reducing the distance the data has to travel and improving load times.
- CDNs can also provide additional security features, such as DDoS protection, SSL/TLS encryption, and Web Application Firewalls (WAFs) to protect against malicious attacks.

## A CDN like CloudFront has two purposes:
- Speed: Users get content from the server closest to them, which is faster than getting it from the origin server.
- Security: The origin server is hidden from the public internet, and only the CDN can access it. This is a security measure that can help prevent DDoS attacks and other malicious activity.

- Images and videos are certainly common, but in reality any static asset is a good fit for a CDN:
```
Images
HTML
CSS
JS
```

- Serving files from Direct S3 is not a good idea because it is slow and expensive. A CDN is a better option for serving static assets. Because, CDNs cache the content and serve it from the edge locations, which are closer to the users, resulting in faster load times and reduced latency. Additionally, CDNs can handle high traffic loads and provide better scalability compared to serving files directly from S3. While, if served from S3 any request will have to go all the way to the S3 bucket (origin region), which can be slow and expensive, especially for large files or high traffic websites.

```
By default, S3 does not store multiple versions of an object. If you upload a file to a key that already contains an object, the old object is overwritten.

Bucket versioning is an optional feature where the bucket stores multiple versions of an object. It helps:
Prevent accidental deletion
Rollback to previous versions of files
Store multiple versions of files in the same key
```
