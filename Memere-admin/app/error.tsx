"use client";

import { useEffect } from "react";
import { ErrorState } from "@/components/ui/error-state";
import { Button } from "@/components/ui/button";

export default function RootError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("[RootError]", error);
  }, [error]);

  return (
    <main className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md text-center">
        <ErrorState error={error} reset={reset} />
        <div className="mt-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              window.location.href = "/login";
            }}
          >
            Return to Login
          </Button>
        </div>
      </div>
    </main>
  );
}
