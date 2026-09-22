"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from "@/components/ui/dialog";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import { CreateCourseInputSchema, type CreateCourseInput } from "@/lib/api/schemas";

export function CreateCourseDialog() {
  const [open, setOpen] = useState(false);
  const [tagsInput, setTagsInput] = useState("");
  const router = useRouter();
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateCourseInput>({
    resolver: zodResolver(CreateCourseInputSchema),
    defaultValues: { price: 0, currency: "ETB", language: "en", is_free: false, level: "beginner", grade: 12 },
  });

  const isFree = watch("is_free");

  async function onSubmit(values: CreateCourseInput) {
    try {
      const tags = tagsInput
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean);

      const payload: CreateCourseInput = {
        title: values.title.trim(),
        description: values.description.trim(),
        short_description: values.short_description?.trim() || undefined,
        subject: values.subject.trim(),
        grade: Number(values.grade) || 12,
        level: values.level || "beginner",
        language: values.language || "en",
        currency: values.currency || "ETB",
        price: isFree ? 0 : (typeof values.price === "number" && !Number.isNaN(values.price) ? values.price : 0),
        is_free: isFree,
      };

      await clientAction("/api/teacher/courses", payload);
      toast.success("Course created.");
      setOpen(false);
      reset();
      setTagsInput("");
      await queryClient.invalidateQueries({ queryKey: ["teacher-courses"] });
      await queryClient.invalidateQueries({ queryKey: ["courses"] });
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create course.");
    }
  }

  return (
    <>
      <Button size="sm" onClick={() => setOpen(true)}>
        <Plus className="h-4 w-4 mr-1.5" />
        New Course
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-lg max-h-[85vh] overflow-y-auto" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Create Course</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="title">Course Title</Label>
              <Input id="title" placeholder="e.g. Grade 12 National Exam Preparation - Mathematics" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="short_description">Short Summary / Tagline (optional)</Label>
              <Input
                id="short_description"
                placeholder="A brief 1-sentence overview for course cards"
                {...register("short_description")}
              />
              {errors.short_description && <p className="text-xs text-destructive">{errors.short_description.message}</p>}
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="description">Full Description</Label>
              <textarea
                id="description"
                rows={3}
                placeholder="What will students learn in this course?"
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("description")}
              />
              {errors.description && <p className="text-xs text-destructive">{errors.description.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="subject">Subject</Label>
                <Input id="subject" placeholder="e.g. Mathematics" {...register("subject")} />
                {errors.subject && <p className="text-xs text-destructive">{errors.subject.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label>Grade Level</Label>
                <Select defaultValue="12" onValueChange={(v) => setValue("grade", Number(v))}>
                  <SelectTrigger><SelectValue placeholder="Select grade" /></SelectTrigger>
                  <SelectContent>
                    {Array.from({ length: 12 }, (_, i) => i + 1).map((g) => (
                      <SelectItem key={g} value={String(g)}>Grade {g}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.grade && <p className="text-xs text-destructive">{errors.grade.message}</p>}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Difficulty Level</Label>
                <Select defaultValue="beginner" onValueChange={(v) => setValue("level", v as CreateCourseInput["level"])}>
                  <SelectTrigger><SelectValue placeholder="Select level" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="beginner">Beginner</SelectItem>
                    <SelectItem value="intermediate">Intermediate</SelectItem>
                    <SelectItem value="advanced">Advanced</SelectItem>
                  </SelectContent>
                </Select>
                {errors.level && <p className="text-xs text-destructive">{errors.level.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label>Language</Label>
                <Select defaultValue="en" onValueChange={(v) => setValue("language", v)}>
                  <SelectTrigger><SelectValue placeholder="Select language" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="en">English (EN)</SelectItem>
                    <SelectItem value="am">Amharic (AM)</SelectItem>
                    <SelectItem value="om">Afaan Oromo (OM)</SelectItem>
                    <SelectItem value="ti">Tigrinya (TI)</SelectItem>
                  </SelectContent>
                </Select>
                {errors.language && <p className="text-xs text-destructive">{errors.language.message}</p>}
              </div>
            </div>

            <div className="rounded-lg border p-3 flex flex-col gap-3 bg-muted/20">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="is_free"
                  checked={isFree}
                  onChange={(e) => {
                    setValue("is_free", e.target.checked);
                    if (e.target.checked) setValue("price", 0);
                  }}
                  className="rounded h-4 w-4"
                />
                <Label htmlFor="is_free" className="text-xs font-semibold cursor-pointer">
                  This is a Free Course (Open to all students)
                </Label>
              </div>

              {!isFree && (
                <div className="grid grid-cols-2 gap-3 pt-1">
                  <div className="grid gap-1.5">
                    <Label htmlFor="price" className="text-xs">Price</Label>
                    <Input
                      id="price"
                      type="number"
                      min="0"
                      step="any"
                      placeholder="e.g. 250"
                      {...register("price", { valueAsNumber: true })}
                    />
                    {errors.price && <p className="text-xs text-destructive">{errors.price.message}</p>}
                  </div>
                  <div className="grid gap-1.5">
                    <Label className="text-xs">Currency</Label>
                    <Select defaultValue="ETB" onValueChange={(v) => setValue("currency", v)}>
                      <SelectTrigger className="text-xs"><SelectValue /></SelectTrigger>
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
              <Label htmlFor="tags" className="text-xs">Tags / Keywords (comma-separated, optional)</Label>
              <Input
                id="tags"
                placeholder="e.g. National Exam, Calculus, Unit 1"
                value={tagsInput}
                onChange={(e) => setTagsInput(e.target.value)}
                className="text-xs"
              />
            </div>

            <DialogFooter className="pt-2">
              <Button type="button" variant="outline" onClick={() => { setOpen(false); reset(); setTagsInput(""); }}>
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Creating…</> : "Create Course"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}
