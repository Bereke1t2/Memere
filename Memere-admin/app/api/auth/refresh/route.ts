import { NextResponse } from "next/server";
import { refresh } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { clearAuthCookies, getRefreshToken, setAuthCookies } from "@/lib/auth/cookies";

export async function POST() {
  const refreshToken = await getRefreshToken();
  if (!refreshToken) {
    return NextResponse.json(
      { code: "NO_REFRESH_TOKEN", message: "No refresh token present." },
      { status: 401 }
    );
  }

  try {
    const auth = await refresh(refreshToken);

    await setAuthCookies({
      access: auth.access_token,
      refresh: auth.refresh_token,
      expiresIn: auth.expires_in,
    });

    return NextResponse.json({ ok: true });
  } catch (err) {
    await clearAuthCookies();
    if (err instanceof ApiError) {
      return NextResponse.json(
        { code: err.code, message: err.message },
        { status: 401 }
      );
    }
    return NextResponse.json(
      { code: "REFRESH_FAILED", message: "Session could not be renewed." },
      { status: 401 }
    );
  }
}
