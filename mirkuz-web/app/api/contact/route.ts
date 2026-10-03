import { NextRequest, NextResponse } from "next/server";
import { sendContactInquiry } from "@/lib/resend";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const { name, email, phone, stream, message } = body;

    if (!name || !email || !message) {
      return NextResponse.json(
        {
          success: false,
          error: "Name, email, and message are required.",
        },
        { status: 400 }
      );
    }

    const result = await sendContactInquiry({
      name: name.trim(),
      email: email.trim().toLowerCase(),
      phone: phone?.trim(),
      stream: stream?.trim(),
      message: message.trim(),
    });

    return NextResponse.json({
      success: true,
      message: "Your message has been sent successfully!",
      data: result,
    });
  } catch (error: any) {
    console.error("Contact inquiry error:", error);
    return NextResponse.json(
      {
        success: false,
        error: error.message || "Failed to send contact inquiry.",
      },
      { status: 500 }
    );
  }
}
