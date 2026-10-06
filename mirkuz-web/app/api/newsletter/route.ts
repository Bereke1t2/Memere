import { NextRequest, NextResponse } from "next/server";
import { sendWelcomeNewsletterEmail } from "@/lib/resend";

export const dynamic = "force-dynamic";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const { email } = body;

    if (!email || typeof email !== "string" || !email.includes("@")) {
      return NextResponse.json(
        { success: false, error: "Please provide a valid email address." },
        { status: 400 }
      );
    }

    const result = await sendWelcomeNewsletterEmail(email.trim().toLowerCase());

    return NextResponse.json({
      success: true,
      message: "Successfully subscribed! A welcome email has been sent.",
      data: result,
    });
  } catch (error: any) {
    console.error("Newsletter subscription error:", error);
    return NextResponse.json(
      {
        success: false,
        error: error.message || "Failed to process newsletter subscription.",
      },
      { status: 500 }
    );
  }
}
