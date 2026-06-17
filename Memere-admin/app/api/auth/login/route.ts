import { NextRequest, NextResponse } from "next/server";
import { login } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { setAuthCookies } from "@/lib/auth/cookies";

export async function POST(req: NextRequest) {
  try {
    const { email, password } = await req.json();

    const auth = await login(email, password);

    const role = auth.user?.role;
    if (role !== "admin" && role !== "teacher") {
      return NextResponse.json(
        {
          code: "NOT_STAFF",
          message: "Only admins and teachers can access this panel.",
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
