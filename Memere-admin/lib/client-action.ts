"use client";

import { toast } from "sonner";

function handleSessionExpired() {
  toast.error("Your session has expired. Please log in again.");
  window.location.replace("/login");
}

/**
 * Makes a POST to a Next.js Route Handler from a client component.
 * Handles 401 (session expired) transparently: toasts and redirects to /login.
 * Returns undefined for 204 No Content, or the parsed JSON body.
 */
export async function clientAction<T = void>(
  url: string,
  body?: Record<string, unknown>,
  method: "POST" | "PUT" | "DELETE" | "GET" = "POST"
): Promise<T> {
  const res = await fetch(url, {
    method,
    headers: body !== undefined ? { "Content-Type": "application/json" } : {},
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (res.status === 401) {
    handleSessionExpired();
    return undefined as unknown as T;
  }

  if (res.status === 204) return undefined as unknown as T;

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    throw new Error((data as { message?: string }).message ?? "Action failed.");
  }

  return data as T;
}
