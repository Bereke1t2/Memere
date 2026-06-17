"use client";

import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { useSearchParams } from "next/navigation";

export interface CursorPage<T> {
  items: T[];
  next: string | null;
}

/**
 * Fetches a cursor-paginated list from a Next.js Route Handler.
 *
 * The hook reads ALL current URL search params (q, role, limit, after, …) and
 * passes them as query-string to `route`. The Route Handler runs server-side,
 * reads the httpOnly access cookie, and calls the Go backend — so tokens never
 * reach the browser.
 *
 * When any URL param changes (filter or cursor), TanStack Query automatically
 * refetches because params are part of the query key.
 */
export function useCursorList<T>({
  route,
  queryKey,
  initialData,
}: {
  route: string;
  queryKey: string[];
  initialData?: CursorPage<T>;
}): UseQueryResult<CursorPage<T>> {
  const searchParams = useSearchParams();
  const params = Object.fromEntries(searchParams.entries());

  return useQuery<CursorPage<T>>({
    queryKey: [...queryKey, params],
    queryFn: async () => {
      const qs = new URLSearchParams(params).toString();
      const res = await fetch(`${route}?${qs}`);

      if (res.status === 401) {
        const { toast } = await import("sonner");
        toast.error("Your session has expired. Please log in again.");
        window.location.replace("/login");
        return { items: [], next: null } as CursorPage<T>;
      }

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(
          (body as { message?: string }).message ?? "Failed to fetch"
        );
      }
      return res.json() as Promise<CursorPage<T>>;
    },
    initialData,
    staleTime: 30_000,
  });
}
