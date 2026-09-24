"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  MoreHorizontal,
  Pencil,
  Trash2,
  Globe,
  Loader2,
  Image as ImageIcon,
} from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useQueryClient } from "@tanstack/react-query";
import { clientAction } from "@/lib/client-action";
import {
  CreateCourseInputSchema,
  type CreateCourseInput,
  type Course,
} from "@/lib/api/schemas";
import { TeacherThumbnailDialog } from "./teacher-thumbnail-dialog";

interface TeacherCourseActionsProps {
  course: Course;
}

export function TeacherCourseActions({ course }: TeacherCourseActionsProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [dialog, setDialog] = useState<
    "edit" | "publish" | "delete" | "thumbnail" | null
  >(null);
  const [busy, setBusy] = useState(false);

  const existingTags = Array.isArray(course.metadata?.tags)
    ? (course.metadata.tags as string[]).join(", ")
    : "";
  const [tagsInput, setTagsInput] = useState(existingTags);

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<CreateCourseInput>({
    resolver: zodResolver(CreateCourseInputSchema),
    defaultValues: {
      title: course.title,
      description: course.description ?? "",
      short_description: course.short_description ?? "",
      subject: course.subject,
      grade: course.grade,
      level: (course.level as CreateCourseInput["level"]) ?? "beginner",
      price: course.price ?? 0,
      currency: course.currency ?? "ETB",
      is_free: course.is_free ?? false,
      language: course.language ?? "en",
    },
  });

  const isFree = watch("is_free");

  async function invalidate() {
    await queryClient.invalidateQueries({ queryKey: ["teacher-courses"] });
    router.refresh();
  }

  async function handleEdit(values: CreateCourseInput) {
    const tags = tagsInput
      .split(",")
      .map((t) => t.trim())
      .filter(Boolean);

    const payload = {
      ...values,
      price: isFree ? 0 : values.price,
      is_free: isFree,
    };

    const res = await fetch(`/api/teacher/courses/${course.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(
        (data as { message?: string }).message ?? "Failed to update course."
      );
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
      const res = await fetch(`/api/teacher/courses/${course.id}`, {
        method: "DELETE",
      });
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
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={(e) => e.stopPropagation()}
          >
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            onClick={(e) => {
              e.stopPropagation();
              setDialog("edit");
            }}
          >
            <Pencil className="h-4 w-4 mr-2" /> Edit Details
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={(e) => {
              e.stopPropagation();
              setDialog("thumbnail");
            }}
          >
            <ImageIcon className="h-4 w-4 mr-2" /> Change Thumbnail
          </DropdownMenuItem>
          {!course.is_published && (
            <DropdownMenuItem
              onClick={(e) => {
                e.stopPropagation();
                setDialog("publish");
              }}
            >
              <Globe className="h-4 w-4 mr-2" /> Publish
            </DropdownMenuItem>
          )}
          <DropdownMenuSeparator />
          {!course.is_published && (
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={(e) => {
                e.stopPropagation();
                setDialog("delete");
              }}
            >
              <Trash2 className="h-4 w-4 mr-2" /> Delete
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Edit dialog */}
      <Dialog
        open={dialog === "edit"}
        onOpenChange={(v) => !v && setDialog(null)}
      >
        <DialogContent
          className="max-w-lg max-h-[85vh] overflow-y-auto"
          aria-describedby={undefined}
        >
          <DialogHeader>
            <DialogTitle>Edit Course Details</DialogTitle>
          </DialogHeader>
          <form
            onSubmit={handleSubmit(handleEdit)}
            className="flex flex-col gap-4 py-2"
          >
            <div className="grid gap-1.5">
              <Label htmlFor="edit-title">Course Title</Label>
              <Input id="edit-title" {...register("title")} />
              {errors.title && (
                <p className="text-xs text-destructive">
                  {errors.title.message}
                </p>
              )}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="edit-short-desc">Short Summary / Tagline</Label>
              <Input
                id="edit-short-desc"
                placeholder="Brief 1-sentence summary"
                {...register("short_description")}
              />
              {errors.short_description && (
                <p className="text-xs text-destructive">
                  {errors.short_description.message}
                </p>
              )}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="edit-desc">Full Description</Label>
              <textarea
                id="edit-desc"
                rows={3}
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("description")}
              />
              {errors.description && (
                <p className="text-xs text-destructive">
                  {errors.description.message}
                </p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="edit-subject">Subject</Label>
                <Input id="edit-subject" {...register("subject")} />
                {errors.subject && (
                  <p className="text-xs text-destructive">
                    {errors.subject.message}
                  </p>
                )}
              </div>
              <div className="grid gap-1.5">
                <Label>Grade Level</Label>
                <Select
                  defaultValue={String(course.grade ?? 12)}
                  onValueChange={(v) => setValue("grade", Number(v))}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {Array.from({ length: 12 }, (_, i) => i + 1).map((g) => (
                      <SelectItem key={g} value={String(g)}>
                        Grade {g}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.grade && (
                  <p className="text-xs text-destructive">
                    {errors.grade.message}
                  </p>
                )}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Difficulty Level</Label>
                <Select
                  defaultValue={course.level ?? "beginner"}
                  onValueChange={(v) =>
                    setValue("level", v as CreateCourseInput["level"])
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="beginner">Beginner</SelectItem>
                    <SelectItem value="intermediate">Intermediate</SelectItem>
                    <SelectItem value="advanced">Advanced</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-1.5">
                <Label>Language</Label>
                <Select
                  defaultValue={course.language ?? "en"}
                  onValueChange={(v) => setValue("language", v)}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="en">English</SelectItem>
                    <SelectItem value="am">Amharic</SelectItem>
                    <SelectItem value="om">Afaan Oromoo</SelectItem>
                    <SelectItem value="ti">Tigrinya</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="border rounded-md p-3 bg-muted/20 flex flex-col gap-3">
              <div className="flex items-center justify-between">
                <Label htmlFor="edit-is-free" className="text-xs font-semibold cursor-pointer">
                  Free Course
                </Label>
                <input
                  type="checkbox"
                  id="edit-is-free"
                  {...register("is_free")}
                  className="rounded"
                />
              </div>

              {!isFree && (
                <div className="grid grid-cols-2 gap-3 pt-1 border-t">
                  <div className="grid gap-1.5">
                    <Label htmlFor="edit-price" className="text-xs">
                      Price
                    </Label>
                    <Input
                      id="edit-price"
                      type="number"
                      min="0"
                      className="text-xs"
                      {...register("price", { valueAsNumber: true })}
                    />
                  </div>
                  <div className="grid gap-1.5">
                    <Label className="text-xs">Currency</Label>
                    <Select
                      defaultValue={course.currency ?? "ETB"}
                      onValueChange={(v) => setValue("currency", v)}
                    >
                      <SelectTrigger className="text-xs">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="ETB">ETB (Ethiopian Birr)</SelectItem>
                        <SelectItem value="USD">USD ($)</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              )}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="edit-tags" className="text-xs">
                Tags / Keywords (comma-separated)
              </Label>
              <Input
                id="edit-tags"
                placeholder="e.g. National Exam, Calculus, Unit 1"
                value={tagsInput}
                onChange={(e) => setTagsInput(e.target.value)}
                className="text-xs"
              />
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setDialog(null)}
              >
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

      {/* Thumbnail dialog */}
      <TeacherThumbnailDialog
        course={course}
        open={dialog === "thumbnail"}
        onOpenChange={(v) => !v && setDialog(null)}
      />

      {/* Publish confirm */}
      <Dialog
        open={dialog === "publish"}
        onOpenChange={(v) => !v && setDialog(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Publish course?</DialogTitle>
            <DialogDescription>
              <strong>{course.title}</strong> will be visible to all students. This
              cannot be undone from the panel.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialog(null)}>
              Cancel
            </Button>
            <Button onClick={handlePublish} disabled={busy}>
              {busy ? (
                <>
                  <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                  Publishing…
                </>
              ) : (
                "Publish"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete confirm */}
      <Dialog
        open={dialog === "delete"}
        onOpenChange={(v) => !v && setDialog(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete course?</DialogTitle>
            <DialogDescription>
              <strong>{course.title}</strong> will be permanently deleted. This
              cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialog(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={busy}>
              {busy ? (
                <>
                  <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                  Deleting…
                </>
              ) : (
                "Delete"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
