import "server-only";

import { cookies } from "next/headers";

export const ACCESS_COOKIE = "mm_access";
export const REFRESH_COOKIE = "mm_refresh";

const IS_PROD = process.env.NODE_ENV === "production";
const REFRESH_MAX_AGE = 60 * 60 * 24 * 30; // 30 days

const COOKIE_BASE = {
  httpOnly: true,
  secure: IS_PROD,
  sameSite: "lax" as const,
  path: "/",
};

export async function getAccessToken(): Promise<string | undefined> {
  const jar = await cookies();
  return jar.get(ACCESS_COOKIE)?.value;
}

export async function getRefreshToken(): Promise<string | undefined> {
  const jar = await cookies();
  return jar.get(REFRESH_COOKIE)?.value;
}

export async function setAuthCookies(tokens: {
  access: string;
  refresh: string;
  expiresIn: number;
}): Promise<void> {
  const jar = await cookies();
  jar.set(ACCESS_COOKIE, tokens.access, {
    ...COOKIE_BASE,
    maxAge: tokens.expiresIn,
  });
  jar.set(REFRESH_COOKIE, tokens.refresh, {
    ...COOKIE_BASE,
    maxAge: REFRESH_MAX_AGE,
  });
}

export async function clearAuthCookies(): Promise<void> {
  const jar = await cookies();
  jar.set(ACCESS_COOKIE, "", { ...COOKIE_BASE, maxAge: 0 });
  jar.set(REFRESH_COOKIE, "", { ...COOKIE_BASE, maxAge: 0 });
}
