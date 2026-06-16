import { NextRequest, NextResponse } from "next/server";

const ACCESS_COOKIE = "mm_access";
const REFRESH_COOKIE = "mm_refresh";

export function middleware(req: NextRequest) {
  const access = req.cookies.get(ACCESS_COOKIE)?.value;
  const refresh = req.cookies.get(REFRESH_COOKIE)?.value;

  if (!access && !refresh) {
    const loginUrl = new URL("/login", req.url);
    loginUrl.searchParams.set("from", req.nextUrl.pathname);
    return NextResponse.redirect(loginUrl);
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
