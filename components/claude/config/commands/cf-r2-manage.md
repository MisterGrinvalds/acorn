---
description: Manage CloudFlare R2 object storage
argument-hint: [action: list|create|upload|download]
allowed-tools: Read, Bash
---

## Task

Help the user manage CloudFlare R2 object storage buckets.

## Actions

Based on `$ARGUMENTS`:

### list
List R2 buckets:

```bash
wrangler r2 bucket list
cf_r2_buckets  # dotfiles function
wrr2list       # alias
```

### create
Create a new bucket:

```bash
wrangler r2 bucket create my-bucket
cf_r2_create my-bucket  # dotfiles function
```

### upload
Upload files to bucket:

```bash
# Upload single file
wrangler r2 object put my-bucket/path/file.txt --file=./local-file.txt

# Upload with content type
wrangler r2 object put my-bucket/image.png --file=./image.png --content-type="image/png"
```

### download
Download files:

```bash
wrangler r2 object get my-bucket/path/file.txt --file=./downloaded.txt
```

## Bucket Operations

### List Objects
```bash
wrangler r2 object list my-bucket
wrangler r2 object list my-bucket --prefix="uploads/"
```

### Delete Object
```bash
wrangler r2 object delete my-bucket/path/file.txt
```

### Delete Bucket
```bash
# Bucket must be empty
wrangler r2 bucket delete my-bucket
```

## Worker Integration

### wrangler.toml Binding
```toml
[[r2_buckets]]
binding = "MY_BUCKET"
bucket_name = "my-bucket"

# Preview bucket for local dev
[[r2_buckets]]
binding = "MY_BUCKET"
bucket_name = "my-bucket-dev"
preview_bucket_name = "my-bucket-dev"
```

### Worker Code
```typescript
export interface Env {
  MY_BUCKET: R2Bucket;
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    // GET - Download
    if (request.method === "GET") {
      const object = await env.MY_BUCKET.get(url.pathname.slice(1));
      if (!object) return new Response("Not found", { status: 404 });
      return new Response(object.body, {
        headers: { "Content-Type": object.httpMetadata?.contentType || "" },
      });
    }

    // PUT - Upload
    if (request.method === "PUT") {
      await env.MY_BUCKET.put(url.pathname.slice(1), request.body, {
        httpMetadata: { contentType: request.headers.get("Content-Type") || "" },
      });
      return new Response("Uploaded", { status: 201 });
    }

    return new Response("Method not allowed", { status: 405 });
  },
};
```

## R2 API Operations

| Method | Description |
|--------|-------------|
| `put(key, value, options)` | Upload object |
| `get(key)` | Download object |
| `head(key)` | Get object metadata |
| `delete(key)` | Delete object |
| `list(options)` | List objects |

## S3 Compatibility

R2 is S3-compatible. To use with S3 tools:

```bash
# Generate S3 API credentials in dashboard
# Then use with aws-cli:
aws s3 ls s3://my-bucket \
  --endpoint-url https://<account-id>.r2.cloudflarestorage.com
```

## Use Cases

- **Static assets**: Images, videos, downloads
- **Backups**: Database backups, logs
- **User uploads**: File storage for apps
- **CDN origin**: Origin for cached content

## Best Practices

1. Use meaningful bucket names
2. Organize with prefixes (folders)
3. Set appropriate content types
4. Use presigned URLs for secure uploads
5. Enable access logging for auditing

## Dotfiles Integration

- `cf_r2_buckets` - List R2 buckets
- `cf_r2_create <bucket>` - Create bucket
- `wrr2` - wrangler r2 (alias)
- `wrr2list` - wrangler r2 bucket list (alias)
