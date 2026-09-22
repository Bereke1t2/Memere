"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Plus,
  Loader2,
  BookOpen,
  ChevronDown,
  ChevronRight,
  Pencil,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import { clientAction } from "@/lib/client-action";
import { LessonsList } from "./lessons-list";
import type { Section, Lesson, Quiz } from "@/lib/api/schemas";

const SectionFormSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  description: z.string().max(2000).optional(),
  order_index: z.number().int().min(0).optional(),
  is_published: z.boolean().optional(),
});
type SectionFormInput = z.infer<typeof SectionFormSchema>;

interface SectionsListProps {
  courseId: string;
  sections: Section[];
  lessonsBySectionId: Record<string, Lesson[]>;
  quizzes?: Quiz[];
  canEdit?: boolean;
}

function EditSectionDialog({
  section,
  open,
  onOpenChange,
}: {
  section: Section;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<SectionFormInput>({
    resolver: zodResolver(SectionFormSchema),
    defaultValues: {
      title: section.title,
      description: section.description ?? "",
      order_index: section.order_index ?? section.order ?? 0,
      is_published: section.is_published ?? true,
    },
  });

  async function onSubmit(values: SectionFormInput) {
    try {
      await clientAction(`/api/teacher/sections/${section.id}`, values, "PUT");
      toast.success("Section updated successfully.");
      onOpenChange(false);
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update section.");
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md" aria-describedby={undefined}>
        <DialogHeader>
          <DialogTitle>Edit Section</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
          <div className="grid gap-1.5">
            <Label htmlFor="ed-sec-title">Title</Label>
            <Input id="ed-sec-title" {...register("title")} />
            {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
          </div>

          <div className="grid gap-1.5">
            <Label htmlFor="ed-sec-desc">Description (optional)</Label>
            <Textarea
              id="ed-sec-desc"
              rows={3}
              placeholder="Brief overview of what this section covers…"
              {...register("description")}
            />
          </div>

          <div className="grid grid-cols-2 gap-3 items-center">
            <div className="grid gap-1.5">
              <Label htmlFor="ed-sec-order">Order index</Label>
              <Input
                id="ed-sec-order"
                type="number"
                min="0"
                {...register("order_index", { valueAsNumber: true })}
              />
            </div>
            <div className="flex items-center gap-2 pt-5">
              <input
                type="checkbox"
                id="ed-sec-pub"
                {...register("is_published")}
                className="rounded"
              />
              <Label htmlFor="ed-sec-pub" className="cursor-pointer text-sm">
                Published
              </Label>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                  Saving…
                </>
              ) : (
                "Save Changes"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function DeleteSectionDialog({
  section,
  open,
  onOpenChange,
}: {
  section: Section;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const router = useRouter();
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    setDeleting(true);
    try {
      await clientAction(`/api/teacher/sections/${section.id}`, undefined, "DELETE");
      toast.success(`Section "${section.title}" deleted.`);
      onOpenChange(false);
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete section.");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={(o) => !deleting && onOpenChange(o)}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle className="text-sm font-semibold flex items-center gap-2 text-destructive">
            <Trash2 className="h-4 w-4" /> Delete Section?
          </DialogTitle>
          <DialogDescription className="text-xs pt-1">
            This will permanently remove <strong>{section.title}</strong> and all lessons inside
            it. This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onOpenChange(false)}
            disabled={deleting}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={handleDelete}
            disabled={deleting}
          >
            {deleting ? (
              <>
                <Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
                Deleting…
              </>
            ) : (
              "Delete Section"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function SectionRow({
  section,
  lessons,
  quizzes,
  index,
  canEdit,
}: {
  section: Section;
  lessons: Lesson[];
  quizzes?: Quiz[];
  index: number;
  canEdit: boolean;
}) {
  const [expanded, setExpanded] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  return (
    <div className="border rounded-lg overflow-hidden">
      <div className="flex items-center justify-between px-4 py-3 hover:bg-muted/30 transition-colors">
        <button
          type="button"
          className="flex-1 flex items-center gap-3 text-left min-w-0"
          onClick={() => setExpanded((e) => !e)}
        >
          {expanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0" />
          ) : (
            <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
          )}
          <span className="w-6 h-6 rounded-full bg-muted flex items-center justify-center text-xs font-medium text-muted-foreground shrink-0">
            {section.order_index ?? section.order ?? index + 1}
          </span>
          <div className="flex flex-col min-w-0">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium truncate">{section.title}</span>
              {section.is_published === false && (
                <Badge variant="outline" className="text-xs h-5">
                  Draft
                </Badge>
              )}
            </div>
            {section.description && (
              <span className="text-xs text-muted-foreground line-clamp-1">
                {section.description}
              </span>
            )}
          </div>
        </button>

        <div className="flex items-center gap-2 shrink-0 ml-3">
          <span className="text-xs text-muted-foreground tabular-nums">
            {lessons.length} lesson{lessons.length !== 1 ? "s" : ""}
          </span>
          {canEdit && (
            <div className="flex items-center gap-1 border-l pl-2 ml-1">
              <Button
                type="button"
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0 text-muted-foreground hover:text-foreground"
                title="Edit Section"
                onClick={(e) => {
                  e.stopPropagation();
                  setEditOpen(true);
                }}
              >
                <Pencil className="h-3.5 w-3.5" />
              </Button>
              <Button
                type="button"
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                title="Delete Section"
                onClick={(e) => {
                  e.stopPropagation();
                  setDeleteOpen(true);
                }}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </div>
          )}
        </div>
      </div>

      {expanded && (
        <div className="border-t px-4 py-3 bg-muted/10">
          <LessonsList
            sectionId={section.id}
            lessons={lessons}
            quizzes={quizzes}
            canEdit={canEdit}
          />
        </div>
      )}

      {canEdit && (
        <>
          <EditSectionDialog
            section={section}
            open={editOpen}
            onOpenChange={setEditOpen}
          />
          <DeleteSectionDialog
            section={section}
            open={deleteOpen}
            onOpenChange={setDeleteOpen}
          />
        </>
      )}
    </div>
  );
}

export function SectionsList({
  courseId,
  sections,
  lessonsBySectionId,
  quizzes,
  canEdit = true,
}: SectionsListProps) {
  const [open, setOpen] = useState(false);
  const router = useRouter();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<SectionFormInput>({
    resolver: zodResolver(SectionFormSchema),
    defaultValues: {
      order_index: sections.length + 1,
      is_published: true,
    },
  });

  async function onSubmit(values: SectionFormInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/sections`, values);
      toast.success("Section added.");
      setOpen(false);
      reset({ order_index: sections.length + 2, is_published: true });
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
            <p className="text-sm text-muted-foreground">
              No sections yet. Add your first section to get started.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="flex flex-col gap-2">
          {sections.map((section, idx) => (
            <SectionRow
              key={section.id}
              section={section}
              lessons={lessonsBySectionId[section.id] ?? []}
              quizzes={quizzes}
              index={idx}
              canEdit={canEdit}
            />
          ))}
        </div>
      )}

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Add Section</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="section-title">Title</Label>
              <Input
                id="section-title"
                placeholder="e.g. Introduction to Algebra"
                {...register("title")}
              />
              {errors.title && (
                <p className="text-xs text-destructive">{errors.title.message}</p>
              )}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="section-desc">Description (optional)</Label>
              <Textarea
                id="section-desc"
                rows={3}
                placeholder="Brief overview of what this section covers…"
                {...register("description")}
              />
            </div>

            <div className="grid grid-cols-2 gap-3 items-center">
              <div className="grid gap-1.5">
                <Label htmlFor="section-order">Order index (optional)</Label>
                <Input
                  id="section-order"
                  type="number"
                  min="0"
                  placeholder="Auto"
                  {...register("order_index", { valueAsNumber: true })}
                />
              </div>
              <div className="flex items-center gap-2 pt-5">
                <input
                  type="checkbox"
                  id="section-pub"
                  {...register("is_published")}
                  className="rounded"
                />
                <Label htmlFor="section-pub" className="cursor-pointer text-sm">
                  Published
                </Label>
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setOpen(false);
                  reset();
                }}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                    Adding…
                  </>
                ) : (
                  "Add Section"
                )}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
