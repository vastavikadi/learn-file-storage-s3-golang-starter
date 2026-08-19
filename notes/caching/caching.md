# Caching Fundamentals - Browser Caching

> This lesson introduces the concept of caching in web applications, focusing on browser caching behavior. It explains how browsers store downloaded resources locally to avoid re-downloading on subsequent visits, and demonstrates a classic caching problem where updated content isn't reflected due to stale cache.

## Prerequisites
- Basic understanding of HTTP request/response cycle
- Familiarity with browser developer tools (Network tab)
- Understanding of static vs dynamic content

## Key Concepts

### Cache

A cache is temporary storage used to keep copies of frequently accessed data for faster retrieval. In web contexts, browsers cache resources (HTML, CSS, JS, images, videos) on the user's device after the first download.

**Examples:**
- Browser downloads boots-image-horizontal.png on first visit
- Browser stores image in local cache (memory/disk)
- On refresh, browser serves image from cache instead of network

> ⚠️ **Gotchas:**
> - Cache is not permanent - can be cleared by user, browser, or OS
> - Cache has size limits - older entries evicted when full
> - Different browsers have different cache storage mechanisms

### Browser Caching Behavior

On first visit to a web application, the browser must download all resources required to render the page. It then stores these resources locally. On subsequent visits or refreshes, the browser checks its cache first and serves cached copies when available, dramatically reducing load times and bandwidth usage.

**Examples:**
- First visit: 50 resources downloaded, 2MB transferred
- Second visit: 2 resources downloaded (HTML + changed CSS), 50KB transferred
- Cache hit: Resource served from local storage (near-instant)
- Cache miss: Resource fetched from network (slower)

> ⚠️ **Gotchas:**
> - Browser decides what to cache based on HTTP headers (Cache-Control, ETag, Last-Modified)
> - Some resources may not be cached if headers forbid it
> - Incognito/private mode typically disables persistent caching

### Cache Invalidation Problem

When a resource on the server changes but the URL remains the same, browsers may continue serving the stale cached version. This is one of the 'two hard things in computer science' (cache invalidation and naming things). The lesson exercise demonstrates this: uploading a new image with the same filename causes the browser to show the old cached image on refresh.

**Examples:**
- Upload boots-image-vertical.png replacing boots-image-horizontal.png
- Same URL: /video/123/image.png
- Browser sees same URL, serves cached horizontal image
- User sees wrong image until cache cleared or expires

> ⚠️ **Gotchas:**
> - Simple refresh (F5/Cmd+R) often uses cached version
> - Hard refresh (Ctrl+Shift+R / Cmd+Shift+R) bypasses cache
> - Cache-Control headers control this behavior but aren't foolproof
> - Users may not know to hard refresh

### Cache Hit vs Cache Miss

A cache hit occurs when the requested resource is found in the cache and can be served immediately. A cache miss occurs when the resource is not in cache (first visit, expired, evicted, or explicitly uncached) and must be fetched from the origin server.

**Examples:**
- Cache hit: Image loaded from disk in 5ms
- Cache miss: Image downloaded from server in 200ms
- Hit ratio: Percentage of requests served from cache

> ⚠️ **Gotchas:**
> - High hit ratio = better performance
> - First visit is always cache misses
> - Aggressive caching increases hit ratio but risks staleness

## Mental Models

### 💡 Scrooge McDuck's Money Bin

The lesson's analogy: browsers hoard resources like Scrooge hoards gold. They keep everything locally to avoid 'spending' network requests.

**Analogy:** Network request = spending money. Cache = keeping gold in your own vault. Browser wants to spend as little as possible.

### 💡 Pantry vs Grocery Store

Cache is your kitchen pantry. Origin server is the grocery store. First time cooking a recipe (first visit), you buy all ingredients. Next time, you check pantry first (cache hit) before going to store (cache miss).

**Analogy:** Pantry check = instant. Grocery trip = slow. But pantry might have expired ingredients (stale cache).

### 💡 Library Book Reserve System

Browser cache is like keeping a personal copy of a library book. You borrow once, keep it at home. When you need it again, you use your copy. But if the library updates the edition, your home copy is outdated.

**Analogy:** Borrowing = network request. Home copy = cache. New edition = server update. You won't know unless you check with library.

## Code Examples

### Cache-Control Header Examples

Cache-Control headers control browser caching behavior. These are set by the server in HTTP responses.

```http
# No caching - always revalidate with server
```

### Cache-Control: no-store

Prevents any caching. Browser must fetch from server every time. Use for sensitive data.

```http
Cache-Control: no-store
```

### Cache-Control: max-age

Cache for 1 year (31536000 seconds). 'public' allows CDN/proxy caching. Good for versioned static assets.

```http
Cache-Control: public, max-age=31536000
```

### Cache-Control: must-revalidate

Cache for 1 hour, then MUST check with server before using stale copy. Prevents serving stale content after expiry.

```http
Cache-Control: max-age=3600, must-revalidate
```

### ETag and Last-Modified Headers

ETag is a content hash. Browser sends If-None-Match on subsequent requests. Server returns 304 Not Modified if unchanged.

```http
ETag: W/abc123def456
```

### Cache Busting via Query String

Changing the query string forces browser to treat as new URL, bypassing cache. Simple but not ideal for CDNs.

```html
<img src='/video/123/image.png?v=2' alt='Boots image'>
```

### Cache Busting via Filename Hash

Best practice: include content hash in filename. When file changes, hash changes, URL changes, cache bypassed automatically.

```html
<img src='/video/123/image.a1b2c3d4.png' alt='Boots image'>
```

## Common Mistakes

- ❌ Assuming refresh always fetches new content - normal refresh uses cache
- ❌ Not implementing cache busting for updated static assets
- ❌ Setting Cache-Control: no-cache everywhere, losing all performance benefits
- ❌ Forgetting that CDNs and proxies also cache - not just browsers
- ❌ Using query strings for cache busting on CDNs that ignore query strings
- ❌ Not understanding difference between no-cache (revalidate) and no-store (don't store)
- ❌ Caching HTML documents aggressively - prevents updates from being seen

## Key Takeaways

- ✅ Cache = temporary local storage for faster subsequent access
- ✅ Browsers automatically cache resources on first visit
- ✅ Cache hits are near-instant; cache misses require network round-trip
- ✅ The cache invalidation problem: same URL + changed content = stale cache
- ✅ Hard refresh (Ctrl+Shift+R) bypasses browser cache
- ✅ Cache-Control headers are the primary mechanism for controlling caching
- ✅ Cache busting (changing URL when content changes) solves invalidation
- ✅ Filename hashing is the most robust cache busting strategy
- ✅ Balance caching aggressiveness with content freshness needs

## Practice Questions

**Q1:** In the lesson exercise, you upload boots-image-vertical.png to replace boots-image-horizontal.png, then refresh the browser. What happens and why?

<details><summary>Answer</summary>

The browser displays the old horizontal image because the URL hasn't changed. The browser sees the same URL, finds a valid cached copy, and serves it from cache (cache hit) instead of downloading the new vertical image. This is the classic cache invalidation problem.

</details>

**Q2:** What is the difference between a normal refresh (F5) and a hard refresh (Ctrl+Shift+R / Cmd+Shift+R)?

<details><summary>Answer</summary>

Normal refresh (F5) uses cached resources when valid (sends If-None-Match/If-Modified-Since). Hard refresh bypasses cache entirely, forcing re-download of all resources with Cache-Control: no-cache or similar headers.

</details>

**Q3:** Why is Cache-Control: no-store different from Cache-Control: no-cache?

<details><summary>Answer</summary>

no-store prevents ANY caching - browser must not write to cache at all. no-cache allows caching but requires revalidation with server before each use (sends conditional request). no-cache still stores the response; no-store does not.

</details>

**Q4:** What are two cache busting strategies, and which is better for production?

<details><summary>Answer</summary>

1) Query string: /image.png?v=2 - simple but some CDNs ignore query strings. 2) Filename hashing: /image.a1b2c3d4.png - better because URL path changes, works with all CDNs, enables long-term caching (max-age=1year) since new content = new URL.

</details>

**Q5:** If you set Cache-Control: public, max-age=31536000 on your HTML document, what problem will users experience?

<details><summary>Answer</summary>

Users will see stale HTML for up to a year. Even if you deploy updates, their browser won't request the new HTML until the year expires. HTML should typically have short max-age or no-cache with revalidation.

</details>

## Concept Relationships

- **Browser Cache** → **Cache-Control Headers**: Headers control browser caching behavior (max-age, no-store, must-revalidate)
- **Cache Invalidation** → **Cache Busting**: Cache busting (URL change) is the primary solution to invalidation problem
- **Browser Cache** → **CDN Cache**: Both cache resources; CDN sits between browser and origin, adds another cache layer
- **ETag/Last-Modified** → **Conditional Requests**: Validators enable 304 Not Modified responses for efficient revalidation

---
*Generated by Boot.dev Notes*