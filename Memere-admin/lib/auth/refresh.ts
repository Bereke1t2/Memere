import "server-only";

import { env } from "@/lib/env";
import {
  ACCESS_COOKIE,
  REFRESH_COOKIE,
  getRefreshToken,
} from "@/lib/auth/cookies";
import { cookies } from "next/headers";

const IS_PROD = process.env.NODE_ENV === "production";

/**
 * Exchanges the mm_refresh cookie for a new access token by calling the
 * backend directly. Updates mm_access cookie when called from a Route Handler
 * or Server Action context (cookie.set is a no-op in plain Server Components).
 * Returns the new access token, or undefined on failure.
 */
export async function refreshAccessToken(): Promise<string | undefined> {
  const refreshToken = await getRefreshToken();
  if (!refreshToken) return undefined;

  try {
    const res = await fetch(`${env.API_BASE_URL}/api/v1/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
      cache: "no-store",
    });

    if (!res.ok) return undefined;

    const json = await res.json();
    const newAccess: string | undefined = json.access_token;
    const newRefresh: string | undefined = json.refresh_token;
    const expiresIn: number = json.expires_in ?? 900;

    if (!newAccess) return undefined;

    const jar = await cookies();
    try {
      jar.set(ACCESS_COOKIE, newAccess, {
        httpOnly: true,
        secure: IS_PROD,
        sameSite: "lax",
        path: "/",
        maxAge: expiresIn,
      });
      if (newRefresh) {
        jar.set(REFRESH_COOKIE, newRefresh, {
          httpOnly: true,
          secure: IS_PROD,
          sameSite: "lax",
          path: "/",
          maxAge: 60 * 60 * 24 * 30,
        });
      }
    } catch {
      // cookies().set() throws in Server Component rendering — the new token
      // is still returned so the current request can retry successfully.
    }

    return newAccess;
  } catch {
    return undefined;
  }
}
