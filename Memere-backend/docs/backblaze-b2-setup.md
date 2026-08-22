# Backblaze B2 as the object store

This backend can use a **Backblaze B2 bucket** as the storage layer for videos,
PDFs, and images. B2 exposes an **S3-compatible API**, so the existing S3 client
is reused unchanged — you only point it at B2's endpoint and hand it B2
application-key credentials. Upload, download, presign, and delete logic are
identical to the AWS S3 path.

> **Secrets policy.** `B2_APPLICATION_KEY_ID` and `B2_APPLICATION_KEY` are
> server-side only. Never put real values in source, `.env.example`, logs, or
> Git. Put real values only in your untracked `.env` or a secret manager. The
> backend never returns or logs them.

---

## 1. Create a bucket

1. Sign in at <https://www.backblaze.com/> → **B2 Cloud Storage → Buckets**.
2. **Create a Bucket.** Give it a globally-unique name (this is `B2_BUCKET`).
3. Keep **Files in Bucket are: Private**. The app never serves B2 directly —
   every download flows through a short-lived, backend-signed URL, so the bucket
   must stay private.

## 2. Note the endpoint and region

Open the bucket; its details show an **Endpoint** like:

```
s3.us-west-004.backblazeb2.com
```

- `B2_ENDPOINT` = that endpoint (scheme optional; the backend prepends `https://`).
- `B2_REGION`  = the region segment in the middle, e.g. `us-west-004`. It must
  match the endpoint or SigV4 signing fails.

## 3. Create an application key

1. **Account → Application Keys → Add a New Application Key.**
2. Restrict it to the bucket from step 1 (least privilege) and grant **Read and
   Write**.
3. On create, B2 shows two values **once**:
   - **keyID**        → `B2_APPLICATION_KEY_ID`
   - **applicationKey** → `B2_APPLICATION_KEY` (copy it now; it is never shown again)

## 4. Fill in your `.env`

```env
STORAGE_PROVIDER=b2
B2_APPLICATION_KEY_ID=<paste real keyID>
B2_APPLICATION_KEY=<paste real applicationKey>
B2_ENDPOINT=s3.us-west-004.backblazeb2.com
B2_REGION=us-west-004
B2_BUCKET=<your bucket name>
```

(In `.env.example` these stay as `YOUR_B2_...` placeholders.)

## 5. Run

```bash
make run
```

Uploads and playback now use B2. Nothing else in the pipeline changes — the CDN
signer sits in front of the store as before; with `CDN_DOMAIN` empty it delegates
to B2 pre-signed URLs.

---

## Why the code disables request checksums

`aws-sdk-go-v2` (v1.32+) adds a CRC32 integrity **trailer** to uploads by
default. Backblaze B2 rejects those trailers, so the B2 branch constructs the S3
client with:

```go
o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
```

This also keeps pre-signed PUT URLs clean (no checksum headers the mobile client
would have to reproduce). AWS S3 and MinIO leave this on (`DisableChecksums:false`).

## How it works at runtime

| Concern | Behavior |
|---|---|
| **Video upload (presigned)** | `POST /lessons/:id/videos/upload-url` returns a pre-signed B2 PUT URL; the client uploads the bytes directly. |
| **Video upload (proxied)** | `POST /lessons/:id/videos/upload` streams the file through the backend to B2 via `PutObject`. |
| **Playback / download** | The stream/download endpoints return short-lived B2 pre-signed GET URLs (default 2h). |
| **Authorization** | Every non-public request is authorized by the app first; a signed URL is minted only after the access check passes. |

## Browser direct uploads — bucket CORS (live upload %)

The admin panel uploads a video **directly from the browser to B2** using the
pre-signed PUT URL. That direct transfer is the only path that reports a **true,
live upload percentage** (the browser measures the bytes as they leave). If the
browser's cross-origin `PUT` is blocked, the admin silently falls back to the
Next.js **server proxy**: the file is uploaded to the Node server first and then
to B2 server-side — the browser can only measure the browser→server leg, so its
progress bar can't reflect the (often slower) server→B2 leg.

A cross-origin `PUT`/`GET` from the admin origin to `*.backblazeb2.com` requires
**CORS rules on the bucket**. B2's S3-compatible API does **not** implement
`PutBucketCors`, so this is configured with B2's own tools, not the S3 client:

**Option A — B2 web console.** Open the bucket → **CORS Rules** → add a rule that
allows your admin origin(s) to perform `s3_put`, `s3_get`, and `s3_head`.

**Option B — `b2` CLI.** Save the rules to `cors.json`:

```json
[
  {
    "corsRuleName": "memereAdminBrowserUploads",
    "allowedOrigins": ["https://admin.memere.et", "http://localhost:3000"],
    "allowedOperations": ["s3_put", "s3_get", "s3_head"],
    "allowedHeaders": ["*"],
    "exposeHeaders": ["etag"],
    "maxAgeSeconds": 3600
  }
]
```

then apply them (use whichever form your installed CLI supports):

```bash
# b2 CLI v3+
b2 bucket update --cors-rules "$(cat cors.json)" <your-bucket-name>
# older b2 CLI
b2 update-bucket --corsRules "$(cat cors.json)" <your-bucket-name> allPrivate
```

- List **exact** origins — the real admin domain plus `http://localhost:3000`
  for local dev. Do **not** use `"*"` (it mirrors the API's fail-closed
  `CORS_ALLOWED_ORIGINS` policy and keeps the private bucket locked down).
- The bucket stays **Private**; CORS only governs which browser origins may use a
  URL the backend already signed — it does not make objects public.
- Playback/download also benefit: `s3_get`/`s3_head` let the browser range-request
  signed media without a proxy hop.

Without these rules the app still works — it just falls back to the proxy and the
admin's percentage is only approximate for large videos.

## Switching providers

- Back to **AWS S3**: `STORAGE_PROVIDER=s3`, leave `S3_ENDPOINT` empty, set the
  `AWS_*` keys.
- **MinIO** (local dev): `STORAGE_PROVIDER=minio`, `S3_ENDPOINT=http://localhost:9000`,
  `S3_USE_PATH_STYLE=true`, keys `minioadmin`.
- **Google Drive**: `STORAGE_PROVIDER=gdrive` (see `docs/google-drive-setup.md`).

The storage seam (`service.ObjectStore`) is swappable without touching usecases.
