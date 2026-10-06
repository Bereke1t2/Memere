import { NextRequest, NextResponse } from "next/server";
import { sendAccountDeletionRequest } from "@/lib/resend";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const { email, reason } = body;

    if (!email) {
      return NextResponse.json(
        {
          success: false,
          error: "Email address is required.",
        },
        { status: 400 }
      );
    }

    const trimmedEmail = email.trim().toLowerCase();

    // Trigger account deletion request email notification
    await sendAccountDeletionRequest({
      email: trimmedEmail,
      reason: reason?.trim(),
    });

    return NextResponse.json({
      success: true,
      message: "Account deletion request received and queued for processing.",
    });
  } catch (error: any) {
    console.error("Account deletion route error:", error);
    return NextResponse.json(
      {
        success: false,
        error: error.message || "Failed to submit account deletion request.",
      },
      { status: 500 }
    );
  }
}
