"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Image as ImageIcon, Upload, Loader2, X, Sparkles } from "lucide-react";
import { toast } from "sonner";
import type { Course } from "@/lib/api/schemas";

interface TeacherThumbnailDialogProps {
  course: Course;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function TeacherThumbnailDialog({
  course,
  open,
  onOpenChange,
}: TeacherThumbnailDialogProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  // Active display image: staged file preview, or existing course thumbnail
  const currentThumb = previewUrl || course.thumbnail_url || null;

  function resetSelection() {
    setSelectedFile(null);
    setPreviewUrl(null);
    setIsDragging(false);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen && !isUploading) {
      resetSelection();
    }
    onOpenChange(nextOpen);
  }

  function processFile(file: File) {
    if (!file.type.startsWith("image/")) {
      toast.error("Please select a valid image file (PNG, JPG, WebP).");
      return;
    }
    const maxBytes = 5 * 1024 * 1024; // 5MB limit
    if (file.size > maxBytes) {
      toast.error(
        `Image must be under 5MB. Selected file is ${(file.size / (1024 * 1024)).toFixed(1)}MB.`
      );
      return;
    }
    setSelectedFile(file);
    setPreviewUrl(URL.createObjectURL(file));
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (file) processFile(file);
  }

  function handleDragOver(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }

  function handleDragLeave(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    const file = e.dataTransfer.files?.[0];
    if (file) processFile(file);
  }

  async function handleUpload() {
    if (!selectedFile) return;

    setIsUploading(true);
    try {
      const formData = new FormData();
      formData.append("file", selectedFile);

      const res = await fetch(`/api/teacher/courses/${course.id}/thumbnail`, {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.message || "Thumbnail upload failed.");
      }

      toast.success("Course thumbnail updated successfully!");
      await queryClient.invalidateQueries({ queryKey: ["teacher-courses"] });
      await queryClient.invalidateQueries({ queryKey: ["courses"] });
      router.refresh();
      handleOpenChange(false);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to upload thumbnail."
      );
    } finally {
      setIsUploading(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-md" aria-describedby="thumb-desc">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ImageIcon className="h-5 w-5 text-primary" />
            Course Thumbnail
          </DialogTitle>
          <DialogDescription id="thumb-desc" className="text-xs">
            Upload a cover image for <strong>{course.title}</strong>. Recommended ratio: 16:9 (1280×720). Max size: 5MB.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-4 py-2">
          {/* Dropzone & Preview Box */}
          <div
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            onClick={() => fileInputRef.current?.click()}
            className={`group relative rounded-xl border-2 border-dashed aspect-video overflow-hidden cursor-pointer transition-all duration-200 flex flex-col items-center justify-center p-2 ${
              isDragging
                ? "border-primary bg-primary/10 ring-2 ring-primary/20"
                : "border-border hover:border-primary/60 hover:bg-muted/40 bg-muted/20"
            }`}
          >
            {currentThumb ? (
              <>
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={currentThumb}
                  alt="Course thumbnail preview"
                  className="w-full h-full object-cover rounded-lg"
                />
                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-1.5 text-white p-4 text-center rounded-lg">
                  <Upload className="h-6 w-6" />
                  <p className="text-xs font-medium">Click or drag image to replace</p>
                </div>
              </>
            ) : (
              <div className="flex flex-col items-center justify-center text-muted-foreground gap-2 p-4 text-center">
                <div className="p-3 rounded-full bg-muted border">
                  <Upload className="h-6 w-6 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-xs font-semibold text-foreground">
                    Click to browse or drag image here
                  </p>
                  <p className="text-[11px] text-muted-foreground mt-0.5">
                    PNG, JPG, WebP up to 5MB (16:9 recommended)
                  </p>
                </div>
              </div>
            )}

            <input
              ref={fileInputRef}
              type="file"
              accept="image/png, image/jpeg, image/webp, image/jpg"
              className="hidden"
              disabled={isUploading}
              onChange={handleFileChange}
            />
          </div>

          {/* Staged File Info Banner */}
          {selectedFile && (
            <div className="flex items-center justify-between gap-2 p-2.5 rounded-lg border bg-muted/40 text-xs">
              <div className="flex items-center gap-2 min-w-0">
                <Sparkles className="h-4 w-4 text-primary shrink-0" />
                <span className="font-medium truncate">{selectedFile.name}</span>
                <span className="text-muted-foreground shrink-0 tabular-nums">
                  ({(selectedFile.size / (1024 * 1024)).toFixed(2)} MB)
                </span>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="h-6 w-6 shrink-0 text-muted-foreground hover:text-foreground"
                onClick={(e) => {
                  e.stopPropagation();
                  resetSelection();
                }}
                disabled={isUploading}
              >
                <X className="h-3.5 w-3.5" />
              </Button>
            </div>
          )}
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isUploading}
          >
            Cancel
          </Button>

          {selectedFile ? (
            <Button
              type="button"
              onClick={handleUpload}
              disabled={isUploading}
            >
              {isUploading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                  Saving…
                </>
              ) : (
                "Save Thumbnail"
              )}
            </Button>
          ) : (
            <Button
              type="button"
              variant="secondary"
              onClick={() => fileInputRef.current?.click()}
            >
              <Upload className="h-3.5 w-3.5 mr-1.5" />
              {currentThumb ? "Change Image" : "Select Image"}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
