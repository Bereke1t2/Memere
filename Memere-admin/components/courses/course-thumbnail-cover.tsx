"use client";

import { useState } from "react";
import { Camera, Image as ImageIcon } from "lucide-react";
import { TeacherThumbnailDialog } from "./teacher-thumbnail-dialog";
import type { Course } from "@/lib/api/schemas";

interface CourseThumbnailCoverProps {
  course: Course;
  canEdit?: boolean;
}

export function CourseThumbnailCover({
  course,
  canEdit = false,
}: CourseThumbnailCoverProps) {
  const [open, setOpen] = useState(false);

  if (!course.thumbnail_url && !canEdit) return null;

  return (
    <>
      <button
        type="button"
        disabled={!canEdit}
        onClick={() => canEdit && setOpen(true)}
        className={`relative h-20 w-32 shrink-0 rounded-lg overflow-hidden border bg-muted group text-left ${
          canEdit
            ? "cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            : "cursor-default"
        }`}
        title={canEdit ? "Click to change course thumbnail" : undefined}
      >
        {course.thumbnail_url ? (
          /* eslint-disable-next-line @next/next/no-img-element */
          <img
            src={course.thumbnail_url}
            alt={course.title}
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="h-full w-full flex flex-col items-center justify-center text-muted-foreground gap-1 p-1 bg-muted/30">
            <ImageIcon className="h-5 w-5 opacity-40" />
            <span className="text-[10px] text-center font-medium">Add Thumbnail</span>
          </div>
        )}

        {canEdit && (
          <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-0.5 text-white">
            <Camera className="h-4 w-4" />
            <span className="text-[10px] font-medium tracking-tight">Change</span>
          </div>
        )}
      </button>

      {canEdit && (
        <TeacherThumbnailDialog
          course={course}
          open={open}
          onOpenChange={setOpen}
        />
      )}
    </>
  );
}
