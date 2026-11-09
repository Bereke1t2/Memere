"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { clientAction } from "@/lib/client-action";

/**
 * Publishes or unpublishes an exam. Publishing is what makes a mock exam
 * reachable by students, so unpublishing exists to pull a live exam back for
 * correction rather than leaving a mistake in front of candidates.
 */
export function ExamPublishToggle({
  examId,
  isPublished,
}: {
  examId: string;
  isPublished: boolean;
}) {
  const [working, setWorking] = useState(false);
  const router = useRouter();

  async function toggle() {
    setWorking(true);
    try {
      await clientAction(
        `/api/teacher/exams/${examId}/${isPublished ? "unpublish" : "publish"}`,
        {}
      );
      toast.success(isPublished ? "Exam unpublished." : "Exam published.");
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to change publish state.");
    } finally {
      setWorking(false);
    }
  }

  return (
    <Button size="sm" variant={isPublished ? "outline" : "default"} onClick={toggle} disabled={working}>
      {working ? (
        <Loader2 className="mr-1.5 h-4 w-4 animate-spin" />
      ) : isPublished ? (
        <EyeOff className="mr-1.5 h-4 w-4" />
      ) : (
        <Eye className="mr-1.5 h-4 w-4" />
      )}
      {isPublished ? "Unpublish" : "Publish"}
    </Button>
  );
}
