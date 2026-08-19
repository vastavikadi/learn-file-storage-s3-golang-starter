# Cache Busting

> Cache busting is a technique to force browsers to fetch the latest version of a resource by modifying the resource URL (typically via query strings) so the browser treats it as a new resource while the server ignores the query parameters and serves the same file.

## Prerequisites
- Understanding of HTTP requests and responses
- Basic knowledge of browser developer tools (Network tab)
- Familiarity with JavaScript template literals and Date.now()
- Concept of static file serving

## Key Concepts

### Browser Caching

Browsers automatically cache static assets (images, CSS, JS) to improve load times and reduce bandwidth usage. When a browser requests a URL it has seen before, it may serve the cached copy instead of making a network request.

**Examples:**
- First visit: Browser downloads image.jpg (200 OK, caches it)
- Second visit: Browser serves image.jpg from cache (no network request)
- Result: Faster page loads, lower data costs for users

> ⚠️ **Gotchas:**
> - Cached content becomes stale when the underlying file changes on the server
> - Users may see outdated images, styles, or scripts
> - Cache duration depends on HTTP headers (Cache-Control, Expires, ETag)

### Cache Busting

A technique to invalidate browser cache by changing the resource URL. The browser sees a new URL and fetches fresh content, while the server ignores the query string and serves the same underlying file.

**Examples:**
- Original: http://example.com/image.jpg
- Busted: http://example.com/image.jpg?v=1
- Busted again: http://example.com/image.jpg?v=2

> ⚠️ **Gotchas:**
> - Only works if server ignores query strings for static assets (most do)
> - Doesn't work if server uses query strings for routing or processing
> - Each unique URL creates a separate cache entry, increasing storage

### Query String Cache Busting

Appending a query parameter (like ?v=timestamp or ?version=1) to a static asset URL. This is the simplest client-side cache busting method because it requires no server-side changes.

**Examples:**
- Version-based: image.jpg?version=1 → image.jpg?version=2
- Timestamp-based: image.jpg?v=1704067200000
- Hash-based: image.jpg?v=a1b2c3d4 (content hash)

> ⚠️ **Gotchas:**
> - Some CDNs/proxies may not cache URLs with query strings by default
> - Timestamp approach busts cache on EVERY page load (defeats caching benefits)
> - Version/hash approach only busts when content actually changes

### Client-Side vs Server-Side Cache Busting

Client-side: Modify URLs in frontend code (JavaScript, HTML templates). Server-side: Configure build tools to append content hashes to filenames (e.g., image.a1b2c3.jpg) or set appropriate Cache-Control headers.

**Examples:**
- Client-side: thumbnailImg.src = video.thumbnail_url + '?v=' + Date.now()
- Server-side (build tool): main.js → main.a1b2c3d4.js
- Server-side (headers): Cache-Control: max-age=31536000, immutable

> ⚠️ **Gotchas:**
> - Client-side timestamp busting defeats caching entirely - use sparingly
> - Server-side filename hashing is preferred for production apps
> - Query string busting is a quick fix; filename hashing is a robust solution

## Mental Models

### 💡 The Library Card Catalog Analogy

Think of the browser cache as a library's card catalog. The URL is the call number. When you change the call number slightly (add a sticky note with a version), the librarian (browser) thinks it's a different book and fetches it from the stacks (server). But the stacks (server) ignore the sticky note and give you the same physical book.

**Analogy:** URL = Call number, Query string = Sticky note on card, Server = Library stacks that ignore sticky notes

### 💡 The Cache Key Concept

Browsers use the full URL (including query string) as the cache key. Two URLs that differ only in query string are completely different cache entries. The server's routing logic typically strips query strings before looking up files, so both URLs map to the same physical resource.

**Analogy:** Cache key = Full URL string, Server lookup = Path only (no query string)

## Code Examples

### Basic Cache Busting with Timestamp (Lesson Exercise)

Uses template literals (backticks) and Date.now() to append a unique query parameter on every page load. This ensures the browser always fetches the latest thumbnail after an upload.

```javascript
// Original code - no cache busting, browser may serve stale thumbnail thumbnailImg.src = video.thumbnail_url; // Updated with cache busting - appends current timestamp in milliseconds thumbnailImg.src = `${video.thumbnail_url}?v=${Date.now()}`; // Result: https://example.com/thumb.jpg?v=1704067200123
```

### Version-Based Cache Busting (Better for Production)

Uses a version number instead of timestamp. Cache is only busted when you intentionally increment the version, preserving caching benefits for unchanged content.

```javascript
// Maintain a version number that increments when content changes const THUMBNAIL_VERSION = 2; // Increment this when thumbnail is updated thumbnailImg.src = `${video.thumbnail_url}?v=${THUMBNAIL_VERSION}`; // Result: https://example.com/thumb.jpg?v=2 // Only busts cache when version actually changes
```

### Content Hash Cache Busting (Build-Time, Most Robust)

Build tools compute a hash of file contents and embed it in the filename. This is the gold standard because: 1) Cache busts only when content changes, 2) Long-term caching with immutable flag is safe, 3) No query string issues with CDNs.

```javascript
// This is typically done by build tools (Webpack, Vite, etc.) // Filename includes content hash: thumbnail.a1b2c3d4.jpg // HTML references the hashed filename: <img src='/assets/thumbnail.a1b2c3d4.jpg' alt='Video thumbnail'> // When content changes, hash changes, filename changes, cache busted automatically
```

### Cache Busting with Fetch API

Same principle applied to API calls. Useful when API responses are incorrectly cached by browser or intermediate proxies.

```javascript
async function fetchFreshData(url) { const bustedUrl = `${url}?t=${Date.now()}`; const response = await fetch(bustedUrl); return response.json(); } // Usage: fetchFreshData('/api/user/profile').then(data => console.log(data));
```

### Server-Side: Ignoring Query Strings (Express.js Example)

Most static file servers (Express.static, Nginx, Apache, S3, CloudFront) ignore query strings by default when resolving file paths. This is what makes query string cache busting work.

```javascript
const express = require('express'); const app = express(); // Serve static files, ignoring query strings app.use(express.static('public', { // This is default behavior - query strings don't affect file lookup etag: true, lastModified: true, maxAge: '1y', immutable: true // Safe with filename hashing })); // Request to /image.jpg?v=123 serves /public/image.jpg
```

## Common Mistakes

- ❌ Using Date.now() for all static assets - defeats caching entirely, hurts performance
- ❌ Assuming query string busting works on all servers - some APIs process query strings
- ❌ Forgetting that CDNs/proxies may have different caching rules for query strings
- ❌ Not clearing browser cache during development - changes won't appear
- ❌ Using random values instead of deterministic ones - prevents any caching benefit
- ❌ Applying cache busting to HTML documents - breaks navigation caching, use Cache-Control headers instead

## Key Takeaways

- ✅ Browser caching improves performance but causes stale content issues when files change
- ✅ Cache busting works by changing the URL so browser fetches fresh, while server serves same file
- ✅ Query string parameters (?, &, =) are ignored by most static file servers - perfect for cache busting
- ✅ Timestamp-based busting (Date.now()) is simple but defeats caching on every load - use only for frequently changing content like user-uploaded thumbnails
- ✅ Version-based or content-hash-based busting is preferred for production - only invalidates when content actually changes
- ✅ Build tools (Webpack, Vite, esbuild) automate content-hash cache busting via filename hashing
- ✅ Always test cache busting in incognito/private mode or with dev tools 'Disable cache' checked

## Practice Questions

**Q1:** Why does adding a query string like ?v=123 to an image URL cause the browser to fetch a fresh copy?

<details><summary>Answer</summary>

Browsers use the complete URL (including query string) as the cache key. Changing the query string creates a new cache key, so the browser treats it as a different resource and makes a new network request. The server typically ignores query strings when resolving static file paths, so it serves the same underlying file.

</details>

**Q2:** What's the difference between using Date.now() and a version number for cache busting? When would you use each?

<details><summary>Answer</summary>

Date.now() generates a unique value on every page load, busting cache every time - useful for user-generated content that changes unpredictably (like uploaded thumbnails). A version number only changes when you increment it - preserves caching benefits for unchanged content, better for deployed assets (CSS, JS, images) where you control when changes happen.

</details>

**Q3:** Why might query string cache busting not work with some CDNs or proxy servers?

<details><summary>Answer</summary>

Some CDNs and proxy servers are configured to ignore query strings when caching (treating /image.jpg?v=1 and /image.jpg?v=2 as the same cache key) or to not cache URLs with query strings at all. This behavior varies by provider and configuration. Filename hashing (image.a1b2c3.jpg) works universally because the path itself changes.

</details>

**Q4:** In the lesson exercise, why do we append the timestamp to the thumbnail URL specifically, rather than all images on the page?

<details><summary>Answer</summary>

Thumbnails are user-uploaded content that can change while the URL stays the same. Other images (logos, UI icons) are deployed with the application and only change on new deployments. Cache busting everything with timestamps would defeat caching for static assets, hurting performance. Targeted cache busting preserves benefits where possible.

</details>

**Q5:** What is the 'gold standard' for cache busting in production applications, and why?

<details><summary>Answer</summary>

Content-hash-based filename versioning (e.g., main.a1b2c3d4.js) is the gold standard because: 1) Cache key changes only when file content changes, 2) Enables long-term caching with 'immutable' flag, 3) Works with all CDNs/proxies without special configuration, 4) Automated by build tools, 5) No query string edge cases.

</details>

## Concept Relationships

- **Browser Caching** → **Cache Busting**: Cache busting is a direct response to browser caching behavior - it exploits how cache keys work to force fresh fetches
- **Query String Cache Busting** → **Server Static File Handling**: Relies on servers ignoring query strings for static asset resolution - this is default behavior for most web servers and CDNs
- **Timestamp Cache Busting** → **Performance Trade-offs**: Defeats caching entirely - use only when content freshness is more critical than performance (user uploads)
- **Content Hash Cache Busting** → **Build Tools**: Requires build-time processing - Webpack, Vite, Rollup, esbuild all support this via plugins/configuration
- **Cache-Control Headers** → **Cache Busting**: Alternative/complementary approach - proper headers (max-age, immutable, must-revalidate) reduce need for aggressive cache busting

---
*Generated by Boot.dev Notes*