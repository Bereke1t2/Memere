"use client";

import { usePathname, useRouter } from "next/navigation";
import { Menu } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme-toggle";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { useBreadcrumb } from "@/lib/breadcrumb-context";
import type { User } from "@/lib/api/schemas";

const ROUTE_LABELS: Record<string, string> = {
  "/": "Dashboard",
  "/users": "Users",
  "/courses": "Courses",
  "/payments": "Payments",
  "/revenue": "Revenue",
  "/announcements": "Announcements",
};

function initials(user: User): string {
  return `${user.first_name[0] ?? ""}${user.last_name[0] ?? ""}`.toUpperCase();
}

function breadcrumb(pathname: string): { parent?: string; current: string } {
  const label = ROUTE_LABELS[pathname];
  if (label) return { parent: label === "Dashboard" ? undefined : "Dashboard", current: label };
  const parts = pathname.split("/").filter(Boolean);
  return { parent: ROUTE_LABELS[`/${parts[0]}`] ?? parts[0], current: parts.at(-1) ?? "Page" };
}

interface HeaderProps {
  user: User;
  onMenuClick: () => void;
}

export function Header({ user, onMenuClick }: HeaderProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { label: contextLabel } = useBreadcrumb();

  // Detail pages set a contextLabel (e.g. user name, course title) via BreadcrumbSetter.
  const raw = breadcrumb(pathname);
  const crumb = contextLabel
    ? { parent: raw.parent, current: contextLabel }
    : raw;

  async function handleLogout() {
    await fetch("/api/auth/logout", { method: "POST" });
    router.push("/login");
  }

  return (
    <header className="flex h-14 shrink-0 items-center border-b bg-card px-4 gap-3">
      <Button
        variant="ghost"
        size="icon"
        className="md:hidden"
        aria-label="Open navigation menu"
        onClick={onMenuClick}
      >
        <Menu className="h-4 w-4" />
      </Button>

      <nav aria-label="Breadcrumb" className="flex-1">
        <ol className="flex items-center gap-1.5 text-sm">
          {crumb.parent ? (
            <>
              <li className="text-muted-foreground">{crumb.parent}</li>
              <li aria-hidden="true" className="text-muted-foreground">/</li>
              <li className="font-medium capitalize">{crumb.current}</li>
            </>
          ) : (
            <li className="font-medium">{crumb.current}</li>
          )}
        </ol>
      </nav>

      <ThemeToggle />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            className="h-8 w-8 rounded-full p-0"
            aria-label="Open user menu"
          >
            <Avatar className="h-8 w-8">
              <AvatarFallback className="text-xs">{initials(user)}</AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-52">
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col gap-0.5">
              <span className="text-sm font-medium">
                {user.first_name} {user.last_name}
              </span>
              <span className="text-xs text-muted-foreground">{user.email}</span>
              <span className="text-xs text-muted-foreground capitalize">{user.role}</span>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onSelect={handleLogout}
            className="text-destructive focus:text-destructive cursor-pointer"
          >
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
  );
}
