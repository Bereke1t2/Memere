import { Button } from "@/components/ui/button";
import { ApiError, friendlyMessage } from "@/lib/api/errors";

interface ErrorStateProps {
  error: Error & { digest?: string };
  reset?: () => void;
}

export function ErrorState({ error, reset }: ErrorStateProps) {
  const message =
    error instanceof ApiError
      ? friendlyMessage(error)
      : "Something went wrong. Please try again.";

  return (
    <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
      <p className="text-sm text-muted-foreground max-w-sm">{message}</p>
      {reset && (
        <Button onClick={reset} variant="outline" size="sm">
          Try again
        </Button>
      )}
    </div>
  );
}
