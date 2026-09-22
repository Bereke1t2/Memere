import "server-only";

import { redirect } from "next/navigation";
import { me } from "@/lib/api/endpoints";
import type { User } from "@/lib/api/schemas";

export interface Session {
  user: User;
}

/**
 * Returns the current staff session by calling GET /auth/me with the
 * httpOnly access cookie. Returns null if unauthenticated or on any error.
 * This is defense-in-depth: a forged cookie value still fails the backend check.
 */
export async function getSession(): Promise<Session | null> {
  try {
    const user = await me();
    if (!user || !user.role) return null;
    return { user };
  } catch {
    // ApiError (4xx/5xx), network failure, or missing env — return null safely.
    return null;
  }
}

/**
 * Asserts the request is from an authenticated admin in Server Components.
 * Redirects to /login otherwise — throws Next.js redirect.
 */
export async function requireAdmin(): Promise<Session> {
  const session = await getSession();
  if (!session || session.user.role !== "admin") {
    redirect("/login");
  }
  return session;
}

/**
 * Allows both admin and teacher in Server Components.
 * Redirects to /login for anyone else — throws Next.js redirect.
 */
export async function requireStaff(): Promise<Session> {
  const session = await getSession();
  if (
    !session ||
    (session.user.role !== "admin" && session.user.role !== "teacher")
  ) {
    redirect("/login");
  }
  return session;
}

/**
 * Safe session checker specifically for Next.js Route Handlers (app/api/*).
 * Never throws NEXT_REDIRECT so route handlers can return clean JSON responses (401/403).
 */
export async function getRouteStaffSession(): Promise<{ user: User } | null> {
  const session = await getSession();
  if (
    !session ||
    (session.user.role !== "admin" && session.user.role !== "teacher")
  ) {
    return null;
  }
  return session;
}

/**
 * Safe admin session checker specifically for Next.js Route Handlers (app/api/*).
 * Never throws NEXT_REDIRECT so route handlers can return clean JSON responses (401/403).
 */
export async function getRouteAdminSession(): Promise<{ user: User } | null> {
  const session = await getSession();
  if (!session || session.user.role !== "admin") {
    return null;
  }
  return session;
}

export function isAdmin(session: Session): boolean {
  return session.user.role === "admin";
}
