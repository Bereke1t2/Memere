import "server-only";

import { redirect } from "next/navigation";
import { me } from "@/lib/api/endpoints";
import type { User } from "@/lib/api/schemas";

export interface Session {
  user: User;
}

/**
 * Returns the current admin session by calling GET /auth/me with the
 * httpOnly access cookie. Returns null if unauthenticated or on any error.
 * This is defense-in-depth: a forged cookie value still fails the backend check.
 */
export async function getSession(): Promise<Session | null> {
  try {
    const user = await me();
    return { user };
  } catch {
    // ApiError (4xx/5xx) or network failure — both mean no valid session.
    return null;
  }
}

/**
 * Asserts the request is from an authenticated admin.
 * Redirects to /login otherwise — never throws.
 */
export async function requireAdmin(): Promise<Session> {
  const session = await getSession();
  if (!session) redirect("/login");
  if (session.user.role !== "admin") redirect("/login");
  return session;
}

/**
 * Allows both admin and teacher. Redirects to /login for anyone else.
 */
export async function requireStaff(): Promise<Session> {
  const session = await getSession();
  if (!session) redirect("/login");
  if (session.user.role !== "admin" && session.user.role !== "teacher")
    redirect("/login");
  return session;
}

export function isAdmin(session: Session): boolean {
  return session.user.role === "admin";
}
