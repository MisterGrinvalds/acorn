---
description: Manage CloudFlare KV key-value storage
argument-hint: [action: list|create|get|put|delete]
allowed-tools: Read, Bash
---

## Task

Help the user manage CloudFlare KV key-value namespaces.

## Actions

Based on `$ARGUMENTS`:

### list
List KV namespaces:

```bash
wrangler kv namespace list
cf_kv_namespaces  # dotfiles function
wrkvlist          # alias
```

### create
Create a namespace:

```bash
wrangler kv namespace create MY_CACHE
cf_kv_create MY_CACHE  # dotfiles function

# Create preview namespace (for local dev)
wrangler kv namespace create MY_CACHE --preview
```

### get
Get a value:

```bash
wrangler kv key get --namespace-id=<id> "my-key"

# From binding name
wrangler kv key get --binding=MY_CACHE "my-key"
```

### put
Set a value:

```bash
wrangler kv key put --namespace-id=<id> "my-key" "my-value"

# With TTL (seconds)
wrangler kv key put --namespace-id=<id> "my-key" "my-value" --ttl=3600

# With expiration timestamp
wrangler kv key put --namespace-id=<id> "my-key" "my-value" --expiration=1700000000
```

### delete
Delete a key:

```bash
wrangler kv key delete --namespace-id=<id> "my-key"
```

## List Keys in Namespace

```bash
wrangler kv key list --namespace-id=<id>
wrangler kv key list --namespace-id=<id> --prefix="user:"
```

## Bulk Operations

```bash
# Bulk put from JSON file
wrangler kv bulk put --namespace-id=<id> ./data.json

# data.json format:
# [
#   {"key": "key1", "value": "value1"},
#   {"key": "key2", "value": "value2", "expiration_ttl": 3600}
# ]

# Bulk delete
wrangler kv bulk delete --namespace-id=<id> ./keys.json
# ["key1", "key2", "key3"]
```

## Worker Integration

### wrangler.toml Binding
```toml
[[kv_namespaces]]
binding = "MY_CACHE"
id = "abc123..."
preview_id = "def456..."  # For local dev
```

### Worker Code
```typescript
export interface Env {
  MY_CACHE: KVNamespace;
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const key = url.pathname.slice(1);

    // GET - Read
    if (request.method === "GET") {
      const value = await env.MY_CACHE.get(key);
      if (!value) return new Response("Not found", { status: 404 });
      return new Response(value);
    }

    // PUT - Write
    if (request.method === "PUT") {
      const value = await request.text();
      await env.MY_CACHE.put(key, value, {
        expirationTtl: 3600, // 1 hour
      });
      return new Response("Stored", { status: 201 });
    }

    // DELETE
    if (request.method === "DELETE") {
      await env.MY_CACHE.delete(key);
      return new Response("Deleted");
    }

    return new Response("Method not allowed", { status: 405 });
  },
};
```

## KV API Operations

| Method | Description |
|--------|-------------|
| `get(key, options)` | Read value |
| `getWithMetadata(key)` | Read value with metadata |
| `put(key, value, options)` | Write value |
| `delete(key)` | Delete key |
| `list(options)` | List keys |

### Get Options
```typescript
// Get as text (default)
const text = await env.KV.get("key");

// Get as JSON
const json = await env.KV.get("key", { type: "json" });

// Get as ArrayBuffer
const buffer = await env.KV.get("key", { type: "arrayBuffer" });

// Get as stream
const stream = await env.KV.get("key", { type: "stream" });
```

### Put Options
```typescript
await env.KV.put("key", "value", {
  expirationTtl: 3600,     // Seconds until expiration
  expiration: 1700000000,   // Unix timestamp
  metadata: { foo: "bar" }, // Custom metadata
});
```

## Use Cases

- **Caching**: API responses, computed values
- **Configuration**: Feature flags, settings
- **Sessions**: User session data
- **Rate limiting**: Request counters
- **A/B testing**: Experiment assignments

## Best Practices

1. Use meaningful key prefixes: `user:123`, `cache:api:users`
2. Set appropriate TTLs for cache data
3. Use JSON for structured data
4. Keep values under 25MB limit
5. Use metadata for key attributes

## Dotfiles Integration

- `cf_kv_namespaces` - List KV namespaces
- `cf_kv_create <namespace>` - Create namespace
- `wrkv` - wrangler kv (alias)
- `wrkvlist` - wrangler kv namespace list (alias)
