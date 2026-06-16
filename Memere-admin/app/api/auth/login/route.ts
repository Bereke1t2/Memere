import { NextRequest, NextResponse } from "next/server";
import { login } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { setAuthCookies } from "@/lib/auth/cookies";

export async function POST(req: NextRequest) {
  try {
    const { email, password } = await req.json();

    const auth = await login(email, password);

    if (auth.user?.role !== "admin") {
      return NextResponse.json(
        {
          code: "NOT_ADMIN",
          message: "This account is not an administrator.",
        },
        { status: 403 }
      );
    }

    await setAuthCookies({
      access: auth.access_token,
      refresh: auth.refresh_token,
      expiresIn: auth.expires_in,
    });

    return NextResponse.json({ ok: true });
  } catch (err) {
    if (err instanceof ApiError) {
      const status = err.status === 401 ? 401 : err.status;
      return NextResponse.json(
        { code: err.code, message: err.message },
        { status }
      );
    }
    return NextResponse.json(
      { code: "INTERNAL_ERROR", message: "An unexpected error occurred." },
      { status: 500 }
    );
  }
}
