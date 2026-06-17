"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { GraduationCap, AlertCircle, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";

const loginSchema = z.object({
  email: z.string().email("Enter a valid email address."),
  password: z.string().min(1, "Password is required."),
});

type LoginValues = z.infer<typeof loginSchema>;

const CODE_COPY: Record<string, string> = {
  INVALID_CREDENTIALS: "Email or password is incorrect.",
  NOT_ADMIN: "This account is not an administrator.",
  FORBIDDEN: "This account is not an administrator.",
};

export function LoginForm() {
  const router = useRouter();
  const [formError, setFormError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginValues>({ resolver: zodResolver(loginSchema) });

  async function onSubmit(values: LoginValues) {
    setFormError(null);
    try {
      const res = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(values),
      });

      if (res.ok) {
        router.push("/");
        router.refresh();
        return;
      }

      const body = await res.json().catch(() => ({}));
      const msg =
        CODE_COPY[body?.code as string] ??
        body?.message ??
        "Login failed. Please try again.";
      setFormError(msg);
    } catch {
      setFormError("Network error. Please check your connection.");
    }
  }

  return (
    <Card className="w-full max-w-md border-0 shadow-xl">
      <CardContent className="pt-8 pb-8 px-8">
        <div className="flex flex-col items-center gap-3 pb-6">
          <div className="rounded-2xl bg-gradient-to-br from-indigo-500 to-violet-600 p-3">
            <GraduationCap className="h-8 w-8 text-white" />
          </div>
          <div className="text-center">
            <h1 className="text-2xl font-bold">Memere Admin</h1>
            <p className="mt-1 text-sm text-muted-foreground">Sign in to your admin account</p>
          </div>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4">
          <div>
            <Label htmlFor="email" className="text-sm font-medium text-foreground mb-1.5 block">
              Email address
            </Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="admin@example.com"
              className="h-11"
              aria-invalid={!!errors.email}
              {...register("email")}
            />
            {errors.email && (
              <p className="mt-1.5 text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div>
            <Label htmlFor="password" className="text-sm font-medium text-foreground mb-1.5 block">
              Password
            </Label>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              className="h-11"
              aria-invalid={!!errors.password}
              {...register("password")}
            />
            {errors.password && (
              <p className="mt-1.5 text-sm text-destructive">{errors.password.message}</p>
            )}
          </div>

          {formError && (
            <div
              role="alert"
              className="flex items-center gap-2 rounded-lg bg-destructive/10 px-3 py-2.5 text-sm text-destructive"
            >
              <AlertCircle className="h-4 w-4 shrink-0" />
              {formError}
            </div>
          )}

          <Button
            type="submit"
            className="w-full h-11 bg-indigo-600 hover:bg-indigo-700 text-white font-medium"
            disabled={isSubmitting}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Loading...
              </>
            ) : (
              "Sign in"
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
