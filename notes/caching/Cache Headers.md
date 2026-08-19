# Cache Headers - HTTP Caching Control

> This lesson covers the Cache-Control HTTP header and its directives for controlling browser and intermediary caching behavior. It explains common directives (no-store, max-age, stale-while-revalidate, no-cache) and demonstrates implementing a no-cache middleware in Go.

## Prerequisites
- Basic HTTP request/response cycle
- Go HTTP handler/middleware pattern
- Understanding of browser caching behavior

## Key Concepts

### Cache-Control Header

The Cache-Control header is the primary mechanism for servers to communicate caching policies to clients (browsers) and intermediate caches (CDNs, proxies). It replaces older headers like Expires and Pragma. When a server includes Cache-Control in a response, it instructs the client how to handle caching of that specific resource.

**Examples:**
- Cache-Control: no-store
- Cache-Control: max-age=3600
- Cache-Control: stale-while-revalidate=60
- Cache-Control: no-cache

> ⚠️ **Gotchas:**
> - Cache-Control is a response header (sent by server), not a request header
> - Browsers generally respect Cache-Control but are not strictly required to
> - Multiple directives can be combined: Cache-Control: no-cache, max-age=3600

### no-store Directive

The most restrictive caching directive. It instructs caches (browser, CDN, proxy) to NOT store any part of the response. The response must be fetched from the origin server every time. Use for sensitive data (auth tokens, personal info) or rapidly changing content where stale data is unacceptable.

**Examples:**
- Cache-Control: no-store
- Cache-Control: no-store, private  // private reinforces no shared caching

> ⚠️ **Gotchas:**
> - no-store prevents caching entirely - no disk cache, no memory cache
> - Even back/forward navigation may re-fetch the resource
> - Does not prevent the browser from displaying the content initially

### max-age Directive

Specifies the maximum time (in seconds) a response can be considered fresh. After this time expires, the cache must revalidate with the origin server before serving the cached copy. max-age=3600 means cache for 1 hour. This is the most common directive for static assets.

**Examples:**
- Cache-Control: max-age=3600        // 1 hour
- Cache-Control: max-age=86400       // 1 day
- Cache-Control: max-age=31536000    // 1 year (for versioned static assets)
- Cache-Control: public, max-age=3600  // public allows CDN caching

> ⚠️ **Gotchas:**
> - max-age is relative to response time, not absolute
> - Clock skew between client and server can affect expiration
> - max-age=0 is effectively 'must revalidate' but still allows caching

### no-cache Directive

Despite the name, no-cache DOES allow caching. It means 'cache this response, but you MUST revalidate with the origin server before serving the cached copy.' The cache stores the response but checks with the server (via ETag or Last-Modified) on every request. If the server returns 304 Not Modified, the cached version is used.

**Examples:**
- Cache-Control: no-cache
- Cache-Control: no-cache, max-age=3600  // cache for 1hr but always revalidate

> ⚠️ **Gotchas:**
> - The name is misleading - it does NOT mean 'don't cache'
> - Requires server to support conditional requests (ETag/Last-Modified)
> - Adds latency on every request due to revalidation round-trip
> - Use when you want caching benefits but need freshness guarantees

### stale-while-revalidate Directive

Allows serving stale (expired) content while asynchronously revalidating in the background. The value is seconds. During this window, if a request comes in and the cache is stale, the stale content is served immediately while a background revalidation fetch occurs. Improves perceived performance by avoiding blocking on revalidation.

**Examples:**
- Cache-Control: max-age=60, stale-while-revalidate=300  // fresh 60s, serve stale up to 5min while revalidating
- Cache-Control: stale-while-revalidate=86400  // serve stale up to 1 day while revalidating

> ⚠️ **Gotchas:**
> - Only works in browsers that support it (modern browsers do)
> - Requires max-age to be set (defines when content becomes stale)
> - Background revalidation may still serve stale content to subsequent requests until complete
> - Not supported by all CDNs/proxies

### Middleware Pattern for Cache Headers in Go

HTTP middleware in Go wraps handlers to add cross-cutting concerns like caching headers. The middleware sets headers before calling the next handler. This centralizes cache policy and ensures consistent application across routes.

**Examples:**
- // Middleware that sets no-store on all responses
- func noCacheMiddleware(next http.Handler) http.Handler {
-     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
-         w.Header().Set('Cache-Control', 'no-store')
-         next.ServeHTTP(w, r)
-     })
- }

> ⚠️ **Gotchas:**
> - Header must be set BEFORE calling next.ServeHTTP()
> - w.Header().Set() overwrites; use w.Header().Add() for multiple values
> - Middleware order matters - cache middleware should wrap the handler

## Mental Models

### 💡 Cache-Control as a Contract

Think of Cache-Control as a contract between server and client. The server says 'here's how you may use this response.' The client (browser) agrees to follow the rules. Directives are clauses in this contract: no-store = 'don't keep a copy', max-age = 'this copy expires at X', no-cache = 'keep a copy but ask permission before using it.'

### 💡 The Library Book Analogy

no-store: Reference book - must read in library, cannot check out. max-age: Regular book - check out for 2 weeks (3600 seconds). no-cache: Book you can check out, but must call library before each re-read to confirm it's still current. stale-while-revalidate: You can keep reading your expired copy while the library checks for a new edition in the background.

### 💡 Freshness vs Validation

Two-phase caching model: 1) Freshness (max-age) - how long until the cached copy is considered 'stale'. 2) Validation (no-cache, must-revalidate) - what happens when stale. Freshness avoids network requests entirely. Validation makes a conditional request (If-None-Match) - cheap if 304 Not Modified, full download if changed.

## Code Examples

### Go Middleware: no-store for All Responses

This middleware ensures every response from the wrapped handlers includes Cache-Control: no-store, preventing any caching of sensitive API responses.

```go
package mainimport (    'net/http')// noCacheMiddleware sets Cache-Control: no-store on all responsesfunc noCacheMiddleware(next http.Handler) http.Handler {    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        // Set header BEFORE calling next handler        w.Header().Set('Cache-Control', 'no-store')        next.ServeHTTP(w, r)    })}func main() {    mux := http.NewServeMux()    mux.HandleFunc('/api/data', handleData)    // Wrap all routes with no-cache middleware    wrappedMux := noCacheMiddleware(mux)    http.ListenAndServe(':8080', wrappedMux)}func handleData(w http.ResponseWriter, r *http.Request) {    w.Write([]byte('sensitive data'))}
```

### Go Middleware: max-age for Static Assets

Conditional middleware that applies long-term caching (1 year) only to static assets. The 'immutable' directive tells caches the content will never change, avoiding revalidation entirely.

```go
func staticCacheMiddleware(next http.Handler) http.Handler {    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        // Only apply to static asset paths        if strings.HasPrefix(r.URL.Path, '/static/') {            w.Header().Set('Cache-Control', 'public, max-age=31536000, immutable')        }        next.ServeHTTP(w, r)    })}
```

### Go Middleware: no-cache with ETag Support

no-cache requires ETag or Last-Modified for revalidation. This middleware adds both the directive and an ETag generator.

```go
func mustRevalidateMiddleware(next http.Handler) http.Handler {    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        w.Header().Set('Cache-Control', 'no-cache')        w.Header().Set('ETag', generateETag(r.URL.Path))        next.ServeHTTP(w, r)    })}func generateETag(path string) string {    // In practice, compute hash of file content    return fmt.Sprintf('W/%%22%x%%22', hashContent(path))}
```

### JavaScript: Reverting to Direct Image Source Assignment

The assignment reverts client-side cache busting (query string) in favor of server-controlled caching via Cache-Control headers.

```javascript
// Before (cache busting with query string)function viewVideo(video) {    thumbnailImg.src = video.thumbnail_url + '?t=' + Date.now();}// After (relying on server Cache-Control headers)function viewVideo(video) {    thumbnailImg.src = video.thumbnail_url;}
```

### Combined Cache-Control Directives

Real-world Cache-Control headers often combine multiple directives for precise control.

```http
# Multiple directives separated by commasCache-Control: no-cache, max-age=3600, must-revalidateCache-Control: public, max-age=3600, stale-while-revalidate=86400Cache-Control: private, no-store, max-age=0# private: only browser cache (not CDN/proxy)# public: any cache may store# must-revalidate: strict revalidation after expiry# immutable: content never changes (use with long max-age)
```

## Common Mistakes

- ❌ Confusing no-cache with no-store - no-cache allows caching but requires revalidation; no-store forbids caching entirely
- ❌ Setting Cache-Control on requests instead of responses - it's a response header
- ❌ Forgetting that max-age=0 with no-cache still caches the response (just forces revalidation)
- ❌ Not setting ETag/Last-Modified when using no-cache - revalidation can't work without validators
- ❌ Applying no-store to static assets that should be cached - hurts performance unnecessarily
- ❌ Assuming all intermediaries (CDNs, corporate proxies) respect Cache-Control - some strip or ignore headers
- ❌ Using stale-while-revalidate without max-age - it has no effect without a freshness lifetime

## Key Takeaways

- ✅ Cache-Control is the primary HTTP header for controlling caching behavior
- ✅ no-store = never cache (most restrictive); use for sensitive data
- ✅ max-age=N = cache for N seconds; use for static assets with versioned filenames
- ✅ no-cache = cache but revalidate every time; misleading name, requires ETag/Last-Modified
- ✅ stale-while-revalidate = serve stale content while background revalidation occurs; improves perceived performance
- ✅ Middleware pattern in Go allows centralized cache policy application across routes
- ✅ Server-controlled caching via headers is superior to client-side cache busting (query strings) because it works for all clients, not just your JavaScript

## Practice Questions

**Q1:** What is the difference between Cache-Control: no-cache and Cache-Control: no-store?

<details><summary>Answer</summary>

no-store prohibits caching entirely - the response must not be stored in any cache. no-cache allows caching but requires revalidation with the origin server before each use. Despite the name, no-cache does NOT mean 'don't cache'.

</details>

**Q2:** If a response has Cache-Control: max-age=3600, stale-while-revalidate=300, how long can a browser serve a stale response while revalidating in the background?

<details><summary>Answer</summary>

300 seconds (5 minutes). After max-age expires (1 hour), the response becomes stale. For the next 300 seconds, the browser can serve the stale response immediately while triggering a background revalidation. After 300 seconds, the browser must wait for revalidation before serving.

</details>

**Q3:** Why must Cache-Control headers be set before calling next.ServeHTTP() in Go middleware?

<details><summary>Answer</summary>

Because the next handler writes the response body and headers. Once the response is written (or headers are flushed), you cannot modify headers. Setting Cache-Control before ensures it's included in the final response.

</details>

**Q4:** A server returns Cache-Control: no-cache but no ETag or Last-Modified header. What happens on subsequent requests?

<details><summary>Answer</summary>

The browser will still cache the response but cannot perform conditional revalidation (no validator to send in If-None-Match/If-Modified-Since). It will likely make unconditional requests every time, downloading the full response each time - defeating the purpose of no-cache.

</details>

**Q5:** When should you use 'private' vs 'public' in Cache-Control?

<details><summary>Answer</summary>

Use 'private' for user-specific content (dashboard, profile pages) that should only be cached in the user's browser, not in shared caches (CDNs, corporate proxies). Use 'public' for truly static, non-personalized assets (CSS, JS, images) that can be cached by any intermediary.

</details>

## Concept Relationships

- **Cache-Control** → **ETag/Last-Modified**: no-cache and must-revalidate require validators (ETag or Last-Modified) to perform conditional revalidation efficiently
- **max-age** → **stale-while-revalidate**: stale-while-revalidate only applies after max-age expires; it defines the stale-serving window
- **Middleware pattern** → **Cache-Control**: Go middleware is the idiomatic way to apply consistent Cache-Control headers across multiple routes
- **Query string cache busting** → **Cache-Control headers**: Cache-Control headers replace client-side cache busting techniques; server-controlled caching is more robust

---
*Generated by Boot.dev Notes*