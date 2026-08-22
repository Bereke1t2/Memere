# Google Drive as the centralized file store

This backend can use **one admin-owned Google Drive account** as the storage
layer for **PDFs, images (course thumbnails), and videos**, instead of S3/MinIO.

Key properties of this design (do not weaken these):

- **Students and teachers never sign in with Google.** Only the backend talks to
  Drive, using one admin account's OAuth **refresh token**.
- **The Flutter app never receives any Google credential.** It only ever gets
  normal app JWTs and short-lived, backend-signed URLs.
- **Drive is never made public.** Every file request is authorized by the app
  first; bytes are then streamed *through* the backend (or via a short-lived,
  HMAC-signed `/media` proxy URL). No Drive "anyone with the link" sharing.
- The Postgres database stays the source of truth; a small `storage.objects`
  table maps stable object keys to opaque Drive file IDs.

> **Secrets policy.** `GOOGLE_CLIENT_SECRET`, the refresh token, and any access
> token are server-side only. Never put real values in source, `.env.example`,
> logs, or Git. Put real values only in your untracked `.env` or a secret
> manager. The backend never returns or logs them.

---

## 1. Create a Google Cloud project

1. Go to <https://console.cloud.google.com/> and create (or pick) a project,
   e.g. `memere-storage`.
2. **APIs & Services → Library →** search **"Google Drive API" → Enable**.

## 2. Configure the OAuth consent screen

1. **APIs & Services → OAuth consent screen.**
2. User type: **External** is fine (you will be the only user). Fill in the app
   name, your support email, and developer email.
3. **Scopes:** you don't have to pre-add them here, but the app will request:
   - `https://www.googleapis.com/auth/drive` (create/read/delete files it owns)
   - `openid`, `email` (only to record *which* account connected, for the admin UI)
4. **Test users:** add the **admin Google account** that will own the Drive
   files. (While the app is in "Testing", only listed test users can connect —
   that's exactly what we want. Refresh tokens for a Testing app can expire after
   7 days; if you want a long-lived token, click **Publish app** once configured.)

## 3. Create the OAuth client credentials

1. **APIs & Services → Credentials → Create Credentials → OAuth client ID.**
2. Application type: **Web application.**
3. **Authorized redirect URIs — add exactly the callback URL the backend uses:**
   - Local dev: `http://localhost:8080/api/v1/admin/google/callback`
   - Production: `https://YOUR_API_DOMAIN/api/v1/admin/google/callback`
   This must match `GOOGLE_REDIRECT_URI` **character-for-character**.
4. Create. Copy the **Client ID** and **Client secret** into your `.env`:

   ```env
   STORAGE_PROVIDER=gdrive
   GOOGLE_CLIENT_ID=<paste real client id>
   GOOGLE_CLIENT_SECRET=<paste real client secret>
   GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/admin/google/callback
   APP_PUBLIC_URL=http://localhost:8080
   ```

   (In `.env.example` these stay as `YOUR_GOOGLE_CLIENT_ID` etc.)

## 4. Choose the Drive root folder

1. In the **admin account's** Google Drive, create a folder, e.g. `Memere`.
2. Open it; the URL is `https://drive.google.com/drive/folders/<FOLDER_ID>`.
3. Copy `<FOLDER_ID>` into `.env`:

   ```env
   GOOGLE_DRIVE_ROOT_FOLDER_ID=<paste real folder id>
   ```

All objects are created inside this folder.

## 5. Run the database migration

The Drive backend needs two tables (object index + encrypted credential row):

```bash
make migrate-up   # applies migrations/0023_google_drive_storage.up.sql
```

## 6. Connect the admin account (one time)

The refresh token is obtained **once**, through an admin-only page — you never
paste it anywhere.

1. Start the backend (`make run`) with the env vars above set.
2. As an **admin** user of the app, obtain an admin JWT and call:

   ```
   GET /api/v1/admin/google/connect
   Authorization: Bearer <admin access token>
   ```

   It returns `{ "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?..." }`.
3. Open that `auth_url` in a browser, sign in as the **admin Google account**
   (the test user from step 2), and grant access. Google redirects to
   `GOOGLE_REDIRECT_URI`; the backend exchanges the code, encrypts the refresh
   token (AES-256-GCM, key derived from `JWT_SECRET`), stores it in
   `storage.google_credentials`, and shows a plain "Connected" page.
4. Verify:

   ```
   GET /api/v1/admin/google/status   →   { "connected": true }
   ```

That's it — uploads and playback now use Drive. Re-running the connect flow
replaces the stored token (e.g. after a revoke).

> **Alternative (no connect page):** if you already have a refresh token from a
> secret manager, set `GOOGLE_REFRESH_TOKEN` in the environment. The DB value
> (from the connect flow) takes precedence when both are present.

---

## How it works at runtime

| Concern | Behavior |
|---|---|
| **PDF upload** | `POST /api/v1/lessons/:id/pdf` (teacher/admin) streams the file to Drive under `lessons/<id>/notes.pdf`. |
| **PDF download** | `GET /api/v1/lessons/:id/pdf` streams the bytes back through the backend. |
| **Video upload** | `POST /api/v1/lessons/:id/videos/upload` (teacher/admin, multipart `file`) streams the MP4 straight to Drive; the lesson video is marked ready with **no transcoding**. |
| **Video playback** | The existing stream/download endpoints return a short-lived, HMAC-signed `/api/v1/media?key=…&exp=…&sig=…` URL. The player streams it with HTTP **Range** so seeking works. |
| **Thumbnail upload** | `POST /api/v1/courses/:id/thumbnail` (teacher/admin, multipart `file`) stores the image and sets the course thumbnail URL. |
| **Thumbnail serve** | `GET /api/v1/courses/:id/thumbnail` is public (marketing content), streamed through the backend and cacheable. |
| **Authorization** | Every non-public request is authorized by the app first; the signed `/media` URL is minted **only after** the video access check passes, and expires quickly. |

### Bandwidth & quota notes

Because video is proxied through the backend (not served from a CDN), egress
flows through your server. Google Drive also enforces a per-account **download
cap (~750 GB/day)**. For a large audience, plan to move hot video to a CDN later;
the storage seam (`service.ObjectStore`) is swappable without touching usecases.

## Switching back to S3/MinIO

Set `STORAGE_PROVIDER=s3` (or `minio`) and restore the `AWS_*`/`S3_*` keys. The
Drive-only routes (`/media`, `/admin/google/*`) simply aren't registered.
