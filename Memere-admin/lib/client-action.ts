"use client";

import { toast } from "sonner";

function handleSessionExpired() {
  toast.error("Your session has expired. Please log in again.");
  setTimeout(() => {
    window.location.href = "/login";
  }, 500);
}

/**
 * Makes a request to a Next.js Route Handler from a client component.
 * Handles 401 (session expired) transparently: toasts, redirects to /login, and throws.
 * Returns undefined for 204 No Content, or the parsed JSON body.
 */
export async function clientAction<T = void>(
  url: string,
  body?: Record<string, unknown>,
  method: "POST" | "PUT" | "DELETE" | "GET" = "POST"
): Promise<T> {
  let res: Response;
  try {
    res = await fetch(url, {
      method,
      headers: body !== undefined ? { "Content-Type": "application/json" } : {},
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch (err) {
    const message =
      err instanceof Error ? err.message : "Network error. Please check your connection.";
    throw new Error(message);
  }

  if (res.status === 401) {
    handleSessionExpired();
    throw new Error("Your session has expired. Please log in again.");
  }

  if (res.status === 204) {
    return undefined as unknown as T;
  }

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    let errorMsg =
      (data as { message?: string; error?: string })?.message ||
      (data as { message?: string; error?: string })?.error ||
      `Request failed with status ${res.status}`;

    const details = (data as { details?: unknown })?.details;
    if (details && typeof details === "object") {
      const fieldErrors = Object.entries(details as Record<string, unknown>)
        .map(([field, msgs]) => `${field}: ${Array.isArray(msgs) ? msgs.join(", ") : String(msgs)}`)
        .filter(Boolean)
        .join("; ");
      if (fieldErrors) {
        errorMsg += ` (${fieldErrors})`;
      }
    }
    throw new Error(errorMsg);
  }

  return data as T;
}
