import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function UserNotFound() {
  return (
    <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
      <p className="text-lg font-medium">User not found</p>
      <p className="text-sm text-muted-foreground">
        This user may have been deleted or the ID is invalid.
      </p>
      <Button asChild variant="outline" size="sm">
        <Link href="/users">Back to users</Link>
      </Button>
    </div>
  );
}
