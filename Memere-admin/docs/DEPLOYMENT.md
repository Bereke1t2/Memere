# Memere Admin — Deployment Runbook

## Environment contract

| Variable | Required | Scope | Notes |
|---|---|---|---|
| `API_BASE_URL` | Yes | Server-only | Go backend base URL (e.g. `https://api.memere.app`) |
| `COOKIE_SECRET` | Yes | Server-only | ≥ 32 random bytes — signs httpOnly session cookies |
| `NODE_ENV` | Yes | Server-only | Set to `production` |
| `NEXT_PUBLIC_APP_NAME` | No | Public | Display name shown in the browser tab |

Generate a strong secret:
```bash
openssl rand -hex 32
```

**Never set `NEXT_PUBLIC_API_BASE_URL` or expose `COOKIE_SECRET` to the browser.**

---

## CORS coordination with the backend

All Go API calls originate **server-side** (Next.js Route Handlers and Server
Components call the backend from the Node.js process — never from the browser).
The browser never contacts the Go backend directly.

**CORS action required: none for standard deploys.**

If you ever add a direct browser→backend call (strongly discouraged — violates
the token-safety architecture), you would need to add the admin origin to the
backend's allowed origins:

```yaml
# k8s/configmaps/app-config.yaml
CORS_ALLOWED_ORIGINS: "https://admin.memere.app"
```

Confirm the backend is reachable from the admin host:
```bash
curl -sI https://api.memere.app/api/v1/health
```

---

## Path A — Vercel (recommended)

1. **Import** the `Memere-admin` repo in [vercel.com](https://vercel.com).

2. **Framework preset**: Next.js (auto-detected).

3. **Environment variables** — add in the Vercel dashboard under
   *Settings → Environment Variables* (mark as **Production** only):

   | Name | Value |
   |---|---|
   | `API_BASE_URL` | `https://api.memere.app` |
   | `COOKIE_SECRET` | `$(openssl rand -hex 32)` |
   | `NODE_ENV` | `production` |
   | `NEXT_PUBLIC_APP_NAME` | `Memere Admin` |

4. **Deploy**: push to `main` — Vercel builds and deploys automatically.

5. **Custom domain**: in *Settings → Domains* add `admin.memere.app`.
   Point a CNAME at `cname.vercel-dns.com` in your DNS provider.

6. **Verify**:
   ```bash
   curl -sI https://admin.memere.app/login | grep -iE 'content-security|x-frame'
   ```

---

## Path B — Self-hosted Docker + nginx + TLS

### 1. Build and push the image

```bash
# On CI or your machine
make docker-build
docker tag memere-admin:local registry.example.com/memere-admin:latest
docker push registry.example.com/memere-admin:latest
```

### 2. Create the env file on the server

```bash
# /srv/memere-admin/.env.production  (chmod 600, root-owned)
API_BASE_URL=https://api.memere.app
COOKIE_SECRET=<output of: openssl rand -hex 32>
NODE_ENV=production
NEXT_PUBLIC_APP_NAME=Memere Admin
```

### 3. Docker Compose

```yaml
# /srv/memere-admin/docker-compose.yml
version: "3.9"
services:
  admin:
    image: registry.example.com/memere-admin:latest
    restart: unless-stopped
    env_file: .env.production
    expose:
      - "3000"
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/login"]
      interval: 30s
      timeout: 5s
      retries: 3
```

```bash
docker compose up -d
```

### 4. nginx reverse proxy + TLS (Certbot)

```nginx
# /etc/nginx/sites-available/admin.memere.app
server {
    listen 443 ssl http2;
    server_name admin.memere.app;

    ssl_certificate     /etc/letsencrypt/live/admin.memere.app/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/admin.memere.app/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass         http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header   Upgrade $http_upgrade;
        proxy_set_header   Connection 'upgrade';
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}

server {
    listen 80;
    server_name admin.memere.app;
    return 301 https://$host$request_uri;
}
```

```bash
# Obtain certificate
certbot --nginx -d admin.memere.app
nginx -t && systemctl reload nginx
```

### 5. Verify

```bash
curl -sI https://admin.memere.app/login | grep -iE 'content-security|x-frame|set-cookie'
```

Expected: security headers present, no `Set-Cookie` on the login page (cookies
are only set after a successful `POST /api/auth/login`).

---

## Rollback

**Vercel**: use the *Deployments* tab → "Redeploy" a previous build.

**Self-hosted**:
```bash
docker compose pull
docker compose up -d --no-deps admin
# or pin an image tag and re-deploy
```
