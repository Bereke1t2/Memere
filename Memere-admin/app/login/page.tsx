import { redirect } from "next/navigation";
import { getSession } from "@/lib/auth/session";
import { LoginForm } from "@/components/auth/login-form";
import type { Metadata } from "next";

// Always dynamic — reads httpOnly cookies to check existing session.
export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Sign In — Memere Admin" };

export default async function LoginPage() {
  const session = await getSession();
  if (session) redirect("/");

  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-muted/40 p-4">
      <LoginForm />
    </main>
  );
}
