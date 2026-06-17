import { NextRequest, NextResponse } from "next/server";

const ACCESS_COOKIE = "mm_access";
const REFRESH_COOKIE = "mm_refresh";

function csrfError() {
  return new NextResponse(
    JSON.stringify({ message: "Cross-site request blocked." }),
    { status: 403, headers: { "Content-Type": "application/json" } }
  );
}

export function proxy(req: NextRequest) {
  const access = req.cookies.get(ACCESS_COOKIE)?.value;
  const refresh = req.cookies.get(REFRESH_COOKIE)?.value;

  if (!access && !refresh) {
    const loginUrl = new URL("/login", req.url);
    loginUrl.searchParams.set("from", req.nextUrl.pathname);
    return NextResponse.redirect(loginUrl);
  }

  // CSRF: for state-changing requests, verify the Origin header matches our host.
  // Same-origin fetch() may omit Origin — we allow that and rely on SameSite=Lax
  // cookies as the primary backstop. If Origin IS present it must match our host.
  const method = req.method.toUpperCase();
  if (method !== "GET" && method !== "HEAD" && method !== "OPTIONS") {
    const origin = req.headers.get("origin");
    if (origin) {
      const host = req.headers.get("host");
      try {
        if (new URL(origin).host !== host) return csrfError();
      } catch {
        return csrfError();
      }
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all paths except:
     *  - /login (public auth page)
     *  - /api/auth/* (auth Route Handlers)
     *  - /_next/* (Next.js internals)
     *  - /favicon.ico, /public assets
     */
    "/((?!login|api/auth|_next/static|_next/image|favicon\\.ico|.*\\.svg).*)",
  ],
};
