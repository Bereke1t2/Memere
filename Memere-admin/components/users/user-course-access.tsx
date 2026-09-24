"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { ShieldCheck, CheckSquare, Square, KeyRound, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { UserCourseAccessItem } from "@/lib/api/schemas";

interface UserCourseAccessProps {
  userId: string;
}

export function UserCourseAccess({ userId }: UserCourseAccessProps) {
  const queryClient = useQueryClient();
  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  const { data, isLoading, error } = useQuery<{ student_id: string; courses: UserCourseAccessItem[] }>({
    queryKey: ["user-course-access", userId],
    queryFn: async () => {
      const res = await fetch(`/api/admin/users/${userId}/course-access`);
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to load course access list");
      }
      return res.json();
    },
  });

  const courses = data?.courses ?? [];
  const paidCoursesWithoutAccess = courses.filter((c) => !c.is_free && !c.has_access);

  const grantMutation = useMutation({
    mutationFn: async (payload: { course_ids?: string[]; all_courses?: boolean }) => {
      const res = await fetch(`/api/admin/users/${userId}/grant-courses`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to grant course access");
      }
      return res.json();
    },
    onSuccess: (_, variables) => {
      toast.success(
        variables.all_courses
          ? "Access to all courses granted successfully."
          : `Access granted to ${variables.course_ids?.length ?? 0} course(s).`
      );
      setSelectedIds([]);
      queryClient.invalidateQueries({ queryKey: ["user-course-access", userId] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  const revokeMutation = useMutation({
    mutationFn: async (courseId: string) => {
      const res = await fetch(`/api/admin/users/${userId}/course-access/${courseId}`, {
        method: "DELETE",
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to revoke course access");
      }
    },
    onSuccess: () => {
      toast.success("Course access revoked successfully.");
      queryClient.invalidateQueries({ queryKey: ["user-course-access", userId] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  function toggleCourse(id: string) {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  }

  function toggleSelectAll() {
    const allGrantableIds = paidCoursesWithoutAccess.map((c) => c.course_id);
    if (selectedIds.length === allGrantableIds.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(allGrantableIds);
    }
  }

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-48 mb-2" />
          <Skeleton className="h-4 w-72" />
        </CardHeader>
        <CardContent className="space-y-3">
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-10 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-destructive">Course Access Error</CardTitle>
          <CardDescription>{(error as Error).message}</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div>
            <CardTitle className="flex items-center gap-2 text-lg">
              <KeyRound className="h-5 w-5 text-primary" />
              Course Access Management
            </CardTitle>
            <CardDescription>
              Directly grant or revoke course access for this student.
            </CardDescription>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            {paidCoursesWithoutAccess.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={toggleSelectAll}
                disabled={grantMutation.isPending}
              >
                {selectedIds.length === paidCoursesWithoutAccess.length ? (
                  <>
                    <Square className="h-4 w-4 mr-1.5" />
                    Deselect All
                  </>
                ) : (
                  <>
                    <CheckSquare className="h-4 w-4 mr-1.5" />
                    Select All ({paidCoursesWithoutAccess.length})
                  </>
                )}
              </Button>
            )}
            <Button
              variant="default"
              size="sm"
              disabled={selectedIds.length === 0 || grantMutation.isPending}
              onClick={() => grantMutation.mutate({ course_ids: selectedIds })}
            >
              {grantMutation.isPending && !grantMutation.variables?.all_courses
                ? "Granting…"
                : `Grant Selected (${selectedIds.length})`}
            </Button>
            <Button
              variant="secondary"
              size="sm"
              disabled={grantMutation.isPending || paidCoursesWithoutAccess.length === 0}
              onClick={() => grantMutation.mutate({ all_courses: true })}
            >
              {grantMutation.isPending && grantMutation.variables?.all_courses
                ? "Granting All…"
                : "Grant All Courses"}
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {courses.length === 0 ? (
          <p className="text-sm text-muted-foreground py-4 text-center">
            No courses found on the platform.
          </p>
        ) : (
          <div className="divide-y rounded-md border text-sm">
            {courses.map((c) => {
              const isSelected = selectedIds.includes(c.course_id);
              return (
                <div
                  key={c.course_id}
                  className={`flex items-center justify-between gap-4 p-3 transition-colors ${
                    isSelected ? "bg-muted/50" : ""
                  }`}
                >
                  <div className="flex items-center gap-3 min-w-0">
                    {!c.is_free && !c.has_access ? (
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => toggleCourse(c.course_id)}
                        className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary cursor-pointer"
                        id={`course-${c.course_id}`}
                      />
                    ) : (
                      <div className="w-4" />
                    )}
                    <label
                      htmlFor={`course-${c.course_id}`}
                      className="font-medium truncate cursor-pointer select-none"
                    >
                      {c.title || c.course_title || "Untitled Course"}
                    </label>
                  </div>

                  <div className="flex items-center gap-3 shrink-0">
                    {c.is_free ? (
                      <Badge variant="outline" className="text-blue-700 bg-blue-50 border-blue-200">
                        Free Course
                      </Badge>
                    ) : (
                      <span className="text-xs text-muted-foreground">
                        {c.price} {c.currency}
                      </span>
                    )}

                    {c.has_access ? (
                      <div className="flex items-center gap-2">
                        <Badge
                          variant="secondary"
                          className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950 flex items-center gap-1"
                        >
                          <ShieldCheck className="h-3 w-3" />
                          Access Granted {c.access_source ? `(${c.access_source})` : ""}
                        </Badge>
                        {!c.is_free && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-7 w-7 text-destructive hover:bg-destructive/10"
                            title="Revoke access"
                            disabled={revokeMutation.isPending}
                            onClick={() => revokeMutation.mutate(c.course_id)}
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                    ) : (
                      <Badge variant="outline" className="text-muted-foreground">
                        No Access
                      </Badge>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
