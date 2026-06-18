"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { MoreHorizontal, Pencil, Trash2, Globe, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem,
  DropdownMenuSeparator, DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle,
  DialogFooter, DialogDescription,
} from "@/components/ui/dialog";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { useQueryClient } from "@tanstack/react-query";
import { clientAction } from "@/lib/client-action";
import { CreateCourseInputSchema, type CreateCourseInput, type Course } from "@/lib/api/schemas";

interface TeacherCourseActionsProps {
  course: Course;
}

export function TeacherCourseActions({ course }: TeacherCourseActionsProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [dialog, setDialog] = useState<"edit" | "publish" | "delete" | null>(null);
  const [busy, setBusy] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CreateCourseInput>({
    resolver: zodResolver(CreateCourseInputSchema),
    defaultValues: {
      title: course.title,
      description: course.description,
      subject: course.subject,
      grade: course.grade,
      level: (course.level as CreateCourseInput["level"]) ?? "beginner",
      price: course.price,
      language: course.language ?? "en",
    },
  });

  async function invalidate() {
    await queryClient.invalidateQueries({ queryKey: ["teacher-courses"] });
    router.refresh();
  }

  async function handleEdit(values: CreateCourseInput) {
    const res = await fetch(`/api/teacher/courses/${course.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(values),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error((data as { message?: string }).message ?? "Failed to update course.");
      return;
    }
    toast.success("Course updated.");
    setDialog(null);
    await invalidate();
  }

  async function handlePublish() {
    setBusy(true);
    try {
      await clientAction(`/api/teacher/courses/${course.id}/publish`);
      toast.success("Course published.");
      setDialog(null);
      await invalidate();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to publish.");
    } finally {
      setBusy(false);
    }
  }

  async function handleDelete() {
    setBusy(true);
    try {
      const res = await fetch(`/api/teacher/courses/${course.id}`, { method: "DELETE" });
      if (!res.ok) throw new Error("Failed to delete course.");
      toast.success("Course deleted.");
      setDialog(null);
      await invalidate();
    } catch {
      toast.error("Failed to delete course.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" className="h-8 w-8" onClick={(e) => e.stopPropagation()}>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={(e) => { e.stopPropagation(); setDialog("edit"); }}>
            <Pencil className="h-4 w-4 mr-2" /> Edit
          </DropdownMenuItem>
          {!course.is_published && (
            <DropdownMenuItem onClick={(e) => { e.stopPropagation(); setDialog("publish"); }}>
              <Globe className="h-4 w-4 mr-2" /> Publish
            </DropdownMenuItem>
          )}
          <DropdownMenuSeparator />
          {!course.is_published && (
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={(e) => { e.stopPropagation(); setDialog("delete"); }}
            >
              <Trash2 className="h-4 w-4 mr-2" /> Delete
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Edit dialog */}
      <Dialog open={dialog === "edit"} onOpenChange={(v) => !v && setDialog(null)}>
        <DialogContent className="max-w-lg" aria-describedby={undefined}>
          <DialogHeader><DialogTitle>Edit Course</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(handleEdit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="edit-title">Title</Label>
              <Input id="edit-title" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="edit-desc">Description</Label>
              <textarea
                id="edit-desc"
                rows={3}
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("description")}
              />
              {errors.description && <p className="text-xs text-destructive">{errors.description.message}</p>}
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="edit-subject">Subject</Label>
                <Input id="edit-subject" {...register("subject")} />
              </div>
              <div className="grid gap-1.5">
                <Label>Level</Label>
                <Select
                  defaultValue={course.level ?? "beginner"}
                  onValueChange={(v) => setValue("level", v as CreateCourseInput["level"])}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="beginner">Beginner</SelectItem>
                    <SelectItem value="intermediate">Intermediate</SelectItem>
                    <SelectItem value="advanced">Advanced</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="edit-price">Price (ETB, 0 = free)</Label>
              <Input id="edit-price" type="number" min="0" {...register("price", { valueAsNumber: true })} />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setDialog(null)}>Cancel</Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Saving…</> : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Publish confirm */}
      <Dialog open={dialog === "publish"} onOpenChange={(v) => !v && setDialog(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Publish course?</DialogTitle>
            <DialogDescription>
              <strong>{course.title}</strong> will be visible to all students. This cannot be undone from the panel.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialog(null)}>Cancel</Button>
            <Button onClick={handlePublish} disabled={busy}>
              {busy ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Publishing…</> : "Publish"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete confirm */}
      <Dialog open={dialog === "delete"} onOpenChange={(v) => !v && setDialog(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete course?</DialogTitle>
            <DialogDescription>
              <strong>{course.title}</strong> will be permanently deleted. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialog(null)}>Cancel</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={busy}>
              {busy ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Deleting…</> : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
