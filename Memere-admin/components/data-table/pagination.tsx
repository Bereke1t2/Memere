"use client";

import { useState } from "react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const LIMIT_OPTIONS = ["10", "20", "50"] as const;

interface PaginationProps {
  nextCursor?: string;
  isLoading?: boolean;
}

export function Pagination({ nextCursor, isLoading }: PaginationProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [cursorStack, setCursorStack] = useState<string[]>([]);

  const limit = searchParams.get("limit") ?? "20";
  const currentAfter = searchParams.get("after") ?? "";

  // Derive a key from filter params (excluding cursor/limit).
  const filterKey = (() => {
    const p = new URLSearchParams(searchParams.toString());
    p.delete("after");
    p.delete("limit");
    return p.toString();
  })();

  // Reset cursor stack when filters change — called during render (not in effect).
  const [prevFilterKey, setPrevFilterKey] = useState(filterKey);
  if (filterKey !== prevFilterKey) {
    setPrevFilterKey(filterKey);
    setCursorStack([]);
  }

  function navigate(after: string, newLimit?: string) {
    const params = new URLSearchParams(searchParams.toString());
    if (after) params.set("after", after);
    else params.delete("after");
    if (newLimit) params.set("limit", newLimit);
    router.replace(`${pathname}?${params.toString()}`);
  }

  function goNext() {
    if (!nextCursor) return;
    setCursorStack((s) => [...s, currentAfter]);
    navigate(nextCursor);
  }

  function goPrev() {
    const stack = [...cursorStack];
    const prev = stack.pop() ?? "";
    setCursorStack(stack);
    navigate(prev);
  }

  function setLimit(value: string) {
    setCursorStack([]);
    const params = new URLSearchParams(searchParams.toString());
    params.set("limit", value);
    params.delete("after");
    router.replace(`${pathname}?${params.toString()}`);
  }

  return (
    <div className="flex items-center justify-between py-3">
      <div className="flex items-center gap-2">
        <span className="text-sm text-muted-foreground">Rows per page</span>
        <Select value={limit} onValueChange={setLimit}>
          <SelectTrigger className="h-8 w-16">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {LIMIT_OPTIONS.map((o) => (
              <SelectItem key={o} value={o}>
                {o}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={goPrev}
          disabled={cursorStack.length === 0 || isLoading}
          aria-label="Previous page"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={goNext}
          disabled={!nextCursor || isLoading}
          aria-label="Next page"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
