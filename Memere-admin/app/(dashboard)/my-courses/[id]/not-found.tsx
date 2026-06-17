import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function CourseNotFound() {
  return (
    <div className="flex flex-col items-center justify-center gap-4 py-24 text-center">
      <p className="text-sm font-medium">Course not found</p>
      <p className="text-sm text-muted-foreground">This course may have been deleted or doesn&apos;t belong to your account.</p>
      <Button asChild variant="outline" size="sm">
        <Link href="/my-courses">Back to My Courses</Link>
      </Button>
    </div>
  );
}
