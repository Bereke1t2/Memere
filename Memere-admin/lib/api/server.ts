import "server-only";

import { z } from "zod";
import { env } from "@/lib/env";
import { ApiError } from "./errors";
import { ErrorEnvelopeSchema } from "./schemas";
import { getAccessToken } from "@/lib/auth/cookies";
import { refreshAccessToken } from "@/lib/auth/refresh";

// ---- Core fetch ---------------------------------------------------------------

interface FetchOptions<T extends z.ZodTypeAny> {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  schema?: T;
  /** Override token (e.g. during login before cookies are set). */
  token?: string;
}

export async function apiFetch<T extends z.ZodTypeAny>(
  path: string,
  options: FetchOptions<T> = {}
): Promise<z.infer<T> | undefined> {
  const { method = "GET", body, schema, token } = options;

  const url = `${env.API_BASE_URL}/api/v1${path}`;

  async function doFetch(bearer: string | undefined) {
    const headers: HeadersInit = { "Content-Type": "application/json" };
    if (bearer) headers["Authorization"] = `Bearer ${bearer}`;
    return fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      cache: "no-store",
    });
  }

  const initialToken = token ?? (await getAccessToken());
  let res = await doFetch(initialToken);

  // Single 401 retry via refresh (only when token came from cookie, not caller)
  if (res.status === 401 && !token) {
    const fresh = await refreshAccessToken();
    if (fresh) {
      res = await doFetch(fresh);
    }
    if (res.status === 401) {
      throw new ApiError("SESSION_EXPIRED", "Your session has expired.", 401);
    }
  }

  // 204 No Content — nothing to return
  if (res.status === 204) {
    return undefined;
  }

  const json = await res.json().catch(() => null);

  if (!res.ok) {
    const parsed = ErrorEnvelopeSchema.safeParse(json);
    if (parsed.success) {
      throw new ApiError(
        parsed.data.code,
        parsed.data.message,
        res.status,
        parsed.data.details
      );
    }
    throw new ApiError("UNKNOWN_ERROR", "An unexpected error occurred.", res.status);
  }

  if (schema) {
    const parsed = schema.safeParse(json);
    if (!parsed.success) {
      console.error("[apiFetch] schema mismatch:", parsed.error.flatten());
      throw new ApiError("SCHEMA_ERROR", "Unexpected response shape.", 500);
    }
    return parsed.data as z.infer<T>;
  }

  return json;
}
