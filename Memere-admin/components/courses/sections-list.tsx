"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Plus, Loader2, BookOpen, ChevronDown, ChevronRight } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from "@/components/ui/dialog";
import { clientAction } from "@/lib/client-action";
import { LessonsList } from "./lessons-list";
import type { Section, Lesson } from "@/lib/api/schemas";

const AddSectionSchema = z.object({
  title: z.string().min(1, "Title is required"),
  order: z.number().optional(),
});
type AddSectionInput = z.infer<typeof AddSectionSchema>;

interface SectionsListProps {
  courseId: string;
  sections: Section[];
  lessonsBySectionId: Record<string, Lesson[]>;
  canEdit?: boolean;
}

function SectionRow({
  section,
  lessons,
  index,
  canEdit,
}: {
  section: Section;
  lessons: Lesson[];
  index: number;
  canEdit: boolean;
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="border rounded-lg overflow-hidden">
      <button
        type="button"
        className="w-full flex items-center gap-3 px-4 py-3 hover:bg-muted/30 transition-colors text-left"
        onClick={() => setExpanded((e) => !e)}
      >
        {expanded ? (
          <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
        )}
        <span className="w-6 h-6 rounded-full bg-muted flex items-center justify-center text-xs font-medium text-muted-foreground shrink-0">
          {section.order ?? section.order_index ?? index + 1}
        </span>
        <span className="flex-1 text-sm font-medium">{section.title}</span>
        <span className="text-xs text-muted-foreground tabular-nums">
          {lessons.length} lesson{lessons.length !== 1 ? "s" : ""}
        </span>
      </button>
      {expanded && (
        <div className="border-t px-4 py-3 bg-muted/10">
          <LessonsList sectionId={section.id} lessons={lessons} canEdit={canEdit} />
        </div>
      )}
    </div>
  );
}

export function SectionsList({ courseId, sections, lessonsBySectionId, canEdit = true }: SectionsListProps) {
  const [open, setOpen] = useState(false);
  const router = useRouter();

  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } =
    useForm<AddSectionInput>({ resolver: zodResolver(AddSectionSchema) });

  async function onSubmit(values: AddSectionInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/sections`, {
        ...values,
        is_published: true,
      });
      toast.success("Section added.");
      setOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add section.");
    }
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <p className="text-sm font-semibold">Sections & Lessons</p>
        {canEdit && (
          <Button size="sm" variant="outline" onClick={() => setOpen(true)}>
            <Plus className="h-4 w-4 mr-1.5" />
            Add Section
          </Button>
        )}
      </div>

      {sections.length === 0 ? (
        <Card className="border shadow-none">
          <CardContent className="flex flex-col items-center gap-2 py-10 text-center">
            <BookOpen className="h-8 w-8 text-muted-foreground/40" />
            <p className="text-sm text-muted-foreground">No sections yet. Add your first section to get started.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="flex flex-col gap-2">
          {sections.map((section, idx) => (
            <SectionRow
              key={section.id}
              section={section}
              lessons={lessonsBySectionId[section.id] ?? []}
              index={idx}
              canEdit={canEdit}
            />
          ))}
        </div>
      )}

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-sm" aria-describedby={undefined}>
          <DialogHeader><DialogTitle>Add Section</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="section-title">Title</Label>
              <Input id="section-title" placeholder="e.g. Introduction to Algebra" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="section-order">Order (optional)</Label>
              <Input id="section-order" type="number" min="1" placeholder="Auto" {...register("order", { valueAsNumber: true })} />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); reset(); }}>Cancel</Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Adding…</> : "Add Section"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
