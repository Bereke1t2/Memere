import { NextRequest, NextResponse } from "next/server";
import { logout } from "@/lib/api/endpoints";
import { clearAuthCookies, getRefreshToken } from "@/lib/auth/cookies";

export async function POST() {
  try {
    const refreshToken = await getRefreshToken();
    if (refreshToken) {
      await logout(refreshToken).catch(() => {
        // Backend logout failure is non-fatal — we still clear cookies locally.
      });
    }
  } finally {
    await clearAuthCookies();
  }

  return NextResponse.json({ ok: true });
}
