import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function CourseNotFound() {
  return (
    <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
      <p className="text-lg font-medium">Course not found</p>
      <p className="text-sm text-muted-foreground">
        This course may have been deleted or the ID is invalid.
      </p>
      <Button asChild variant="outline" size="sm">
        <Link href="/courses">Back to courses</Link>
      </Button>
    </div>
  );
}
