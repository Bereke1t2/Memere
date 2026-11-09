import { getRouteStaffSession } from "@/lib/auth/session";
import { canManageContent } from "@/lib/auth/roles";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

/**
 * Wraps a route handler with the staff session + content-permission check every
 * authoring route needs, and maps backend ApiErrors onto the response.
 *
 * These routes proxy the backend's authoring endpoints, several of which return
 * answer keys. Centralizing the guard means a new authoring route cannot ship
 * with the check accidentally omitted — and the backend re-checks ownership
 * regardless, so this is defence in depth rather than the only gate.
 */
export async function withStaff<T>(
  fallbackMessage: string,
  handler: () => Promise<T>
): Promise<Response> {
  try {
    const session = await getRouteStaffSession();
    if (!session) {
      return Response.json({ message: "Unauthorized", code: "UNAUTHORIZED" }, { status: 401 });
    }
    if (!canManageContent(session.user)) {
      return Response.json({ message: "Forbidden" }, { status: 403 });
    }
    const result = await handler();
    if (result === undefined) {
      return new Response(null, { status: 204 });
    }
    return Response.json(result);
  } catch (err) {
    if (err instanceof ApiError) {
      return Response.json({ message: friendlyMessage(err), details: err.details }, { status: err.status });
    }
    const message = err instanceof Error ? err.message : fallbackMessage;
    return Response.json({ message }, { status: 500 });
  }
}

/** The 400 response shape the client's error formatter already understands. */
export function validationResponse(details: unknown): Response {
  return Response.json({ message: "Validation failed", details }, { status: 400 });
}
