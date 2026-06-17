"use client";

import { useRef, useState } from "react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface SearchFilterProps {
  placeholder?: string;
  paramKey?: string;
}

export function SearchFilter({
  placeholder = "Search...",
  paramKey = "q",
}: SearchFilterProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [value, setValue] = useState(searchParams.get(paramKey) ?? "");
  const timerRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  function handleChange(v: string) {
    setValue(v);
    clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => {
      const params = new URLSearchParams(searchParams.toString());
      if (v) params.set(paramKey, v);
      else params.delete(paramKey);
      params.delete("after"); // reset cursor on new filter
      router.replace(`${pathname}?${params.toString()}`);
    }, 400);
  }

  return (
    <div className="relative">
      <Search className="absolute left-2.5 top-2 h-4 w-4 text-muted-foreground pointer-events-none" />
      <Input
        className="pl-8 h-8 w-56"
        placeholder={placeholder}
        value={value}
        onChange={(e) => handleChange(e.target.value)}
      />
    </div>
  );
}

interface SelectFilterOption {
  label: string;
  value: string;
}

interface SelectFilterProps {
  options: SelectFilterOption[];
  placeholder?: string;
  paramKey?: string;
}

const ALL = "__all__";

export function SelectFilter({
  options,
  placeholder = "All",
  paramKey = "filter",
}: SelectFilterProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const current = searchParams.get(paramKey) ?? ALL;

  function handleChange(value: string) {
    const params = new URLSearchParams(searchParams.toString());
    if (value && value !== ALL) params.set(paramKey, value);
    else params.delete(paramKey);
    params.delete("after"); // reset cursor on new filter
    router.replace(`${pathname}?${params.toString()}`);
  }

  return (
    <Select value={current} onValueChange={handleChange}>
      <SelectTrigger className="h-8 w-36">
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value={ALL}>{placeholder}</SelectItem>
        {options.map((o) => (
          <SelectItem key={o.value} value={o.value}>
            {o.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
