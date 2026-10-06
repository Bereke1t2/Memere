import { Resend } from "resend";

// Initialize the Resend client with the production API key
export const resend = new Resend(
  process.env.RESEND_API_KEY || ""
);

// Default sender configurations
export const SENDER_EMAIL =
  process.env.RESEND_FROM_EMAIL || "Mirkuz <noreply@mirkuz.app>";
export const SENDER_FALLBACK =
  process.env.RESEND_ONBOARDING_FALLBACK || "Mirkuz <onboarding@resend.dev>";
export const ADMIN_EMAIL =
  process.env.ADMIN_NOTIFICATION_EMAIL || "info@mirkuz.app";

/**
 * Send Welcome Newsletter Email to a new subscriber
 */
export async function sendWelcomeNewsletterEmail(subscriberEmail: string) {
  const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Welcome to Mirkuz</title>
</head>
<body style="margin: 0; padding: 0; background-color: #F5F7FA; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; color: #4D4D4D;">
  <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #F5F7FA; padding: 40px 16px;">
    <tr>
      <td align="center">
        <!-- Container Card -->
        <table width="100%" max-width="600" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 16px rgba(171, 190, 209, 0.3); border: 1px solid #E8ECF2;">
          
          <!-- Header Bar -->
          <tr>
            <td style="background-color: #4CAF4F; padding: 28px 36px; text-align: left;">
              <table border="0" cellspacing="0" cellpadding="0">
                <tr>
                  <td style="background-color: #ffffff; border-radius: 8px; width: 36px; height: 36px; text-align: center; vertical-align: middle;">
                    <span style="color: #4CAF4F; font-size: 20px; font-weight: bold; line-height: 36px;">M</span>
                  </td>
                  <td style="padding-left: 12px;">
                    <span style="color: #ffffff; font-size: 22px; font-weight: bold; letter-spacing: -0.5px;">Mirkuz</span>
                    <span style="color: #E8F5E9; font-size: 13px; font-weight: 500; margin-left: 6px;">(ምርኩዝ)</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Body Content -->
          <tr>
            <td style="padding: 36px 36px 28px 36px;">
              <h1 style="color: #263238; font-size: 24px; font-weight: bold; margin: 0 0 16px 0; line-height: 1.3;">
                Welcome to Mirkuz Exam Prep! 🇪🇹
              </h1>
              
              <p style="color: #717171; font-size: 15px; line-height: 1.6; margin: 0 0 20px 0;">
                Thank you for subscribing to the <strong>Mirkuz</strong> community newsletter. You’ll be the first to receive national exam strategies, high-yield topic breakdowns, solved matric exam questions, and platform updates.
              </p>

              <!-- Highlight Box -->
              <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #F5F7FA; border-radius: 8px; border-left: 4px solid #4CAF4F; margin-bottom: 24px;">
                <tr>
                  <td style="padding: 16px 20px;">
                    <p style="color: #4D4D4D; font-size: 14px; font-weight: bold; margin: 0 0 6px 0;">What you get with Mirkuz:</p>
                    <ul style="color: #717171; font-size: 13px; margin: 0; padding-left: 20px; line-height: 1.6;">
                      <li>50+ HD video courses across Natural & Social streams</li>
                      <li>10,000+ past national matric questions (2008–2016 E.C.)</li>
                      <li>100% data-free encrypted offline study mode</li>
                      <li>Timed mock exam simulator with instant solutions</li>
                    </ul>
                  </td>
                </tr>
              </table>

              <!-- CTA Button -->
              <table border="0" cellspacing="0" cellpadding="0" style="margin: 28px 0 20px 0;">
                <tr>
                  <td align="center" style="border-radius: 6px; background-color: #4CAF4F;">
                    <a href="https://play.google.com/store/apps" target="_blank" style="font-size: 14px; font-weight: bold; color: #ffffff; text-decoration: none; padding: 14px 28px; display: inline-block; border-radius: 6px;">
                      Download Mirkuz App on Google Play →
                    </a>
                  </td>
                </tr>
              </table>

              <p style="color: #89939E; font-size: 13px; line-height: 1.5; margin: 24px 0 0 0;">
                Have questions or need study guidance? Join our active Telegram student community at <a href="https://t.me/mirkuz_exam" style="color: #4CAF4F; text-decoration: none; font-weight: bold;">t.me/mirkuz_exam</a>.
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background-color: #263238; padding: 24px 36px; text-align: center;">
              <p style="color: #ABBED1; font-size: 12px; margin: 0 0 6px 0;">
                © 2026 Mirkuz (ምርኩዝ) Education Platform. Addis Ababa, Ethiopia.
              </p>
              <p style="color: #717171; font-size: 11px; margin: 0;">
                You received this email because you subscribed on <a href="https://mirkuz.app" style="color: #4CAF4F; text-decoration: none;">mirkuz.app</a>.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>
  `;

  try {
    const result = await resend.emails.send({
      from: SENDER_EMAIL,
      to: subscriberEmail,
      subject: "Welcome to Mirkuz — Grade 12 National Exam Prep",
      html: htmlContent,
    });

    if (result.error) {
      if (result.error.message?.includes("domain") || (result.error as any)?.statusCode === 403) {
        const fallbackResult = await resend.emails.send({
          from: SENDER_FALLBACK,
          to: subscriberEmail,
          subject: "Welcome to Mirkuz — Grade 12 National Exam Prep",
          html: htmlContent,
        });
        if (fallbackResult.error) {
          throw new Error(fallbackResult.error.message);
        }
        return { success: true, data: fallbackResult.data };
      }
      throw new Error(result.error.message);
    }

    return { success: true, data: result.data };
  } catch (error: any) {
    // If domain mirkuz.app is still undergoing DNS verification, retry with onboarding fallback
    if (error?.message?.includes("domain") || error?.statusCode === 403) {
      const fallbackResult = await resend.emails.send({
        from: SENDER_FALLBACK,
        to: subscriberEmail,
        subject: "Welcome to Mirkuz — Grade 12 National Exam Prep",
        html: htmlContent,
      });
      if (fallbackResult.error) {
        throw new Error(fallbackResult.error.message);
      }
      return { success: true, data: fallbackResult.data };
    }
    throw error;
  }
}

/**
 * Send Contact / Demo Request Notification to Admin & Confirmation to User
 */
export async function sendContactInquiry({
  name,
  email,
  phone,
  stream,
  message,
}: {
  name: string;
  email: string;
  phone?: string;
  stream?: string;
  message: string;
}) {
  // 1. Notification to Admin
  const adminHtml = `
    <h2>New Demo / Student Inquiry on mirkuz.app</h2>
    <p><strong>Name:</strong> ${name}</p>
    <p><strong>Email:</strong> ${email}</p>
    <p><strong>Phone:</strong> ${phone || "N/A"}</p>
    <p><strong>Stream:</strong> ${stream || "N/A"}</p>
    <p><strong>Message:</strong></p>
    <blockquote style="background: #f5f7fa; padding: 12px; border-left: 4px solid #4CAF4F;">
      ${message}
    </blockquote>
  `;

  // 2. Confirmation to User
  const userHtml = `
    <div style="font-family: sans-serif; color: #4D4D4D; max-width: 550px;">
      <h2 style="color: #4CAF4F;">Thank you for contacting Mirkuz, ${name}!</h2>
      <p>We have received your message regarding Grade 12 exam preparation. Our admissions and academic advisors will get back to you within 24 hours.</p>
      <p>In the meantime, you can explore courses or download the app at <a href="https://mirkuz.app">mirkuz.app</a>.</p>
      <hr style="border: 0; border-top: 1px solid #E8ECF2; margin: 20px 0;" />
      <p style="font-size: 12px; color: #89939E;">Mirkuz Education Platform • Addis Ababa, Ethiopia</p>
    </div>
  `;

  // Send admin notification
  const adminRes = await resend.emails.send({
    from: SENDER_EMAIL,
    to: ADMIN_EMAIL,
    subject: `New Student Inquiry: ${name} (${stream || "Grade 12"})`,
    html: adminHtml,
  });

  // Send student confirmation
  const userRes = await resend.emails.send({
    from: SENDER_EMAIL,
    to: email,
    subject: "We received your inquiry — Mirkuz Exam Prep",
    html: userHtml,
  });

  if (adminRes.error && userRes.error) {
    throw new Error(adminRes.error.message || userRes.error.message);
  }

  return {
    admin: adminRes.data,
    user: userRes.data,
  };
}

/**
 * Send Account Deletion Request Notification to Admin & Confirmation to User
 */
export async function sendAccountDeletionRequest({
  email,
  reason,
}: {
  email: string;
  reason?: string;
}) {
  const adminHtml = `
    <h2>⚠️ Account Deletion Request Received (Google Play Data Safety)</h2>
    <p>A student has requested permanent deletion of their account and all personal data via mirkuz.app.</p>
    <p><strong>Account Email:</strong> ${email}</p>
    <p><strong>Reason Provided:</strong> ${reason || "None specified"}</p>
    <p><strong>Timestamp:</strong> ${new Date().toISOString()}</p>
    <p>Please process this deletion in the admin dashboard within 48 hours to comply with Google Play User Data policies.</p>
  `;

  const userHtml = `
    <div style="font-family: sans-serif; color: #4D4D4D; max-width: 550px;">
      <h2 style="color: #ef4444;">Account Deletion Request Received</h2>
      <p>Hello,</p>
      <p>We have received your formal request to delete your Mirkuz account linked to <strong>${email}</strong>.</p>
      <p>Your request is currently being processed. All account data, study history, and enrollment records will be permanently removed within 48 hours.</p>
      <p>If you did not request this deletion, please reply to this email immediately at <a href="mailto:support@mirkuz.app">support@mirkuz.app</a>.</p>
      <hr style="border: 0; border-top: 1px solid #E8ECF2; margin: 20px 0;" />
      <p style="font-size: 12px; color: #89939E;">Mirkuz (ምርኩዝ) Education Platform • Addis Ababa, Ethiopia</p>
    </div>
  `;

  try {
    const adminRes = await resend.emails.send({
      from: SENDER_EMAIL,
      to: ADMIN_EMAIL,
      subject: `[ACTION REQUIRED] Account Deletion Request: ${email}`,
      html: adminHtml,
    });

    const userRes = await resend.emails.send({
      from: SENDER_EMAIL,
      to: email,
      subject: "Account Deletion Request Confirmation — Mirkuz",
      html: userHtml,
    });

    return { admin: adminRes.data, user: userRes.data };
  } catch (error: any) {
    console.error("Account deletion email dispatch error:", error);
    // Don't throw if email service fails, to ensure request is recorded
    return { success: true };
  }
}
